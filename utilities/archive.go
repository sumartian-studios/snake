// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package utilities

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Compress(filename string, files []string) error {
	zipFile, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	compressAndWriteFile := func(f *os.File) error {
		info, err := f.Stat()

		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)

		if err != nil {
			return err
		}

		header.Name = f.Name()
		header.Method = zip.Deflate

		writer, err := writer.CreateHeader(header)

		if err != nil {
			return err
		}

		defer f.Close()

		if _, err = io.Copy(writer, f); err != nil {
			return err
		}

		return nil
	}

	for _, file := range files {
		fmt.Println("compress:", file)

		f, err := os.Open(file)

		if err != nil {
			return err
		}

		if err = compressAndWriteFile(f); err != nil {
			return err
		}
	}

	return nil
}

func Decompress(reader *zip.Reader, dest string) error {
	dest = filepath.Clean(dest) + string(os.PathSeparator)

	// Make the destination directory if it does not exist.
	os.MkdirAll(dest, os.ModePerm)

	extractAndWriteFile := func(f *zip.File) error {
		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip: https://snyk.io/research/zip-slip-vulnerability
		// See: https://stackoverflow.com/a/58192644
		if !strings.HasPrefix(path, dest) {
			return fmt.Errorf("%s: illegal file path", path)
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		defer rc.Close()

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}

			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, f := range reader.File {
		err := extractAndWriteFile(f)

		if err != nil {
			return err
		}
	}

	return nil
}
