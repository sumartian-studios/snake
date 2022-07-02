// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	lsTargetsFlag  bool
	lsProfilesFlag bool
	lsOptionsFlag  bool
)

func listTargets() error {
	if app.cfg.Targets == nil {
		return nil
	}

	targets := *app.cfg.Targets

	for _, target := range targets {
		fmt.Println("--", target.Name)
	}

	return nil
}

func listOptions() error {
	if app.cfg.Features == nil {
		return nil
	}

	feats := *app.cfg.Features

	for _, feat := range feats {
		if feat.Key != nil && feat.Description != nil && feat.Value != nil {
			fmt.Println("--", (*feat.Key)+"="+(*feat.Value), fmt.Sprintf("\033[0;90m%s\033[0m", *feat.Description))
		}
	}

	return nil
}

func listProfiles() error {
	var s string

	for i, profile := range app.cfg.Profiles {
		if i == app.db.ProfileIndex {
			s = "-- [x] " + profile.Name + " (current)"
		} else {
			s = "-- [ ] " + profile.Name
		}

		fmt.Println(s)
	}

	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available information",
	RunE: func(c *cobra.Command, args []string) error {
		if err := app.initSlow(); err != nil {
			return err
		}

		if lsTargetsFlag {
			return listTargets()
		} else if lsProfilesFlag {
			return listProfiles()
		} else if lsOptionsFlag {
			return listOptions()
		} else {
			c.Help()
		}

		return nil
	},
}

func init() {
	listCmd.PersistentFlags().BoolVarP(&lsTargetsFlag, "targets", "t", false, "List project build targets")
	listCmd.PersistentFlags().BoolVarP(&lsProfilesFlag, "profiles", "p", false, "List project build profiles")
	listCmd.PersistentFlags().BoolVarP(&lsOptionsFlag, "options", "o", false, "List project build options")
}
