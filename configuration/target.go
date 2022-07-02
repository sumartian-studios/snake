// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Target struct {
	// The name used to identify the target.
	Name string `yaml:"name" jsonschema:"required"`

	// A short description of the target. You can access this variable in C/C++ by
	// including "version.h" and accessing EXE_TARGET_DESCRIPTION.
	Description string `yaml:"description" jsonschema:"required"`

	// The type of target. Type can be "executable", "application", "shared-library",
	// "static-library", "header-library", "plugin", or "test".
	Type string `yaml:"type" jsonschema:"required"`

	// Set this property to true if you want to export a library under the
	// project namespace. This will automatically handle installation.
	Export *bool `yaml:"export"`

	// The target (and its unique dependencies) are only enabled when this condition
	// evaluates to true. Set the value to 'SNAKE_ALWAYS_BUILD' to have no requirements; do
	// not set it to an empty value as that would evaluate to false.
	Requirement string `yaml:"requirement" jsonschema:"required"`

	// The *relative* path to the source and header files. The TARGET_SOURCE_DIR variable
	// is set to the absolute path and can be used by your configuration.
	Path string `yaml:"path" jsonschema:"required"`

	// List of conditional features acquired by this target. These
	// are executed in the order they are provided.
	Features *[]TargetFeature `yaml:"features"`
}
