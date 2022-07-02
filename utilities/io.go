// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package utilities

import (
	"io"
	"os"
)

// SmartLink overwrites existing symlinks.
func SmartLink(src string, dest string) error {
	if _, err := os.Lstat(dest); err == nil {
		os.Remove(dest)
	}

	return os.Symlink(src, dest)
}

func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)

	if err != nil {
		return err
	}

	defer srcFile.Close()

	destFile, err := os.Create(dest)

	if err != nil {
		return err
	}

	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}
