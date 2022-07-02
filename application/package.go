// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"time"

	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the targets for deployment",
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Packaging")

		if err := app.initSlow(); err != nil {
			return err
		}

		if err := app.launch("cmake", "--build", app.db.ProfilePath, "--", "package"); err != nil {
			return err
		}

		return nil
	},
}
