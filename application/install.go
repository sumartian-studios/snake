// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"time"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install targets (i.e. cmake install)",
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Installation")

		if err := app.initSlow(); err != nil {
			return err
		}

		if err := app.launch("cmake", "--build", app.db.ProfilePath, "--", "install"); err != nil {
			return err
		}

		return nil
	},
}
