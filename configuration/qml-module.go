// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type QMLModule struct {
	// QML module URI.
	Uri string `yaml:"uri" jsonschema:"required"`

	// QML module version.
	Version *string `yaml:"version"`

	// QML module prefix.
	Prefix *string `yaml:"prefix"`
}
