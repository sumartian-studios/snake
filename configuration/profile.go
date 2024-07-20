// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Profile struct {
	// Profile name.
	Name string `yaml:"id"`

	// Profile description.
	Description string `yaml:"description"`

	// Build type (ex. Debug or Release).
	Type string `yaml:"type"`

	// Optional operating system (ex. Windows, Linux, macOS).
	System string `yaml:"system"`

	// The C++ compiler.
	Compiler string `yaml:"compiler"`

	// Optional system architecture.
	Arch string `yaml:"arch"`

	// List of option maps.
	Variables []map[string]string `yaml:"options"`

	// List of linker flags.
	LinkFlags []string `yaml:"flags.link"`

	// List of compiler flags.
	CompileFlags []string `yaml:"flags.compile"`
}
