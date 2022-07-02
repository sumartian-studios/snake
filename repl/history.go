// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package repl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/tchap/go-patricia/v2/patricia"
)

// HistoryItem represents a history entry.
type HistoryItem struct {
	// The text displayed by the hint.
	Text string
	// The history time. Used as priority when sorting hints.
	Time int
}

type HistoryTrie struct {
	*patricia.Trie

	// Path to the history file.
	Path string

	// List of hints.
	Hints []HistoryItem

	// Used for navigation (Ctrl+N, Ctrl+P).
	NavigationArray []HistoryItem

	// Return true if navigation has started.
	Navigating bool

	// Current index of visible hint.
	NavigationIndex int

	// The offset at which the user provided prefix ends. Must
	// be less than the length of the visible hint text.
	CurrentHintTextOffset int

	// The visible hint text.
	CurrentHintText string
}

// Load history optionally merging it with the existing trie.
func (h *HistoryTrie) Load(merge bool) error {
	if _, err := os.Stat(h.Path); os.IsNotExist(err) {
		return nil
	}

	file, err := ioutil.ReadFile(h.Path)

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(file))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		scanner.Scan()

		var timestamp int

		if timestamp, err = strconv.Atoi(scanner.Text()); err != nil {
			return err
		}

		prefix := patricia.Prefix(text)

		if merge {
			if existingTime, ok := h.Get(prefix).(int); ok && timestamp > existingTime {
				h.Set(prefix, timestamp)
			}
		} else {
			h.Set(prefix, timestamp)
		}
	}

	return nil
}

// Append a line to history.
func (h *HistoryTrie) Push(prefix patricia.Prefix) {
	h.Set(prefix, int(time.Now().Unix()))
}

// Go to the next history entry and return selected text.
func (h *HistoryTrie) GoToNextEntry() string {
	if h.NavigationIndex > 0 {
		h.NavigationIndex--
	}

	if !h.Navigating {
		h.StartNavigation()
	}

	return h.NavigationArray[h.NavigationIndex].Text
}

// Go to the previous history entry and return selected text.
func (h *HistoryTrie) GoToPrevEntry() string {
	if h.NavigationIndex < len(h.NavigationArray)-1 {
		h.NavigationIndex++
	}

	if !h.Navigating {
		h.StartNavigation()
	}

	return h.NavigationArray[h.NavigationIndex].Text
}

// Start history navigation.
func (h *HistoryTrie) StartNavigation() {
	h.Visit(func(prefix patricia.Prefix, item patricia.Item) error {
		timestamp, ok := item.(int)

		if !ok {
			return errors.New("unable to read time")
		}

		h.NavigationArray = append(h.NavigationArray, HistoryItem{Text: string(prefix), Time: timestamp})

		return nil
	})

	if len(h.NavigationArray) > 0 {
		// Sort entries by recency
		sort.Slice(h.NavigationArray, func(i, j int) bool {
			return h.NavigationArray[i].Time > h.NavigationArray[j].Time
		})
	}

	h.Navigating = true
}

// End history navigation.
func (h *HistoryTrie) EndNavigation() {
	h.NavigationIndex = 0
	h.NavigationArray = []HistoryItem{}
	h.Navigating = false
}

// Save history.
func (h *HistoryTrie) Save() {
	var b bytes.Buffer

	// Another process could have changed the history so we first
	// merge existing values before saving.
	if err := h.Load(true); err != nil {
		panic(err)
	}

	h.Visit(func(prefix patricia.Prefix, item patricia.Item) error {
		timestamp, ok := item.(int)

		if !ok {
			return errors.New("unable to save timestamp")
		}

		b.WriteString(fmt.Sprintf("%s\n%d\n", prefix, timestamp))

		return nil
	})

	if err := ioutil.WriteFile(h.Path, b.Bytes(), 0644); err != nil {
		panic(err)
	}
}
