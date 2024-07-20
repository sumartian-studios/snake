// Copyright (c) 2022-2024 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package repl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strings"
	"unicode"

	"github.com/chzyer/readline"
	"github.com/tchap/go-patricia/v2/patricia"
)

// Read-print-loop manager.
type Repl struct {
	rl        *readline.Instance
	History   HistoryTrie
	Runner    func(args []string) error
	Suggester Suggester
}

// User key input handler.
func (r *Repl) HandleKey(line []rune, pos int, key rune) ([]rune, int, bool) {
	if key == readline.CharPrev {
		t := []rune(r.History.GoToPrevEntry())
		return t, len(t), true
	} else if key == readline.CharNext {
		t := []rune(r.History.GoToNextEntry())
		return t, len(t), true
	} else {
		r.History.EndNavigation()
	}

	if key == readline.CharBckSearch || key == readline.CharFwdSearch {
		return nil, 0, true
	}

	if key == readline.CharLineEnd && len(r.History.CurrentHintText) > 0 {
		line = append(line, []rune(r.History.CurrentHintText)...)
		r.History.CurrentHintText = ""
		r.History.CurrentHintTextOffset = 0
		return line, len(line), true
	}

	if !unicode.IsControl(key) {
		r.History.CurrentHintTextOffset = pos
		r.History.Hints = []HistoryItem{}

		err := r.History.VisitSubtree(patricia.Prefix(string(line)), func(prefix patricia.Prefix, item patricia.Item) error {
			timestamp, ok := item.(int)

			if !ok {
				return errors.New("unable to read time")
			}

			r.History.Hints = append(r.History.Hints, HistoryItem{
				Text: string(prefix[r.History.CurrentHintTextOffset:]),
				Time: timestamp,
			})

			return nil
		})

		if err != nil {
			panic(err)
		}

		if len(r.History.Hints) > 0 {
			// Sort entries by recency
			sort.Slice(r.History.Hints, func(i, j int) bool {
				return r.History.Hints[i].Time > r.History.Hints[j].Time
			})

			r.History.CurrentHintText = r.History.Hints[0].Text

			r.rl.Terminal.Print(fmt.Sprintf("\033[0;90m%s\033[0m%s",
				r.History.CurrentHintText,
				bytes.Repeat([]byte{'\b'}, len(r.History.CurrentHintText))))
		}
	}

	return nil, 0, false
}

// Perform cleanup tasks before quitting.
func (r *Repl) Cleanup() {
	r.History.Save()
	r.rl.Close()
}

// Create a new Repl.
func NewRepl() *Repl {
	r := new(Repl)
	r.History = HistoryTrie{Trie: patricia.NewTrie()}
	return r
}

// Enter the loop.
func (r *Repl) Start() error {
	var err error

	r.Suggester.Create("exit")

	r.rl, err = readline.NewEx(&readline.Config{
		AutoComplete:           &r.Suggester,
		UniqueEditLine:         false,
		Prompt:                 "\033[31m»\033[0m ",
		DisableAutoSaveHistory: true,
		HistoryLimit:           -1,
		HistorySearchFold:      false,
		VimMode:                false,
	})

	if err != nil {
		return err
	}

	if err = r.History.Load(false); err != nil {
		return err
	}

	r.rl.Config.SetListener(r.HandleKey)

	defer r.Cleanup()

	r.rl.CaptureExitSignal()

	var line string
	var commands []string
	var command string

	// We ignore interrupt signals to prevent closing after the subprocess
	// is interrupted.
	signal.Ignore(os.Interrupt)

	for {
		line, err = r.rl.Readline()

		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			return nil
		}

		line = strings.ToValidUTF8(line, "")

		if len(line) == 0 {
			continue
		}

		r.History.Push(patricia.Prefix(line))

		if line == "exit" {
			return nil
		}

		commands = strings.Split(line, "&&")

		for _, command = range commands {
			if line = strings.TrimSpace(command); r.Runner(strings.Fields(line)) != nil {
				break
			}
		}
	}

	return nil
}
