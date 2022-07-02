// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package configuration

type TargetFeature struct {
	Feature `yaml:",inline"`

	// List of libraries used by the target. If the name is separated by '::' then
	// the left side will be used as the package name for find_package. Qt6 components
	// are automatically handled.
	Libraries *[]TargetLibrary `yaml:"libraries"`

	// List of installation rules used by the target.
	Installs *[]Install `yaml:"installs"`

	// List of resources that will be embedded into the target.
	Resources *[]Resource `yaml:"resources"`

	// List of imported plugins.
	Plugins *[]string `yaml:"plugins"`

	// Map of target properties.
	Properties *[]map[string]string `yaml:"properties"`

	// List of enabled tests. Only active when target type is "test".
	Tests *[]Test `yaml:"tests"`

	// Qt QML module definition.
	Module *QMLModule `yaml:"module"`
}
