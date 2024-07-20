// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Configuration struct {
	// Project name.
	Project string `yaml:"Project"`

	// Short project description.
	Description string `yaml:"Description"`

	// Version (ex. 0.0.1).
	Version string `yaml:"Version"`

	// Organization name.
	Organization string `yaml:"Organization"`

	// Organization contact information.
	Contact string `yaml:"Contact"`

	// List of maintainers.
	Maintainers []Maintainer `yaml:"Maintainers"`

	// List of root-level features.
	Features *[]Feature `yaml:"Features"`

	// Project website.
	Site string `yaml:"Site"`

	// Project repository.
	Repository string `yaml:"Repository"`

	// Local path to project logo image.
	Logo string `yaml:"Logo"`

	// Project license.
	License string `yaml:"License"`

	// List o runnable scripts.
	Scripts *[]Script `yaml:"Scripts"`

	// List of build profiles.
	Profiles []Profile `yaml:"Profiles"`

	// List of third-party Conan packages.
	Dependencies *[]Dependency `yaml:"Dependencies"`

	// List of targets.
	Targets *[]Target `yaml:"Targets"`
}
