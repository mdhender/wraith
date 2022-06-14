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

// Game configuration
type Game struct {
	Path         string `json:"path"`      // path to this store file
	GamePath     string `json:"game-path"` // path to store game data
	Name         string `json:"name"`
	Description  string `json:"description"`
	NationsStore string `json:"nations-store"`
}

// LoadGame loads an existing configuration.
// It returns any errors.
func LoadGame(path string) (*Game, error) {
	s := Game{Path: filepath.Clean(path)}
	return &s, s.Read()
}

// Create creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func (s *Game) Create(path string, overwrite bool) error {
	s.Path = filepath.Clean(filepath.Join(path, "game.json"))
	if _, err := os.Stat(s.Path); err == nil {
		if !overwrite {
			return errors.New("game store exists")
		}
	}
	return s.Write()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Game) Read() error {
	b, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Game) Write() error {
	if s.Path == "" {
		return errors.New("missing game store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Path, b, 0600)
}
