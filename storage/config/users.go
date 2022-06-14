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

// Users configuration
type Users struct {
	Path  string `json:"path"` // path to this store file
	Users []User `json:"users"`
}

// User configuration
type User struct {
	Id     string `json:"id"`
	Handle string `json:"handle"`
}

// LoadUsers loads an existing store.
// It returns any errors.
func LoadUsers(path string) (*Users, error) {
	s := Users{Path: filepath.Clean(path)}
	return &s, s.Read()
}

// Create creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func (s *Users) Create(path string, overwrite bool) error {
	s.Path = filepath.Clean(filepath.Join(path, "users.json"))
	if _, err := os.Stat(s.Path); err == nil {
		if !overwrite {
			return errors.New("users store exists")
		}
	}
	return s.Write()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Users) Read() error {
	b, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Users) Write() error {
	if s.Path == "" {
		return errors.New("missing users store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Path, b, 0600)
}
