// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listOptions() error {
	if app.cfg.Features == nil {
		fmt.Println("No options available")
		return nil
	}

	feats := *app.cfg.Features

	for _, feat := range feats {
		var key, description, value = "?", "?", "?"

		if feat.Key != nil {
			key = *feat.Key
			fmt.Println("Warning: key is empty")
		}

		if feat.Description != nil {
			description = *feat.Description
			fmt.Println("Warning: description is empty")
		}

		if feat.Value != nil {
			value = *feat.Value
			fmt.Println("Warning: value is empty")
		}

		fmt.Println("--", key+"="+value, fmt.Sprintf("\033[0;90m%s\033[0m", description))
	}

	return nil
}

var listOptionsCmd = &cobra.Command{
	Use:   "options",
	Short: "List available options",
	RunE: func(c *cobra.Command, args []string) error {
		if err := app.initSlow(); err != nil {
			return err
		}

		return listOptions()
	},
}
