// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"time"

	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate documentation",
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Documentation")

		// TODO: doxygen support...

		return nil
	},
}
