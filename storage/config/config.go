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
	"io/ioutil"
)

type Config struct {
	ConfigFile string `json:"config-file"`
	GamesPath  string `json:"games-path"`
}

// Read loads a configuration from file.
// It returns any errors.
func Read(c *Config) error {
	b, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}

// Write writes a configuration to a JSON file.
// It returns any errors.
func Write(name string, c *Config) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, b, 0600)
}
