// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package repl

import (
	"strings"
)

// Suggestion represents a command-level suggestion which can possess children.
type Suggestion struct {
	Prefix   string
	children []string
}

// Add a new suggestion string.
func (suggestion *Suggestion) AddChild(value string) {
	suggestion.children = append(suggestion.children, value)
}

// Suggester can perform auto-completion.
type Suggester struct {
	Suggestions []*Suggestion
}

// Create a new root suggestion item.
func (suggester *Suggester) Create(s string) *Suggestion {
	suggestion := new(Suggestion)
	suggestion.Prefix = s

	suggester.Suggestions = append(suggester.Suggestions, suggestion)

	return suggestion
}

// Get suggestions/completions.
func (suggester *Suggester) Do(line []rune, pos int) (suggestions [][]rune, offset int) {
	if len(line) == 0 {
		for _, suggestion := range suggester.Suggestions {
			suggestions = append(suggestions, []rune(suggestion.Prefix+" "))
		}

		return suggestions, offset
	}

	str := string(line)

	if i := strings.LastIndex(str, "&&"); i != -1 {
		str = str[i+3:]
	}

	tokens := strings.Split(strings.ToLower(str), " ")

	// This can happen when we complete on just space characters.
	if len(tokens) == 0 {
		return suggestions, offset
	}

	argument := len(tokens) > 1

	var suffix, prefix string
	prefix = tokens[0]

	if argument {
		suffix = tokens[len(tokens)-1]
		offset = len(suffix)
	} else {
		offset = len(str)
	}

	for _, suggestion := range suggester.Suggestions {
		if p := strings.ToLower(suggestion.Prefix); strings.HasPrefix(p, prefix) {
			if argument {
				for _, child := range suggestion.children {
					if s := strings.ToLower(child); strings.HasPrefix(s, suffix) {
						if child[len(child)-1] == '=' {
							suggestions = append(suggestions, []rune(child[len(suffix):]))
						} else {
							suggestions = append(suggestions, []rune(child[len(suffix):]+" "))
						}
					}
				}
			} else {
				suggestions = append(suggestions, []rune(suggestion.Prefix[len(prefix):]+" "))
			}
		}
	}

	return suggestions, offset
}
