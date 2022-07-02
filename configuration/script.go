// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Script struct {
	// Script name.
	Name string `yaml:"name"`

	// Script description.
	Description string `yaml:"description"`

	// List of commands this script will run.
	Commands []string `yaml:"commands"`

	// Optional dependencies.
	Requires *[]string `yaml:"requires"`

	// Optional products (i.e. files this script will produce).
	Products *[]string `yaml:"products"`
}
