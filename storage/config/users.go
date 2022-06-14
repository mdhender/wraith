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
	Store string       `json:"store"` // default path to store data
	Index []UsersIndex `json:"index"`
}

type UsersIndex struct {
	Id     string `json:"id"`
	Handle string `json:"handle"`
}

// CreateUsers creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func CreateUsers(path string, overwrite bool) (*Users, error) {
	s := &Users{
		Store: filepath.Clean(filepath.Join(path, "users")),
		Index: []UsersIndex{},
	}
	if _, err := os.Stat(filepath.Join(s.Store, "store.json")); err == nil {
		if !overwrite {
			return nil, errors.New("users store exists")
		}
	}
	return s, s.Write()
}

// LoadUsers loads an existing store.
// It returns any errors.
func LoadUsers(path string) (*Users, error) {
	s := &Users{
		Store: filepath.Clean(filepath.Join(path, "users")),
		Index: []UsersIndex{},
	}
	return s, s.Read()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Users) Read() error {
	b, err := ioutil.ReadFile(filepath.Join(s.Store, "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Users) Write() error {
	if s.Store == "" {
		return errors.New("missing users store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(s.Store, "store.json"), b, 0600)
}
