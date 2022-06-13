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
)

// Games configuration
type Games struct {
	FileName string `json:"file-name"`
	Games    []Game `json:"games"`
}

// Game configuration
type Game struct {
	FileName    string `json:"file-name"`
	GamePath    string `json:"game-path"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// LoadGames loads an existing configuration.
// It returns any errors.
func LoadGames(filename string) (*Games, error) {
	c := Games{FileName: filename}
	return &c, c.Read()
}

// Read loads a configuration from a JSON file.
// It returns any errors.
func (c *Games) Read() error {
	b, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}

// Write writes a configuration to a JSON file.
// It returns any errors.
func (c *Games) Write() error {
	if c.FileName == "" {
		return errors.New("missing games store file name")
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.FileName, b, 0600)
}

// LoadGame loads an existing configuration.
// It returns any errors.
func LoadGame(filename string) (*Game, error) {
	c := Game{FileName: filename}
	return &c, c.Read()
}

// Read loads a configuration from a JSON file.
// It returns any errors.
func (c *Game) Read() error {
	b, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}

// Write writes a configuration to a JSON file.
// It returns any errors.
func (c *Game) Write() error {
	if c.FileName == "" {
		return errors.New("missing game store file name")
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.FileName, b, 0600)
}
