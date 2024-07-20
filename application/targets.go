// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listTargets() error {
	if app.cfg.Targets == nil {
		fmt.Println("No targets available")
		return nil
	}

	targets := *app.cfg.Targets

	for _, target := range targets {
		fmt.Println("--", target.Name)
	}

	return nil
}

var listTargetsCmd = &cobra.Command{
	Use:   "targets",
	Short: "List available targets",
	RunE: func(c *cobra.Command, args []string) error {
		if err := app.initSlow(); err != nil {
			return err
		}

		return listTargets()
	},
}
