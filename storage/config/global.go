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
	Path       string `json:"path"`        // path to this store
	GamesPath  string `json:"games-path"`  // default path to store game data in
	GamesStore string `json:"games-store"` // path to the games store
	UsersPath  string `json:"users-path"`  // default path to store user data in
	UsersStore string `json:"users-store"` // path to the users store
}

// LoadGlobal loads an existing store.
// It returns any errors.
func LoadGlobal(path string) (*Global, error) {
	c := Global{Path: filepath.Clean(path)}
	return &c, c.Read()
}

// Create creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func (s *Global) Create(path string, overwrite bool) error {
	s.Path = filepath.Clean(filepath.Join(path, ".wraith.json"))
	if _, err := os.Stat(s.Path); err == nil {
		if !overwrite {
			return errors.New("global store exists")
		}
	}
	return s.Write()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Global) Read() error {
	b, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Global) Write() error {
	if s.Path == "" {
		return errors.New("missing global store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Path, b, 0600)
}
