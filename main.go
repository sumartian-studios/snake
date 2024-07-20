// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/sumartian-studios/snake/application"
)

//go:embed distribution/data.zip
var dataZip embed.FS

func main() {
	if err := application.Execute(&dataZip); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
