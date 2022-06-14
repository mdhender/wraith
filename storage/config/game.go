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
	"log"
	"os"
	"path/filepath"
)

// Game configuration
type Game struct {
	Store        string         `json:"store"` // path to store data
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	NationsIndex []NationsIndex `json:"nations-index"`
	TurnsIndex   []TurnsIndex   `json:"turns-index"`
}

type TurnsIndex struct{}

// CreateGame creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func CreateGame(name, descr, path string, overwrite bool) (*Game, error) {
	s := &Game{
		Store:        filepath.Clean(path),
		Name:         name,
		Description:  descr,
		NationsIndex: []NationsIndex{},
		TurnsIndex:   []TurnsIndex{},
	}
	if _, err := os.Stat(filepath.Join(s.Store, "store.json")); err == nil {
		log.Printf("game store exists %q %v\n", filepath.Join(s.Store, "store.json"), overwrite)
		if !overwrite {
			return nil, errors.New("game store exists")
		}
	}
	return s, s.Write()
}

// LoadGame loads an existing store.
// It returns any errors.
func LoadGame(path string) (*Game, error) {
	s := &Game{
		Store:        filepath.Clean(path),
		NationsIndex: []NationsIndex{},
		TurnsIndex:   []TurnsIndex{},
	}
	return s, s.Read()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Game) Read() error {
	b, err := ioutil.ReadFile(filepath.Join(s.Store, "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Game) Write() error {
	if s.Store == "" {
		return errors.New("missing game store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(s.Store, "store.json"), b, 0600)
}
