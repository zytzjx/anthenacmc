// Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	Log "github.com/zytzjx/anthenacmc/loggersys"
	"gopkg.in/yaml.v2"
)

const confVer = 0

var Conf = struct {
	Version int
	CmdArgs []string `json:"cmd_args" yaml:"cmd_args"`
}{
	CmdArgs: []string{},
}

// LoadConfig loads the bee tool configuration.
// It looks for Beefile or bee.json in the current path,
// and falls back to default configuration in case not found.
func LoadConfig() {
	currentPath, err := os.Getwd()
	if err != nil {
		Log.Log.Error(err.Error())
	}

	dir, err := os.Open(currentPath)
	if err != nil {
		Log.Log.Error(err.Error())
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		Log.Log.Error(err.Error())
	}

	for _, file := range files {
		switch file.Name() {
		case "anthena.json":
			{
				err = parseJSON(filepath.Join(currentPath, file.Name()), &Conf)
				if err != nil {
					Log.Log.Errorf("Failed to parse JSON file: %s", err)
				}
				break
			}
		case "anthenafile":
			{
				err = parseYAML(filepath.Join(currentPath, file.Name()), &Conf)
				if err != nil {
					Log.Log.Errorf("Failed to parse YAML file: %s", err)
				}
				break
			}
		}
	}

	// Check format version
	if Conf.Version != confVer {
		Log.Log.Warn("Your configuration file is outdated. Please do consider updating it.")
	}

}

func parseJSON(path string, v interface{}) error {
	var (
		data []byte
		err  error
	)
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func parseYAML(path string, v interface{}) error {
	var (
		data []byte
		err  error
	)
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, v)
	return err
}
