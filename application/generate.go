// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/sumartian-studios/snake/cmake"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Re-generate the CMakeLists.txt",
	RunE: func(c *cobra.Command, args []string) error {
		defer app.timeTrack(time.Now(), "Generation")

		err := app.initSlow()

		if err != nil {
			return err
		}

		cmakeListsTxt := filepath.Join(app.rootDir, "CMakeLists.txt")

		fmt.Println("Generating...", cmakeListsTxt)

		g := new(cmake.Generator)

		g.Context.IfAliasMap = map[string]bool{
			"and":      true,
			"exists":   true,
			"or":       true,
			"not":      true,
			"strequal": true,
			"less":     true,
			"greater":  true,
			"matches":  true,
		}

		g.Buffer = &g.Start

		g.Buffer.WriteString(fmt.Sprintf("# Generated by Snake (%s). You must not modify this file.\n\n", app.Version))

		g.Call("if", "NOT", "DEFINED", "SNAKE_DIR")
		g.Call("message", "STATUS", cmake.Quote("Snake directory is not defined..."))
		g.Call("if", "DEFINED", "NO_SNAKE")
		g.Call("set", "SNAKE_DIR", cmake.Quote("${CMAKE_BINARY_DIR}"), "CACHE", "INTERNAL", cmake.Quote(""))
		g.Call("message", "STATUS", cmake.Quote("Not using Snake... ${SNAKE_DIR}"))
		g.Call("else")
		g.Call("message", "FATAL_ERROR", cmake.Quote("You must re-configure the project using Snake or set NO_SNAKE=on"))
		g.Call("endif")
		g.Call("else")
		g.Call("message", "STATUS", cmake.Quote("Slithering into... ${SNAKE_DIR}"))
		g.Call("endif")

		g.Call("cmake_minimum_required", "VERSION", "3.30.0", "FATAL_ERROR")
		g.Call("project", app.cfg.Project, "VERSION", app.cfg.Version, "LANGUAGES", "CXX")

		g.Call("set", "SNAKE_CONTACT", cmake.Quote(app.cfg.Contact))
		g.Call("set", "SNAKE_ORGANIZATION", cmake.Quote(app.cfg.Organization))
		g.Call("set", "SNAKE_PROJECT_LICENSE", cmake.Quote(app.cfg.License))
		g.Call("set", "SNAKE_PROJECT_REPOSITORY", cmake.Quote(app.cfg.Repository))

		g.Call("set", "CMAKE_PROJECT_HOMEPAGE_URL", cmake.Quote(app.cfg.Site))
		g.Call("set", "CMAKE_PROJECT_DESCRIPTION", cmake.Quote(app.cfg.Description))

		g.Call("include", cmake.Quote("${SNAKE_DIR}/snake.1.cmake"))
		g.Call("include", cmake.Quote("${SNAKE_DIR}/snake.2.cmake"))

		g.Context.LibraryMap = map[string]cmake.PreDependency{}
		g.Context.RequirementMap = map[string]map[string]bool{}

		if app.cfg.Dependencies != nil {
			dependencies := *app.cfg.Dependencies

			for i, d := range dependencies {
				before, _, _ := strings.Cut(d.Package, "/")

				if len(before) < 1 {
					return fmt.Errorf("package name cannot be empty: %s", d.Package)
				}

				for _, imports := range d.Imports {
					g.Context.LibraryMap[imports.Name] = cmake.PreDependency{
						FindPackageString: imports.Declare,
						Dependency:        &dependencies[i],
					}
				}
			}
		}

		if app.cfg.Features != nil {
			feats := *app.cfg.Features

			for _, feat := range feats {
				g.AddGlobalFeature(&feat)
			}
		}

		// The end buffer starts here. Append to g.Start to preprend to g.End.
		g.Buffer = &g.End

		g.Call("include", cmake.Quote("${SNAKE_DIR}/snake.3.cmake"))

		if app.cfg.Targets != nil {
			targets := *app.cfg.Targets
			count := len(targets)

			for i, t := range targets {
				g.AddTarget(&t, i, count)
			}
		}

		g.Buffer = &g.Start

		for k, v := range g.Context.RequirementMap {
			requirements := []string{}

			for kk := range v {
				requirements = append(requirements, "("+kk+")")
			}

			if l := g.Context.LibraryMap[k]; l.Dependency != nil {

				if l.Dependency.From == "pkg" {
					g.Call("if", g.CleanConditional(strings.Join(requirements, " AND ")))
					g.Call("snake_fetch_pkg", cmake.Quote(l.Dependency.Package))
					g.Call("endif")
				} else if l.Dependency.From == "conan" {
					g.Call("if", g.CleanConditional(strings.Join(requirements, " AND ")))
					g.Call("list", "APPEND", "ENABLED_CONAN_PACKAGES", cmake.Quote(l.Dependency.Package))
					g.Call("endif")
				} else if l.Dependency.From == "url" || l.Dependency.From == "git" {
					before, after, found := strings.Cut(l.Dependency.Package, "/")

					if !found {
						return fmt.Errorf("invalid arguments: %s", l.Dependency.Package)
					}

					if l.Dependency.From == "url" {
						g.Call("if", g.CleanConditional(strings.Join(requirements, " AND ")))
						g.Call("snake_fetch_url", cmake.Quote(before),
							cmake.Quote(l.Dependency.Path), cmake.Quote(after))
						g.Call("endif")
					} else {
						g.Call("if", g.CleanConditional(strings.Join(requirements, " AND ")))
						g.Call("snake_fetch_git", cmake.Quote(before),
							cmake.Quote(l.Dependency.Path), cmake.Quote(after))
						g.Call("endif")
					}
				} else if l.Dependency.From == "system" {
					// Nothing to do...
				} else {
					return fmt.Errorf("unsupported dependency provider: %s", l.Dependency.From)
				}

			}
		}

		g.Buffer = &g.End

		// Scripts
		if app.cfg.Scripts != nil {
			scripts := *app.cfg.Scripts
			for _, s := range scripts {
				g.AddScript(&s)
			}
		}

		g.Call("include", cmake.Quote("${SNAKE_DIR}/snake.4.cmake"))

		if err = g.Save(cmakeListsTxt); err != nil {
			return err
		}

		return nil
	},
}
