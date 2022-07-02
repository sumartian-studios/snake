// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"errors"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:                "build",
	Short:              "Build a target using the build system or run a script",
	DisableFlagParsing: true,
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Build")

		err := app.initFast()

		if err != nil {
			return err
		}

		cmd := exec.Command("cmake", append([]string{"--build", app.db.ProfilePath, "--"}, args...)...)
		cmd.Env = os.Environ()

		if app.db.ProfileIndex == -1 {
			return errors.New("you must re-configure this project (snake configure)")
		}

		cmd.Stderr, cmd.Stdout, cmd.Stdin = os.Stderr, os.Stdout, os.Stdin

		return cmd.Run()
	},
}
