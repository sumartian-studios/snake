# Copyright (c) 2022 Sumartian Studios
#
# Snake is free software: you can redistribute it and/or modify it under the
# terms of the MIT license.

PROJECT_VERSION = 1.0.0

# Need to pass version string in ldflags...
build:
	go build -ldflags="-s -w -X 'github.com/sumartian/snake/application.VersionStr=${PROJECT_VERSION}'"

# Need to generate the data.zip...
run:
	go run tools/schema-generator/main.go

# Need to pass data.zip manually via --archive...
configure:
	snake configure --profile default --archive ./distribution/data.zip --update
