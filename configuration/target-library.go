// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type TargetLibrary struct {
	// The link type (private or public).
	Type string `yaml:"type"`

	// The list of targets.
	Targets []string `yaml:"targets"`
}
