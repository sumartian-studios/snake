// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"time"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Scaffold a new project or create files from templates",
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Creation")

		// TODO: ...

		return nil
	},
}
