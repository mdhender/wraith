////////////////////////////////////////////////////////////////////////////////
// wraith - the wraith game engine and server
// Copyright (c) 2022 Michael D. Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
////////////////////////////////////////////////////////////////////////////////

package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Global configuration
type Global struct {
	Self       string `json:"self"` // path to this file
	User       string `json:"user"`
	Password   string `json:"password"`
	Schema     string `json:"schema"`
	SchemaFile string `json:"schema-file"`
}

// CreateGlobal creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func CreateGlobal(filename, user, password, schema, schemaFile string, overwrite bool) (*Global, error) {
	s := &Global{
		Self:       filepath.Clean(filename),
		User:       user,
		Password:   password,
		Schema:     schema,
		SchemaFile: schemaFile,
	}
	if _, err := os.Stat(s.Self); err == nil {
		if !overwrite {
			return nil, errors.New("configuration file exists")
		}
	}
	return s, s.Write()
}

// LoadGlobal loads an existing store.
// It returns any errors.
func LoadGlobal(filename string) (*Global, error) {
	s := &Global{Self: filepath.Clean(filename)}
	return s, s.Read()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Global) Read() error {
	b, err := ioutil.ReadFile(s.Self)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Global) Write() error {
	if s.Self == "" {
		return errors.New("missing global store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Self, b, 0600)
}
