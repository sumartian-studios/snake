// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/sumartian-studios/snake/configuration"
)

var forceUpdateFlag bool
var profileFlag string
var traceFlag bool

// Returns the profile and if it has changed.
func getOrUpdateCurrentProfile() (*configuration.Profile, bool, error) {
	var currentProfile *configuration.Profile = nil
	var currentProfileExists = false

	if currentProfileExists, currentProfile = app.getCurrentProfile(); currentProfileExists {
		fmt.Println("Reusing profile:", currentProfile.Name)
	}

	if len(app.cfg.Profiles) < 1 {
		return nil, false, fmt.Errorf("project does not have any profiles: %s", profileFlag)
	}

	// Look for a profile matching flag name.
	for i, p := range app.cfg.Profiles {
		if p.Name == profileFlag {
			// Requested same profile; do nothing...
			if currentProfile != nil && p.Name == currentProfile.Name {
				return currentProfile, false, nil
			}

			currentProfile = &p
			app.setCurrentProfile(i, currentProfile)

			return currentProfile, true, nil
		}
	}

	// If no profile specified and no previous configuration use the first profile available.
	if len(profileFlag) < 1 {
		currentProfile = &app.cfg.Profiles[0]
		app.setCurrentProfile(0, currentProfile)
		return currentProfile, true, nil
	}

	return nil, false, fmt.Errorf("unable to find profile (see 'snake profiles'): %s", profileFlag)
}

func prettyPrintCMakeTraceResults() error {
	file, err := os.Open(filepath.Join(app.db.ProfilePath, "cmake.trace"))

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	type TraceResult struct {
		File  string   `json:"file"`
		Line  int      `json:"line"`
		Cmd   string   `json:"cmd"`
		Args  []string `json:"args"`
		Time  float64  `json:"time"`
		Frame int      `json:"frame"`
	}

	var lastTime time.Time
	durations := map[string]time.Duration{}

	i := 0

	for scanner.Scan() {
		var trace TraceResult

		if err = json.Unmarshal(scanner.Bytes(), &trace); err != nil {
			return err
		}

		if len(trace.Cmd) < 1 || trace.Time < 1 {
			continue
		}

		sec, dec := math.Modf(trace.Time)
		t := time.Unix(int64(sec), int64(dec*(1e9)))

		if i != 0 {
			d := t.Sub(lastTime)
			durations[trace.Cmd] = durations[trace.Cmd] + d
		}

		i++
		lastTime = t
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)

	for cmd, total := range durations {
		fmt.Fprintln(writer, total.String()+"\t"+cmd)
	}

	writer.Flush()

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the build system",
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Configuration")

		fmt.Println("Configuring...")

		err := app.initSlow()

		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		var cmakeOptions []string

		// Check if the configuration changed and if so regenerate.

		currentProfile, profileChanged, err := getOrUpdateCurrentProfile()

		if err != nil {
			return err
		}

		cmakeOptions = append(cmakeOptions,
			"-DSNAKE_DIR="+app.snakeDir,
			"-B", app.db.ProfilePath, "-S", app.rootDir,
			"-G", "Ninja",
		)

		if app.verbose {
			cmakeOptions = append(cmakeOptions,
				"--warn-uninitialized", "--warn-unused-vars", "--check-system-vars")
		}

		if traceFlag {
			cmakeOptions = append(cmakeOptions, "--trace-format=json-v1",
				"--trace-redirect="+filepath.Join(app.db.ProfilePath, "cmake.trace"))
		}

		if profileChanged {
			if _, err := os.Stat(app.db.ProfilePath); os.IsNotExist(err) {
				forceUpdateFlag = true
			}
		}

		fmt.Println("Configured:", app.db.Configured)
		fmt.Println("Profile Changed:", profileChanged)
		fmt.Println("Force Update:", forceUpdateFlag)

		if !app.db.Configured || forceUpdateFlag {
			fmt.Println("Updating...")

			if err = app.decompress(); err != nil {
				return err
			}

			if !(app.db.Configured && profileChanged) {
				if err = os.RemoveAll(filepath.Join(app.snakeDir, "snake.lock")); err != nil {
					return err
				}
			}

			app.db.Configured = true
			app.storageChanged()
		}

		fmt.Println("Load profile:", currentProfile.Name)

		if len(currentProfile.Type) > 0 {
			cmakeOptions = append(cmakeOptions,
				"-DCMAKE_BUILD_TYPE="+currentProfile.Type)
		}

		if len(currentProfile.LinkFlags) > 0 {
			cmakeOptions = append(cmakeOptions,
				"-DSNAKE_GLOBAL_LINKER_OPTIONS="+strings.Join(
					strings.Split(strings.Join(currentProfile.LinkFlags, " "), " "), ";"))
		}

		if len(currentProfile.CompileFlags) > 0 {
			cmakeOptions = append(cmakeOptions,
				"-DSNAKE_GLOBAL_COMPILE_OPTIONS="+strings.Join(
					strings.Split(strings.Join(currentProfile.CompileFlags, " "), " "), ";"))
		}

		if len(currentProfile.Compiler) > 0 {
			cmakeOptions = append(cmakeOptions,
				"-DCMAKE_CXX_COMPILER="+currentProfile.Compiler)
		}

		for _, mapping := range currentProfile.Variables {
			for k, v := range mapping {
				if app.verbose {
					fmt.Println("set:", k, v)
				}

				cmakeOptions = append(cmakeOptions, fmt.Sprintf("-D%s=%s", k, v))
			}
		}

		for _, arg := range args {
			cmakeOptions = append(cmakeOptions, "-D"+arg)
		}

		if err := app.launch("cmake", cmakeOptions...); err != nil {
			return err
		}

		if err = app.saveStorage(); err != nil {
			return err
		}

		if traceFlag {
			if err = prettyPrintCMakeTraceResults(); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	configureCmd.PersistentFlags().BoolVar(&forceUpdateFlag,
		"update", false,
		"Force download and install dependencies")

	configureCmd.PersistentFlags().BoolVar(&traceFlag,
		"trace", false,
		"Trace the CMake scripts and pretty print elapsed times")

	configureCmd.PersistentFlags().StringVarP(&profileFlag,
		"profile", "p", "",
		"Select the build profile/preset")

}
