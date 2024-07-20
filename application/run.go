// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:                "run",
	Short:              "Run scripts or executable targets",
	DisableFlagParsing: true,
	RunE: func(c *cobra.Command, args []string) error {
		err := app.initFast()

		if err != nil {
			return err
		}

		if len(args) < 1 {
			return errors.New("you must specify a target to run")
		}

		if err = app.launch(filepath.Join(app.db.ProfilePath, "bin", args[0]), args[1:]...); err != nil {
			fmt.Println(err)
		}

		return nil
	},
}
