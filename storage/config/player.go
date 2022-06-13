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

// Players configuration
type Players struct {
	FileName string   `json:"file-name"`
	Players  []Player `json:"players"`
}

// Player configuration
type Player struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// LoadPlayers loads an existing configuration.
// It returns any errors.
func LoadPlayers(filename string) (*Players, error) {
	c := Players{FileName: filename}
	return &c, c.Read()
}

// Read loads a configuration from a JSON file.
// It returns any errors.
func (c *Players) Read() error {
	b, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}

// Write writes a configuration to a JSON file.
// It returns any errors.
func (c *Players) Write() error {
	if c.FileName == "" {
		return errors.New("missing config file name")
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.FileName, b, 0600)
}
