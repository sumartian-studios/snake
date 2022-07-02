// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Feature struct {
	// Other properties will only be evaluated if this property evaluates to true.
	Condition *string `yaml:"if"`

	// List of CMake instructions.
	Scripts *[]string `yaml:"scripts"`

	// List of target definitions.
	Definitions *[]string `yaml:"defines"`

	// The variable identifier.
	Key *string `yaml:"key"`

	// Optional description. This is mandatory for user provided options.
	Description *string `yaml:"description"`

	// The value of the variable.
	Value *string `yaml:"value"`
}
