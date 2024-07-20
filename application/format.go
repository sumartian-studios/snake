// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"time"

	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format all accessible source code in project",
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Formatting")

		if err := app.initSlow(); err != nil {
			return err
		}

		if err := app.launch("cmake", "--build", app.db.ProfilePath, "--", "format"); err != nil {
			return err
		}

		return nil
	},
}
