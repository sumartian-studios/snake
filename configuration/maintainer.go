// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type Maintainer struct {
	// Name of maintainer.
	Name string `yaml:"name"`

	// Contact details.
	Contact string `yaml:"contact"`

	// Role of maintainer.
	Role string `yaml:"role"`
}
