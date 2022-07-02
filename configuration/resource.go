// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Resource struct {
	// Optional prefix that will be prepended to the resource alias.
	Prefix *string `yaml:"prefix"`

	// Optional module name if this resource group describes a QML module.
	Module *string `yaml:"module"`

	// List of files used by this resource group. You may use
	// regular expression.
	Files []string `yaml:"files"`
}
