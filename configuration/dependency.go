// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type DependencyImport struct {
	// The name of the import (ex. mylib::mylib).
	Name string `yaml:"target" jsonschema:"required"`

	// The find_package() argument string (ex. mylib REQUIRED).
	Declare string `yaml:"find" jsonschema:"required"`
}

type Dependency struct {
	// The package manager or source to use when fetching the package.
	From string `yaml:"from"`

	// The package string. This should be in the form of
	// NAME/(VERSION or TAG)
	Package string `yaml:"package"`

	// The URL or path to the resource.
	Path string `yaml:"path"`

	// The dependency imports.
	Imports []DependencyImport `yaml:"imports"`
}
