// Copyright (c) 2022 Sumartian Studios
//
// Snake is free software: you can redistribute it and/or modify it under the
// terms of the MIT license.

package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fxamacker/cbor/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var mutatorFlag string

func mutateGenericStringMap(inputPath string, outputPath string, from string, to string) error {
	inputData, err := ioutil.ReadFile(inputPath)

	if err != nil {
		return err
	}

	var outputData []byte
	var m map[string]interface{}

	switch from {
	case "yaml":
		if err := yaml.Unmarshal(inputData, &m); err != nil {
			return err
		}
	case "cbor":
		if err := cbor.Unmarshal(inputData, &m); err != nil {
			return err
		}
	case "json":
		if err := json.Unmarshal(inputData, &m); err != nil {
			return err
		}
	}

	switch to {
	case "json":
		if outputData, err = json.Marshal(m); err != nil {
			return err
		}
	case "cbor":
		if outputData, err = cbor.Marshal(m); err != nil {
			return err
		}
	case "yaml":
		if outputData, err = yaml.Marshal(m); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(outputPath, outputData, 0644); err != nil {
		return err
	}

	return nil
}

var mutateCmd = &cobra.Command{
	Use:   "mutate input-path output-path",
	Short: "Transform or optimize resource files",
	RunE: func(c *cobra.Command, args []string) error {
		if len(mutatorFlag) == 0 {
			return errors.New("you need to specify a mutator")
		}

		if len(args) < 2 {
			return errors.New("you need to pass input path as 'arg1' and output path as 'arg2'")
		}

		inputPath, outputPath := args[0], args[1]

		switch mutatorFlag {
		case "yaml-to-json":
			if err := mutateGenericStringMap(inputPath, outputPath, "yaml", "json"); err != nil {
				return err
			}
		case "json-to-yaml":
			if err := mutateGenericStringMap(inputPath, outputPath, "json", "yaml"); err != nil {
				return err
			}
		case "json-to-cbor":
			if err := mutateGenericStringMap(inputPath, outputPath, "json", "cbor"); err != nil {
				return err
			}
		case "yaml-to-cbor":
			if err := mutateGenericStringMap(inputPath, outputPath, "yaml", "cbor"); err != nil {
				return err
			}
		case "cbor-to-json":
			if err := mutateGenericStringMap(inputPath, outputPath, "cbor", "json"); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported mutator: %s", mutatorFlag)
		}

		fmt.Println(mutatorFlag, inputPath)

		return nil
	},
}

func init() {
	mutateCmd.PersistentFlags().StringVarP(&mutatorFlag,
		"mutator", "m", "",
		"The type of mutator to use for I/O")
}
