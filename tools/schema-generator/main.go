package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/sumartian/snake/cmake"
	"github.com/sumartian/snake/configuration"
	"github.com/sumartian/snake/utilities"
)

// Generate a JSON schema file from the configuration.
func generateSnakeSchema(src string, dest string) error {
	r := new(jsonschema.Reflector)
	r.Anonymous = true
	r.RequiredFromJSONSchemaTags = true
	r.DoNotReference = true
	r.AllowAdditionalProperties = false
	r.PreferYAMLSchema = true

	r.Mapper = func(t reflect.Type) *jsonschema.Schema {
		if t.Kind() == reflect.String {
			obj := jsonschema.ReflectFromType(t)
			obj.Type = ""
			obj.Version = ""
			return obj
		}

		return nil
	}

	if err := r.AddGoComments("github.com/sumartian/snake", src); err != nil {
		return err
	}

	s := r.Reflect(&configuration.Configuration{})
	s.Version = ""

	data, err := json.Marshal(s)

	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(dest, data, 0644); err != nil {
		return err
	}

	return nil
}

// Dump Snake metadata into zip files for later embedding.
func dump() error {
	root, err := os.Getwd()

	if err != nil {
		return err
	}

	tmpDataDir := filepath.Join(root, "distribution", "snake")

	err = os.MkdirAll(tmpDataDir, os.ModePerm)

	if err != nil {
		return err
	}

	err = generateSnakeSchema("./configuration", filepath.Join(tmpDataDir, "snake.schema.json"))

	if err != nil {
		return err
	}

	filesToZip := []string{
		"./snake.schema.json",
	}

	dataDir := filepath.Join(root, "data")

	err = filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fmt.Println("clean:", path)

		rel := strings.ReplaceAll(path, dataDir, "")

		filesToZip = append(filesToZip, "."+rel)

		if err != nil {
			return err
		}

		if err = os.MkdirAll(filepath.Join(tmpDataDir, filepath.Dir(rel)), os.ModePerm); err != nil {
			return err
		}

		if ext := filepath.Ext(path); ext == ".cmake" {
			data, err := cmake.Minify(path)

			if err != nil {
				return err
			}

			if err = os.WriteFile(filepath.Join(tmpDataDir, rel), data, 0644); err != nil {
				return err
			}
		} else {
			if err = utilities.CopyFile(path, filepath.Join(tmpDataDir, rel)); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err = os.Chdir(tmpDataDir); err != nil {
		return err
	}

	if err = utilities.Compress(filepath.Join(root, "distribution", "data.zip"), filesToZip); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := dump(); err != nil {
		fmt.Println(err)
	}
}
