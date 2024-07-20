// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the exported packages",
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Deployment")

		fmt.Println("Deploying...")

		// TODO...

		return nil
	},
}
