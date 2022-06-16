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
	"path/filepath"
	"strings"
	"unicode"
)

// Users configuration
type Users struct {
	Store string       `json:"store"` // default path to store data
	Index []UsersIndex `json:"index"`
}

type UsersIndex struct {
	Id     int    `json:"id"`
	Handle string `json:"handle"`
}

// AddUser adds a new user to the store
func (e *Engine) AddUser(handle string) error {
	handle = strings.TrimSpace(handle)
	if handle == "" {
		return errors.New("missing handle")
	}

	// check for invalid runes in the handle
	for _, r := range handle {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return errors.New("invalid rune in user handle")
		}
	}

	// check for duplicate handle
	for _, u := range e.stores.users.Index {
		if strings.ToLower(u.Handle) == strings.ToLower(handle) {
			return errors.New("duplicate handle")
		}
	}

	// generate an id for the user
	id := len(e.stores.users.Index)
	for _, u := range e.stores.users.Index {
		if u.Id > id {
			id = u.Id
		}
	}
	id = id + 1

	// add the new user to the users store
	e.stores.users.Index = append(e.stores.users.Index, UsersIndex{
		Id:     id,
		Handle: handle,
	})

	return e.WriteUsers()
}

// ReadUsers loads a store from a JSON file.
// It returns any errors.
func (e *Engine) ReadUsers() error {
	b, err := ioutil.ReadFile(filepath.Join(e.stores.users.Store, "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, e.stores.users)
}

// WriteUsers writes a store to a JSON file.
// It returns any errors.
func (e *Engine) WriteUsers() error {
	b, err := json.MarshalIndent(e.stores.users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(e.stores.users.Store, "store.json"), b, 0600)
}
