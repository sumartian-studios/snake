// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var cleanEverythingFlag bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Reset the current profile directory",
	RunE: func(c *cobra.Command, args []string) error {
		if err := app.initSlow(); err != nil {
			return err
		}

		deleteProfileBuildDir := func(s string) error {
			return os.RemoveAll(s)
		}

		if cleanEverythingFlag {
			for _, profile := range app.cfg.Profiles {
				if err := deleteProfileBuildDir(filepath.Join(app.snakeDir, profile.Name)); err != nil {
					return err
				}
			}
		} else {
			if err := deleteProfileBuildDir(app.db.ProfilePath); err != nil {
				return err
			}
		}

		if err := os.RemoveAll(filepath.Join(app.snakeDir, "snake.lock")); err != nil {
			return err
		}

		app.db.Configured = false
		app.db.ProfilePath = ""
		app.db.ProfileIndex = -1

		app.storageChanged()

		if err := app.saveStorage(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cleanCmd.PersistentFlags().BoolVarP(&cleanEverythingFlag, "all", "a", false, "Reset all profile build directories")
}
