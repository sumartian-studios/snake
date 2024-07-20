// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package cmake

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"unicode"
)

// Reduces the size of a CMake file by removing whitespace and comments.
func Minify(path string) (data []byte, err error) {
	file, err := os.Open(path)

	if err != nil {
		return data, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var buffer bytes.Buffer

	multiline := false

	for scanner.Scan() {
		withoutLineComment, _, _ := strings.Cut(scanner.Text(), "#")

		s := strings.TrimFunc(withoutLineComment, func(r rune) bool {
			return unicode.IsSpace(r)
		})

		if len(s) == 0 {
			continue
		}

		if firstLetter := s[0]; firstLetter == '#' {
			continue
		}

		if lastLetter := s[len(s)-1]; lastLetter != ')' {
			if multiline {
				buffer.WriteString(" ")
			}
			buffer.WriteString(s)
			multiline = true
		} else {
			multiline = false
			buffer.WriteString(s)
			buffer.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return data, err
	}

	return buffer.Bytes(), err
}
