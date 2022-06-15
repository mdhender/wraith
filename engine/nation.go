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

package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Nation configuration
type Nation struct {
	Store       string `json:"store"` // path to store data
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Speciality  string `json:"speciality"`
	Government  struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
	} `json:"government"`
	Homeworld struct {
		Name     string `json:"name"`
		Location struct {
			X     int `json:"x"`
			Y     int `json:"y"`
			Z     int `json:"z"`
			Orbit int `json:"orbit"`
		} `json:"location"`
	} `json:"homeworld"`
}

// CreateNation creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func CreateNation(id int, name, descr, path string, overwrite bool) (*Nation, error) {
	s := &Nation{
		Store:       filepath.Clean(filepath.Join(path, fmt.Sprintf("%d", id))),
		Id:          id,
		Name:        name,
		Description: descr,
	}
	if _, err := os.Stat(filepath.Join(s.Store, "store.json")); err == nil {
		log.Printf("nation store exists %q %v\n", filepath.Join(s.Store, "store.json"), overwrite)
		if !overwrite {
			return nil, errors.New("nation store exists")
		}
	}
	return s, s.Write()
}

// LoadNation loads an existing store.
// It returns any errors.
func LoadNation(path string) (*Nation, error) {
	s := &Nation{
		Store: filepath.Clean(path),
	}
	return s, s.Read()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Nation) Read() error {
	b, err := ioutil.ReadFile(filepath.Join(s.Store, "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Nation) Write() error {
	if s.Store == "" {
		return errors.New("missing nation store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(s.Store, "store.json"), b, 0600)
}
