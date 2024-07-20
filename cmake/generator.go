// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package cmake

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sumartian-studios/snake/configuration"
)

type PreDependency struct {
	FindPackageString string
	Dependency        *configuration.Dependency
}

// Generator is a CMake generator.
type Generator struct {
	// Used to prepend data.
	Start bytes.Buffer
	End   bytes.Buffer

	// Used to append data.
	Buffer *bytes.Buffer

	// The generator context.
	Context struct {
		// The current state of automoc.
		enableAutoMoc bool

		// The current default link type.
		defaultLinkType string

		// Map of dependencies.
		LibraryMap map[string]PreDependency

		// Map of requirements.
		RequirementMap map[string]map[string]bool

		// Map of conditional aliases (ex. and becomes AND)
		IfAliasMap map[string]bool
	}
}

// Call a CMake macro/function.
func (g *Generator) Call(name string, args ...string) {
	g.Buffer.WriteString(name)
	g.Buffer.WriteString("(")

	length := len(args) - 1

	for i, arg := range args {
		g.Buffer.WriteString(arg)

		if i < length {
			g.Buffer.WriteString(" ")
		}
	}

	g.Buffer.WriteString(")\n")
}

func (g *Generator) CleanConditional(s string) string {
	fields := strings.Fields(s)

	for i, f := range fields {
		if g.Context.IfAliasMap[f] {
			fields[i] = strings.ToUpper(f)
		}
	}

	return strings.Join(fields, " ")
}

func (g *Generator) ExistsOr(s *string, defaultValue string) string {
	if s != nil {
		return *s
	} else {
		return defaultValue
	}
}

// Save the buffer to a file.
func (g *Generator) Save(path string) error {
	g.Start.Write(g.End.Bytes())
	return ioutil.WriteFile(path, g.Start.Bytes(), 0644)
}

func (g *Generator) LinkLibrary(t *configuration.Target, lib string) {
	g.Call("target_link_libraries", t.Name, g.Context.defaultLinkType, lib)
}

func (g *Generator) AddTargetLibrary(t *configuration.Target, feat *configuration.TargetFeature, lib string) {
	before, _, found := strings.Cut(lib, "::")

	if found && strings.HasPrefix(before, "Qt") {
		if !g.Context.enableAutoMoc && g.Context.defaultLinkType != "INTERFACE" {
			g.Call("set_target_properties", t.Name, "PROPERTIES", "AUTOMOC", "on")
		}

		g.Context.enableAutoMoc = true
	}

	if l := g.Context.LibraryMap[lib]; l.Dependency != nil {
		if g.Context.RequirementMap[lib] == nil {
			g.Context.RequirementMap[lib] = map[string]bool{}
		}

		g.Context.RequirementMap[lib][t.Requirement] = true

		if feat.Condition != nil {
			g.Context.RequirementMap[lib][*feat.Condition] = true
		}

		g.Call("find_package", l.FindPackageString)
	}

	g.LinkLibrary(t, lib)
}

func (g *Generator) AddTargetFeature(t *configuration.Target, feat *configuration.TargetFeature) {
	if feat.Condition != nil {
		g.Call("if", g.CleanConditional(*feat.Condition))
	}

	if feat.Libraries != nil {
		libraries := *feat.Libraries

		for _, lib := range libraries {
			if t.Type != "header-library" {
				if a := strings.ToLower(lib.Type); a == "private" {
					g.Context.defaultLinkType = "PRIVATE"
				} else {
					g.Context.defaultLinkType = "PUBLIC"
				}
			} else {
				g.Context.defaultLinkType = "INTERFACE"
			}

			for _, lname := range lib.Targets {
				g.AddTargetLibrary(t, feat, lname)
			}
		}
	}

	if feat.Module != nil {
		m := *feat.Module
		g.Call("snake_add_qml_module",
			Quote(t.Name),
			Quote(m.Uri), Quote(g.ExistsOr(m.Version, "")),
			Quote(g.ExistsOr(m.Prefix, "")))
	}

	if feat.Properties != nil {
		properties := *feat.Properties
		for _, group := range properties {
			for k, v := range group {
				g.Call("set_target_properties", t.Name, "PROPERTIES", k, v)
			}
		}
	}

	if feat.Definitions != nil {
		definitions := *feat.Definitions
		for _, d := range definitions {
			g.Call("target_compile_definitions", t.Name, "PUBLIC", d)
		}
	}

	if feat.Scripts != nil {
		g.Buffer.WriteString(strings.Join(*feat.Scripts, "\n"))
		g.Buffer.WriteString("\n")
	}

	if feat.Plugins != nil {
		plugins := *feat.Plugins
		for _, plugin := range plugins {
			g.Call("snake_import_plugin", t.Name, plugin)
		}
	}

	if feat.Tests != nil {
		tests := *feat.Tests
		for _, test := range tests {
			g.Call("add_test", "NAME", Quote(test.Name), "COMMAND", t.Name+" "+strings.Join(test.Functions, " "))
		}
	}

	if feat.Resources != nil {
		resources := *feat.Resources
		for _, resource := range resources {
			var prefix, module string

			if resource.Prefix != nil {
				prefix = *resource.Prefix
			}

			if resource.Module != nil {
				module = *resource.Module
			}

			g.Call("snake_add_resources", Quote(t.Name),
				Quote(strings.Join(resource.Files, ";")), Quote(module), Quote(prefix))
		}
	}

	if feat.Installs != nil {
		installs := *feat.Installs
		for _, install := range installs {
			g.Call("install", install.Type+" "+strings.Join(install.Rules, " "))
		}
	}

	if feat.Condition != nil {
		g.Call("endif")
	}
}

func (g *Generator) AddTarget(t *configuration.Target, i int, count int) {
	g.Call("set", "TARGET_STATUS", fmt.Sprintf("\"[%02d/%02d] %s\"", i+1, count, t.Name))
	g.Call("if", g.CleanConditional(t.Requirement))
	g.Call("print_status", "\"${TARGET_STATUS}\"")

	g.Context.enableAutoMoc = false

	// Because interface libraries can only use interface sources...
	// Otherwise just use public includes.
	g.Context.defaultLinkType = "PRIVATE"

	switch t.Type {
	case "executable":
		g.Call("add_executable", t.Name)
	case "application":
		g.Call("snake_create_graphical_app", t.Name)
	case "static-library":
		g.Call("add_library", t.Name, "STATIC")
	case "shared-library":
		g.Context.defaultLinkType = "PUBLIC"
		g.Call("add_library", t.Name, "${SNAKE_LIB_TYPE}")
	case "header-library":
		g.Context.defaultLinkType = "INTERFACE"
		g.Call("add_library", t.Name, "INTERFACE")
	case "test":
		g.Call("add_executable", t.Name)
	case "plugin":
		g.Call("add_library", t.Name, "MODULE")
	}

	export := "off"

	if t.Export != nil && *t.Export {
		export = "on"
	}

	g.Call("snake_init_target", t.Name, Quote(t.Path),
		g.Context.defaultLinkType, t.Type, Quote(t.Description), export)

	if t.Features != nil {
		features := *t.Features

		for _, feat := range features {
			g.AddTargetFeature(t, &feat)
		}
	}

	g.Call("snake_fini_target", t.Name)
	g.Call("else")
	g.Call("print_dim_status", "\"${TARGET_STATUS} (disabled)\"")
	g.Call("endif")
}

func (g *Generator) AddScript(s *configuration.Script) {
	if len(s.Commands) > 0 {
		a := []string{s.Name, "WORKING_DIRECTORY", "${CMAKE_SOURCE_DIR}"}

		if s.Products != nil {
			a = append(a, "BYPRODUCTS", Quote(strings.Join(*s.Products, " ")))
		}

		if s.Requires != nil {
			a = append(a, "DEPENDS", Quote(strings.Join(*s.Requires, " ")))
		}

		for _, exec := range s.Commands {
			if len(exec) > 0 {
				a = append(a, "COMMAND")
				a = append(a, exec)
			}
		}

		g.Call("add_custom_target", a...)
	}
}

func (g *Generator) AddGlobalFeature(feat *configuration.Feature) {
	if feat.Condition != nil {
		g.Call("if", g.CleanConditional(*feat.Condition))
	}

	if feat.Key != nil && feat.Value != nil {
		if feat.Description != nil {
			g.Call("set", *feat.Key, *feat.Value, "CACHE", "INTERNAL", Quote(""))
		} else {
			g.Call("set", *feat.Key, *feat.Value)
		}
	}

	if feat.Scripts != nil {
		g.Buffer.WriteString(strings.Join(*feat.Scripts, "\n"))
		g.Buffer.WriteString("\n")
	}

	if feat.Definitions != nil {
		g.Call("add_compile_definitions", *feat.Definitions...)
	}

	if feat.Condition != nil {
		g.Call("endif")
	}
}
