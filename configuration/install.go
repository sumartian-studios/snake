// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Install struct {
	// The type of installation rule (ex. FILE or DIRECTORY).
	Type string `yaml:"type" jsonschema:"required"`

	// List of arguments passed to the install() command.
	Rules []string `yaml:"rules" jsonschema:"required"`
}
