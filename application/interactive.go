// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/sumartian/snake/repl"
)

func (app *Application) resetFlags() {
	for _, c := range app.Commands() {
		c.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				f.Value.Set(f.DefValue)
				f.Changed = false
			}
		})
	}
}

func startInteractiveMode(cmd *cobra.Command, args []string) error {
	if err := app.initSlow(); err != nil {
		return err
	}

	app.interactive = true

	r := repl.NewRepl()
	r.History.Path = filepath.Join(app.snakeDir, "snake.history.txt")

	r.Runner = func(args []string) error {
		// Need to reset the flags after each call otherwise
		// our flags will be stuck.
		app.resetFlags()
		app.SetArgs(args)
		return app.Execute()
	}

	for _, cmd := range app.Commands() {
		suggestion := r.Suggester.Create(cmd.Name())

		addFlags := func(f *pflag.Flag) {
			if len(f.Name) > 0 {
				suggestion.AddChild("--" + f.Name)
			}

			if len(f.Shorthand) > 0 {
				suggestion.AddChild("-" + f.Shorthand)
			}
		}

		if !cmd.DisableFlagParsing {
			cmd.LocalFlags().VisitAll(addFlags)
			cmd.InheritedFlags().VisitAll(addFlags)
		}

		switch cmd.Name() {
		case "configure":
			for _, profile := range app.cfg.Profiles {
				suggestion.AddChild("-p " + profile.Name)
			}

			if app.cfg.Features != nil {
				feats := *app.cfg.Features

				for _, feat := range feats {
					if feat.Key != nil && feat.Description != nil {
						suggestion.AddChild(*feat.Key + "=")
					}
				}
			}
		case "build":
			if app.cfg.Targets != nil {
				targets := *app.cfg.Targets
				for _, target := range targets {
					if target.Type != "header-library" {
						suggestion.AddChild(target.Name)
					}
				}
			}

			if app.cfg.Scripts != nil {
				scripts := *app.cfg.Scripts
				for _, script := range scripts {
					suggestion.AddChild(script.Name)
				}
			}
		case "run":
			if app.cfg.Targets != nil {
				targets := *app.cfg.Targets
				for _, target := range targets {
					if target.Type == "executable" || target.Type == "application" {
						suggestion.AddChild(target.Name)
					}
				}
			}
		case "test":
			if app.cfg.Targets != nil {
				targets := *app.cfg.Targets

				for _, target := range targets {
					if target.Type == "test" {

						if target.Features != nil {
							features := *target.Features

							for _, feat := range features {
								if feat.Tests != nil {
									tests := *feat.Tests

									for _, test := range tests {
										suggestion.AddChild(test.Name)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return r.Start()
}
