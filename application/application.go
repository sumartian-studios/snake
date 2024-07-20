// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/spf13/cobra"
	"github.com/sumartian/snake/configuration"
	"gopkg.in/yaml.v3"
)

// This is set by ldflags.
var VersionStr string

// Storage represents a persistent structure.
type Storage struct {
	// True if Conan has been configured.
	Configured bool `json:"Configured"`

	// Path to the active profile.
	ProfilePath string `json:"ProfilePath"`

	// Index of the active profile.
	ProfileIndex int `json:"ProfileIndex"`
}

// Application represents our global state manager.
type Application struct {
	cobra.Command

	// YAML configuration.
	cfg *configuration.Configuration

	// Snake build database.
	db Storage

	// Path to snake directory.
	snakeDir string

	// Path to source directory.
	rootDir string

	// Path to YAML configuration file.
	configPath string

	// Path to the storage file.
	storagePath string

	// Returns true if this is the first launch.
	firstLaunch bool

	// True if application has been initialized.
	running bool

	// True if running in interactive mode.
	interactive bool

	// True if the storage needs to be saved.
	storagePendingSave bool

	// True if verbose mode enabled.
	verbose bool
}

// Global instance of our application.
var app Application

// Creates the build directory if it does not already exist.
func (app *Application) createSnakeDir() error {
	if _, err := os.Stat(app.snakeDir); os.IsNotExist(err) {
		fmt.Println("Creating build directory")
		os.Mkdir(app.snakeDir, os.ModePerm)
	} else if os.IsExist(err) {
		fmt.Println("Build directory found:", app.snakeDir)
	} else {
		return err
	}

	return nil
}

// Common initialization.
func (app *Application) init() error {
	var err error

	if app.rootDir, err = filepath.Abs(app.rootDir); err != nil {
		return err
	}

	if app.snakeDir, err = filepath.Abs(app.snakeDir); err != nil {
		return err
	}

	app.configPath = filepath.Join(app.rootDir, ".snake.yml")
	app.storagePath = filepath.Join(app.snakeDir, "snake.db")

	return nil
}

// Slow path initialization.
func (app *Application) initSlow() error {
	if app.running {
		return nil
	}

	err := app.init()

	if err != nil {
		return err
	}

	if err = app.loadConfiguration(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err = app.createSnakeDir(); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	if _, err := os.Stat(app.snakeDir); os.IsNotExist(err) {
		if err = os.Mkdir(app.snakeDir, os.ModePerm); err != nil {
			return err
		}
	}

	if err := app.loadStorage(); err != nil {
		return err
	}

	return nil
}

// Reload configuration.
func (app *Application) loadConfiguration() error {
	app.cfg = new(configuration.Configuration)

	data, err := ioutil.ReadFile(app.configPath)

	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, app.cfg); err != nil {
		return err
	}

	return nil
}

// Fast path initialization.
func (app *Application) initFast() error {
	if app.running {
		return nil
	}

	if err := app.init(); err != nil {
		return err
	}

	if err := app.loadStorage(); err != nil {
		return err
	}

	return nil
}

// Track the time taken by a function.
func (app *Application) timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

// Launch a subprocess.
func (app *Application) launch(program string, args ...string) error {
	cmd := exec.Command(program, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Load storage from disk into memory.
func (app *Application) loadStorage() error {
	if _, err := os.Stat(app.storagePath); os.IsNotExist(err) {
		app.firstLaunch = true
		return nil
	}

	app.firstLaunch = false

	file, err := os.Open(app.storagePath)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := cbor.NewDecoder(file)

	if err = decoder.Decode(&app.db); err != nil {
		return err
	}

	return nil
}

// Needs to be called before saving storage to disk.
func (app *Application) storageChanged() {
	app.storagePendingSave = true
}

// Save storage to disk.
func (app *Application) saveStorage() error {
	if !app.storagePendingSave {
		return nil
	}

	var b bytes.Buffer

	encoder := cbor.NewEncoder(&b)

	if err := encoder.Encode(&app.db); err != nil {
		return err
	}

	if err := ioutil.WriteFile(app.storagePath, b.Bytes(), 0644); err != nil {
		return err
	}

	app.storagePendingSave = false

	return nil
}

// Start the application and parse command-line arguments.
func Execute() {
	app.Execute()
}

func init() {
	app.Command = cobra.Command{
		Use:                "snake",
		Version:            VersionStr,
		Short:              "Snake is a C++ build system and CI/CD tool designed to interoperate with CMake\nand other third-party applications and libraries.",
		SilenceUsage:       true,
		DisableFlagParsing: false,
		RunE:               startInteractiveMode,
	}

	app.Command.PersistentFlags().BoolP("help", "h", false, "Show help information")
	app.Command.PersistentFlags().BoolP("version", "v", false, "Show version information")

	app.Command.PersistentFlags().StringVar(&app.rootDir, "root-dir", ".", "The root source directory of the project")
	app.Command.PersistentFlags().StringVar(&app.snakeDir, "snake-dir", "build", "The Snake directory used for configuration")
	app.Command.PersistentFlags().BoolVar(&app.verbose, "verbose", false, "Enable verbose logging")

	app.Command.AddCommand(deployCmd, buildCmd, testCmd, configureCmd, installCmd, cleanCmd,
		packageCmd, runCmd, listProfilesCmd, listOptionsCmd, listTargetsCmd, docCmd, mutateCmd, formatCmd, generateCmd, newCmd)
}
