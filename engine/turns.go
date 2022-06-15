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
	"io/ioutil"
	"os"
	"path/filepath"
)

// Turns configuration
type Turns struct {
	Store string       `json:"store"` // path to store data
	Index []TurnsIndex `json:"index"`
}

type TurnsIndex struct {
	Id int `json:"id"` // unique turn number
}

// CreateTurns creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func CreateTurns(path string, overwrite bool) (*Turns, error) {
	s := &Turns{
		Store: filepath.Clean(filepath.Join(path, "turns")),
		Index: []TurnsIndex{},
	}
	if _, err := os.Stat(filepath.Join(s.Store, "store.json")); err == nil {
		if !overwrite {
			return nil, errors.New("turns store exists")
		}
	}
	return s, s.Write()
}

// LoadTurns loads an existing store.
// It returns any errors.
func LoadTurns(path string) (*Turns, error) {
	s := &Turns{
		Store: filepath.Clean(filepath.Join(path, "turns")),
		Index: []TurnsIndex{},
	}
	return s, s.Read()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Turns) Read() error {
	b, err := ioutil.ReadFile(filepath.Join(s.Store, "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Turns) Write() error {
	if s.Store == "" {
		return errors.New("missing turns store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(s.Store, "store.json"), b, 0600)
}
