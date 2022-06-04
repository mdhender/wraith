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

package identity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type IdentityRepository interface {
	GetById(id string) (Identity, bool)
	Add(id Identity) bool
	Remove(id string) bool
}

// store implements a file based data store for identity data.
type store struct {
	sync.Mutex

	// semantic version of the data
	Version string `json:"version"`
	// configuration information for bcrypt package
	BCrypt struct {
		Cost int `json:"cost"`
	} `json:"bcrypt"`
	// map of credentials indexed by id
	Identities map[string]*identityValue `json:"identities"`

	path    string // path to file
	emails  map[string]*identityValue
	handles map[string]*identityValue
}

// identityValue defines information for a user as stored in the file.
type identityValue struct {
	Id           string   `json:"id"`
	Handle       string   `json:"handle"`
	Email        string   `json:"email"`
	Secret       string   `json:"secret,omitempty"`
	HashedSecret string   `json:"hashed_secret"`
	Roles        []string `json:"roles"`
}

// bootstrap creates and initializes a new identity store.
func bootstrap(cost int) (*store, error) {
	if cost < bcrypt.DefaultCost {
		cost = bcrypt.DefaultCost
	}
	if cost < bcrypt.MinCost {
		cost = bcrypt.MinCost
	}

	s := &store{
		Version:    "0.1.0",
		Identities: make(map[string]*identityValue),
		path:       ".",
		emails:     make(map[string]*identityValue),
		handles:    make(map[string]*identityValue),
	}
	s.BCrypt.Cost = cost

	return s, nil
}

// loadStore returns an existing store.
// It forces the `path` field to match the actual file loaded.
func loadStore(path string) (*store, error) {
	s := &store{
		Identities: make(map[string]*identityValue),
	}
	s.pathSet(path) // set the path for the store
	if err := s.fileRead(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *store) create(i identityValue) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.Identities[i.Handle]; ok {
		return errors.New("duplicate handle")
	}
	if _, ok := s.Identities[i.Email]; ok {
		return errors.New("duplicate email")
	}
	v := identityValue{
		Id:           uuid.New().String(),
		Handle:       i.Handle,
		Email:        i.Email,
		HashedSecret: i.HashedSecret,
		Roles:        i.Roles,
	}
	s.Identities[v.Id] = &v
	return nil
}

func (s *store) filterByHandle(handle string) (set []identityValue) {
	s.Lock()
	defer s.Unlock()

	if i, ok := s.handles[handle]; ok {
		set = append(set, identityValue{
			Id:           i.Id,
			Email:        i.Email,
			Handle:       i.Handle,
			HashedSecret: i.HashedSecret,
			Roles:        i.Roles,
		})
	}
	return set
}

func (s *store) filterByEmail(email string) (set []identityValue) {
	s.Lock()
	defer s.Unlock()

	if i, ok := s.emails[email]; ok {
		set = append(set, identityValue{
			Id:           i.Id,
			Email:        i.Email,
			Handle:       i.Handle,
			HashedSecret: i.HashedSecret,
			Roles:        i.Roles,
		})
	}
	return set
}

func (s *store) filterById(id string) (set []identityValue) {
	s.Lock()
	defer s.Unlock()

	if i, ok := s.Identities[id]; ok {
		set = append(set, identityValue{
			Id:           i.Id,
			Email:        i.Email,
			Handle:       i.Handle,
			HashedSecret: i.HashedSecret,
			Roles:        i.Roles,
		})
	}
	return set
}

// fileRead reads the store from a JSON file.
// It returns any errors.
// Warning: the caller should have a lock before reading from the file.
func (s *store) fileRead() error {
	b, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}
	fs := &store{
		Identities: make(map[string]*identityValue),
	}
	if err := json.Unmarshal(b, fs); err != nil {
		return err
	}
	if s.Version != fs.Version {
		return errors.New("incompatible versions")
	}
	fs.emails = make(map[string]*identityValue)
	fs.handles = make(map[string]*identityValue)
	for _, v := range fs.Identities {
		if _, ok := fs.emails[v.Email]; ok {
			return errors.New("duplicate e-mail")
		}
		fs.emails[v.Email] = v
		if _, ok := fs.handles[v.Handle]; ok {
			return errors.New("duplicate handle")
		}
		fs.handles[v.Handle] = v
	}
	s.BCrypt = fs.BCrypt
	s.Identities = fs.Identities

	return nil
}

// fileWrite writes the store to a JSON file.
// It returns any errors.
// Warning: the caller should have a lock before writing to the file.
func (s *store) fileWrite() error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.path, b, 0600)
}

// pathGet returns the current path to the JSON file.
func (s *store) pathGet() string {
	return s.path
}

// pathSet forces the `path` to the value given.
func (s *store) pathSet(path string) {
	s.path = filepath.Clean(path)
}
