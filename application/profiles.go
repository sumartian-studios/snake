// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sumartian-studios/snake/configuration"
)

func listProfiles() error {
	var s string

	if len(app.cfg.Profiles) < 1 {
		fmt.Println("No profiles available")
		return nil
	}

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

var listProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List available profiles",
	RunE: func(c *cobra.Command, args []string) error {
		if err := app.initSlow(); err != nil {
			return err
		}

		return listProfiles()
	},
}

// Returns the current profile if it exists. If it does not exist it returns false and nil.
func (app *Application) getCurrentProfile() (bool, *configuration.Profile) {
	if len(app.db.ProfilePath) > 0 && len(app.cfg.Profiles) > app.db.ProfileIndex {
		return true, &app.cfg.Profiles[app.db.ProfileIndex]
	}

	return false, nil
}

// Sets the current profile at index.
func (app *Application) setCurrentProfile(i int, p *configuration.Profile) {
	app.db.ProfileIndex = i
	app.db.ProfilePath = filepath.Join(app.snakeDir, p.Name)
	app.storageChanged()
}
