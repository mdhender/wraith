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

// Games configuration
type Games struct {
	Store string       `json:"store"` // default path to store data
	Index []GamesIndex `json:"index"`
}

type GamesIndex struct {
	Name  string `json:"name"`  // name of game
	Store string `json:"store"` // path to the game store file
}

// CreateGames creates a new store.
// Assumes that the path to store the data already exists.
// It returns any errors.
func CreateGames(path string, overwrite bool) (*Games, error) {
	s := &Games{
		Store: filepath.Clean(filepath.Join(path, "games")),
		Index: []GamesIndex{},
	}
	if _, err := os.Stat(filepath.Join(s.Store, "store.json")); err == nil {
		if !overwrite {
			return nil, errors.New("games store exists")
		}
	}
	return s, s.Write()
}

// LoadGames loads an existing store.
// It returns any errors.
func LoadGames(path string) (*Games, error) {
	s := &Games{
		Store: filepath.Clean(filepath.Join(path, "games")),
		Index: []GamesIndex{},
	}
	return s, s.Read()
}

// Read loads a store from a JSON file.
// It returns any errors.
func (s *Games) Read() error {
	b, err := ioutil.ReadFile(filepath.Join(s.Store, "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes a store to a JSON file.
// It returns any errors.
func (s *Games) Write() error {
	if s.Store == "" {
		return errors.New("missing games store path")
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(s.Store, "store.json"), b, 0600)
}
