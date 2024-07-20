# Copyright (c) 2022 Sumartian Studios
#
# Snake is free software: you can redistribute it and/or modify it under the
# terms of the MIT license.

# Need to pass version string in ldflags...
build:
	go build -ldflags="-s -w -X 'github.com/sumartian/snake/application.VersionStr=1.0.0'"

# Need to generate the data.zip...
run:
	go run tools/schema-generator/main.go

# Need to pass data.zip manually via --archive...
configure:
	snake configure --profile default --archive ~/code/snake/distribution/data.zip --update
