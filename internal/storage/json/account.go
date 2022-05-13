/*
 * Wraith Game Engine
 * Copyright (c) 2022 Michael D. Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package json

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Accounts defines a list of account data.
type Accounts []Account

// Account defines the format for storing account data as a JSON object.
type Account struct {
	// Id is the unique identifier for an account.
	// example: "00000000-0000-0000-0000-000000000000"
	Id string `json:"id"`
	// Handle is the display name for the account.
	// example: "jimi"
	Handle string `json:"handle"`
	// Salt is the salt used to hash the account's secret.
	// example: [83, 65, 76, 84]
	Salt []byte `json:"salt"`
	// Secret is the hashed account secret.
	// example: [83, 69, 67, 82, 69, 84]
	Secret []byte `json:"secret"`
	// Roles is the list of roles assigned to the account.
	// example: ["foo","bah"]
	Roles []string `json:"roles"`
	// Games is the list of games that the account has access to.
	Games []AccountGame `json:"games"`
}

// AccountGame defines the format for storing game data as a JSON object.
type AccountGame struct {
	// Id is the unique identifier for the game.
	// example: "abcd-ef-aabb"
	Id string `json:"id"`
	// Roles is the list of roles assigned to the account for this game.
	// example: ["player"]
	Roles []string `json:"roles"`
}

// AccountsDriver is used to interact with the JSON database.
type AccountsDriver struct {
	sync.Mutex
	path string // the path to the JSON database
	data Accounts
}

// NewAccountsDriver creates a new JSON database at the given path and returns a AccountsDriver to interact with the database.
func NewAccountsDriver(path string) (*AccountsDriver, error) {
	driver := &AccountsDriver{path: filepath.Clean(path)}
	// load the database if it already exists
	err := driver.read()
	if err != nil {
		if !os.IsNotExist(err) {
			// we don't know how to handle other errors, so pass them back to the caller
			return nil, err
		}
		// the database does not exist, so we must create it
		if err = driver.write(); err != nil {
			return nil, err
		}
	}
	return driver, nil
}

// Create locks the database, creates a new record, then writes the updated database to file.
func (d *AccountsDriver) Create(a Account) error {
	d.Lock()
	defer d.Unlock()

	if a.Id == "" {
		// let caller create their own identifiers, maybe for unit testing?
		a.Id = uuid.New().String()
	} else if a.Id != strings.TrimSpace(a.Id) {
		return errors.New("invalid id")
	}
	if a.Handle != strings.TrimSpace(a.Handle) {
		return errors.New("invalid handle")
	} else if a.Handle == "" {
		return errors.New("missing handle")
	}

	set := d.data
	d.data = nil
	for _, data := range set {
		// forbid duplicate identifiers
		if a.Id == data.Id {
			// restore the database then return the error
			d.data = set
			return errors.New("duplicate id")
		} else if a.Handle == data.Handle {
			// restore the database then return the error
			d.data = set
			return errors.New("duplicate handle")
		}
		d.data = append(d.data, data)
	}
	// create new record and add it to the database
	d.data = append(d.data, a)
	if err := d.write(); err != nil {
		// restore the database then return the error
		d.data = set
		return err
	}
	return nil
}

// Delete locks the database, removes all records that satisfy the filter, then writes the updated database to file.
func (d *AccountsDriver) Delete(filter func(Account) bool) error {
	d.Lock()
	defer d.Unlock()

	set, deleted := d.data, false
	d.data = nil
	for _, data := range set {
		if filter(data) {
			deleted = true
		} else {
			d.data = append(d.data, data)
		}
	}
	// no need to write out the database if no deletes were made
	if !deleted {
		// restore the database then return the error
		d.data = set
		return nil
	}
	if err := d.write(); err != nil {
		// restore the database then return the error
		d.data = set
		return err
	}
	return nil
}

// DeleteById locks the database, removes all records that match the identifier, then writes the updated database to file.
func (d *AccountsDriver) DeleteById(id string) error {
	return d.Delete(func(a Account) bool {
		return a.Id == id
	})
}

// Read locks the database, creates a set of records from the database that satisfy the filter, then returns that set.
func (d *AccountsDriver) Read(filter func(account Account) bool) (set Accounts) {
	d.Lock()
	defer d.Unlock()

	for _, data := range d.data {
		if filter(data) {
			set = append(set, data)
		}
	}
	return set
}

// ReadAll locks the database, creates a set of all records from the database, then returns that set.
func (d *AccountsDriver) ReadAll() (set Accounts) {
	return d.Read(func(account Account) bool {
		return true
	})
}

// Update locks the database, updates all records with a matching identifier, then tries to write the updated database to file.
func (d *AccountsDriver) Update(a Account) error {
	d.Lock()
	defer d.Unlock()

	if a.Handle != strings.TrimSpace(a.Handle) {
		return errors.New("invalid handle")
	} else if a.Handle == "" {
		return errors.New("missing handle")
	}

	set, found := d.data, false
	d.data = nil
	for _, data := range set {
		if a.Id == data.Id {
			// update existing record and write to the database
			d.data = append(d.data, a)
			found = true
		} else {
			// keep the existing record
			d.data = append(d.data, data)
		}
	}
	// no need to write out the database if no updates were made
	if !found {
		// restore the database then return the error
		d.data = set
		return nil
	}
	if err := d.write(); err != nil {
		// restore the database then return the error
		d.data = set
		return err
	}
	return nil
}

// read data from the file
func (d *AccountsDriver) read() error {
	b, err := os.ReadFile(d.path)
	if err != nil {
		return err
	}
	var data struct {
		Version string   `json:"version"`
		Data    Accounts `json:"data"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	d.data = data.Data
	return nil
}

// write marshaled data to the database file
func (d *AccountsDriver) write() error {
	data := struct {
		Version string   `json:"version"`
		Data    Accounts `json:"data"`
	}{
		Version: "0.1.0",
		Data:    d.data,
	}
	b, err := json.MarshalIndent(data, "", "\t")
	if err == nil {
		err = ioutil.WriteFile(d.path, b, 0600)
	}
	return err
}
