# Copyright (c) 2022-2024 Sumartian Studios
#
# Snake is free software: you can redistribute it and/or modify it under the
# terms of the MIT license.

PROJECT_VERSION = 2.0.0
OUTPUT_DIR      = ~/.local/bin
EXAMPLE_DIR     = ./examples

# Need to pass version string in ldflags...
build:
	echo "Building..."
	go build -ldflags="-s -w -X 'github.com/sumartian-studios/snake/application.VersionStr=${PROJECT_VERSION}'"  -o ${OUTPUT_DIR}

# Need to generate the data.zip...
generate-schema:
	go run tools/schema-generator/main.go

# Examples tests...
# ----------------------------------------------------------------------------
configure-example:
	make build
	snake --root-dir=${EXAMPLE_DIR} configure

clean-example:
	snake --root-dir=${EXAMPLE_DIR} clean
