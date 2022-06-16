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
	"io/ioutil"
	"path/filepath"
)

// Nations configuration
type Nations struct {
	Store string         `json:"store"` // path to store data
	Index []NationsIndex `json:"index"`
}

type NationsIndex struct {
	Id    int    `json:"id"`    // unique identifier for nation
	Name  string `json:"name"`  // name of nation
	Store string `json:"store"` // path to the species game data
}

// ReadNations loads a store from a JSON file.
// It returns any errors.
func (e *Engine) ReadNations() error {
	b, err := ioutil.ReadFile(filepath.Join(e.stores.nations.Store, "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, e.stores.nations)
}

// WriteNations writes a store to a JSON file.
// It returns any errors.
func (e *Engine) WriteNations() error {
	b, err := json.MarshalIndent(e.stores.nations, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(e.stores.nations.Store, "store.json"), b, 0600)
}
