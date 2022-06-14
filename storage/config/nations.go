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

// Nations configuration
type Nations struct {
	Store string         `json:"store"` // path to store data
	Index []NationsIndex `json:"index"`
}

type NationsIndex struct {
	Id   int    `json:"id"`   // unique identifier for nation
	Name string `json:"name"` // name of nation
	Path string `json:"path"` // path to the species game data
}

// CreateNations creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func CreateNations(path string, overwrite bool) (*Nations, error) {
	s := &Nations{
		Store: filepath.Clean(filepath.Join(path, "nations")),
		Index: []NationsIndex{},
	}
	if _, err := os.Stat(filepath.Join(s.Store, "store.json")); err == nil {
		if !overwrite {
			return nil, errors.New("nations store exists")
		}
	}
	return s, s.Write()
}

// LoadNations loads an existing store.
// It returns any errors.
func LoadNations(path string) (*Nations, error) {
	s := &Nations{
		Store: filepath.Clean(filepath.Join(path, "nations")),
		Index: []NationsIndex{},
	}
	return s, s.Read()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Nations) Read() error {
	b, err := ioutil.ReadFile(filepath.Join(s.Store, "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Nations) Write() error {
	if s.Store == "" {
		return errors.New("missing nations store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(s.Store, "store.json"), b, 0600)
}
