// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"archive/zip"
	"bytes"
	"fmt"

	"github.com/sumartian/snake/utilities"
)

// Decompress the embedded zip file.
func (app *Application) decompress() error {
	file, err := app.dataZip.ReadFile("distribution/data.zip")

	reader := bytes.NewReader(file)
	zipReader, err := zip.NewReader(reader, int64(len(file)))

	if err != nil {
		return err
	}

	if err = utilities.Decompress(zipReader, app.snakeDir); err != nil {
		return fmt.Errorf("unable to decompress embedded zip: %v", err)
	}

	return nil
}
