// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"time"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:                "test",
	Short:              "Run built unit tests and benchmarks",
	DisableFlagParsing: true,
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Testing")

		if err := app.initFast(); err != nil {
			return err
		}

		if len(args) == 0 {
			args = append(args, ".*")
		} else if args[0][0] == '-' {
			args = append([]string{".*"}, args...)
		}

		opts := append([]string{"--test-dir", app.db.ProfilePath,
			"--output-on-failure", "-R"}, args...)

		if err := app.launch("ctest", opts...); err != nil {
			return err
		}

		return nil
	},
}
