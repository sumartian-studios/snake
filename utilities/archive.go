// Copyright (c) 2022 Sumartian Studios
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

	for _, file := range files {
		fmt.Println("compress:", file)

		f, err := os.Open(file)

		if err != nil {
			return err
		}

		defer f.Close()

		info, err := f.Stat()

		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)

		if err != nil {
			return err
		}

		header.Name = file
		header.Method = zip.Deflate

		writer, err := writer.CreateHeader(header)

		if err != nil {
			return err
		}

		if _, err = io.Copy(writer, f); err != nil {
			return err
		}
	}

	return nil
}

func Decompress(src string, dest string) error {
	r, err := zip.OpenReader(src)

	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		p := filepath.Join(dest, f.Name)

		if strings.Contains(p, "..") {
			return fmt.Errorf("destination path should not contain '..': %s", p)
		}

		fmt.Println("write:", p)

		if f.FileInfo().IsDir() {
			os.MkdirAll(p, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
			return err
		}

		out, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

		if err != nil {
			return err
		}

		rc, err := f.Open()

		if err != nil {
			return err
		}

		_, err = io.Copy(out, rc)

		out.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
