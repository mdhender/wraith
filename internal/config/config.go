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
	ConfigFile string `json:"-"`
	//Accounts map[string]Account `json:"accounts,omitempty"`
	Identity struct {
		Cost       int `json:"cost,omitempty"` // bcrypt's hashing cost
		Repository struct {
			JSONFile string `json:"json_file,omitempty"`
		} `json:"repository,omitempty"`
	} `json:"identity,omitempty"`
	Secrets struct {
		Signing string `json:"signing,omitempty"` // plain-text
		Sysop   string `json:"sysop,omitempty"`   // plain-text
	} `json:"secrets"`
	Server struct {
		Host string `json:"host,omitempty"`
		Port string `json:"port"`
	} `json:"server"`
	//Users map[string]User `json:"users,omitempty"`
}

//type Account struct {
//	Id     string `json:"id,omitempty"`
//	Handle string `json:"handle,omitempty"`
//}

//type Identity struct {
//	Id           string `json:"id,omitempty"`
//	Handle       string `json:"handle,omitempty"`
//	Secret       string `json:"secret,omitempty"`        // plain-text
//	HashedSecret string `json:"hashed_secret,omitempty"` // hashed and b64-encoded
//}

//type User struct {
//	Id    string `json:"id,omitempty"`
//	Email string `json:"email,omitempty"`
//}

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
