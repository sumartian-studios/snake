// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Test struct {
	// Name of this test.
	Name string `yaml:"name"`

	// List of function names.
	Functions []string `yaml:"functions"`
}
