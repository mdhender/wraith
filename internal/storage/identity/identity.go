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

// Package identity implements a store for identification data.
package identity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mdhender/jsonwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Store defines
type Store struct {
	sync.Mutex
	// semantic version of the data
	Version string `json:"version"`
	// path to file
	Filename string `json:"file_name"`
	// configuration information for bcrypt package
	BCrypt struct {
		Cost int `json:"cost"`
	} `json:"bcrypt"`
	// map of credentials indexed by id
	Credentials map[string]Credentials `json:"credentials"`
	Keys        map[string]*Key        `json:"keys"`
	signers     map[string]*jsonwt.Factory
	authSigner  *jsonwt.Factory
}

// Credentials defines information for a user.
type Credentials struct {
	Id           string   `json:"id"`
	Handle       string   `json:"handle"`
	Email        string   `json:"email"`
	Secret       string   `json:"-"`
	HashedSecret string   `json:"hashed_secret"`
	Roles        []string `json:"roles"`
}

// Key defines information for keys used to sign authorization tokens.
type Key struct {
	Id      string `json:"id"`
	Public  []byte `json:"public"`
	Private []byte `json:"private"`
}

// Bootstrap creates and initializes a new Identity store.
func Bootstrap(filename string, cost int, signer *jsonwt.Factory) (*Store, error) {
	if signer == nil {
		return nil, errors.New("missing signer")
	}

	if cost < bcrypt.DefaultCost {
		cost = bcrypt.DefaultCost
	}
	if cost < bcrypt.MinCost {
		cost = bcrypt.MinCost
	}

	s := &Store{
		Version:     "0.1.0",
		Filename:    filepath.Clean(filename),
		Credentials: make(map[string]Credentials),
		Keys:        make(map[string]*Key),
		signers:     make(map[string]*jsonwt.Factory),
	}
	s.BCrypt.Cost = cost
	s.signers[signer.ID()] = signer
	s.authSigner = signer

	if _, err := os.Stat(s.Filename); err == nil {
		return nil, errors.New("cowardly refusing to overwrite existing identity store")
	}
	if err := s.Write(); err != nil {
		return nil, err
	}

	return s, nil
}

// Load returns an existing Identity store.
func Load(filename string, signer *jsonwt.Factory) (*Store, error) {
	if signer == nil {
		return nil, errors.New("missing signer")
	}
	filename = filepath.Clean(filename)

	s := &Store{
		Filename:    filename,
		Credentials: make(map[string]Credentials),
		Keys:        make(map[string]*Key),
		signers:     make(map[string]*jsonwt.Factory),
	}
	if err := s.Read(); err != nil {
		return nil, err
	} else if s.Filename != filename {
		return nil, errors.New("mismatch on filename")
	}
	s.signers[signer.ID()] = signer
	s.authSigner = signer

	return s, nil
}

// AuthenticateCredentials implements this poorly.
// todo: someday worry about timing attacks
func (s *Store) AuthenticateCredentials(credentials struct {
	Handle string
	Secret string // hex-encoded passphrase
}) (string, error) {
	if len(credentials.Handle) != 0 || len(credentials.Secret) == 0 || len(credentials.Secret)%2 != 0 {
		return "", errors.New("invalid credentials")
	}

	s.Lock()
	users := s.fetch(func(u Credentials) bool {
		return credentials.Handle == u.Handle
	})
	s.Unlock()

	if len(users) != 1 {
		return "", errors.New("invalid credentials")
	}

	// bcrypt.CompareHashAndPassword returns an error if the credentials don't match the stored values.
	if err := bcrypt.CompareHashAndPassword([]byte(users[0].HashedSecret), []byte(credentials.Secret)); err != nil {
		return "", errors.New("invalid credentials")
	}

	var claim struct {
		Id    string   `json:"id"`
		Roles []string `json:"roles"`
	}
	claim.Id = users[0].Id
	claim.Roles = []string{"authenticated"}

	//t, err := s.authSigner.Token(91*time.Second, claim)
	//if err != nil {
	//	return "", err
	//}
	//return t.String(), nil

	return claim.Id, nil
}

func (s *Store) Create(user Credentials) error {
	hashedSecretBytes, err := bcrypt.GenerateFromPassword([]byte(user.Secret), s.BCrypt.Cost)
	if err != nil {
		return err
	}
	user.Id = uuid.New().String()
	user.HashedSecret = string(hashedSecretBytes)

	s.Lock()
	defer s.Unlock()

	if len(s.fetch(func(u Credentials) bool {
		return u.Handle == user.Handle || u.Email == user.Email
	})) != 0 {
		return errors.New("duplicate handle")
	}

	if err := s.create(user); err != nil {
		return err
	}

	return s.Write()
}

func (s *Store) Delete(user Credentials) error {
	s.Lock()
	defer s.Unlock()

	if err := s.delete(user.Id); err != nil {
		return err
	}

	return s.Write()
}

func (s *Store) Fetch(filter func(Credentials) bool) []Credentials {
	s.Lock()
	defer s.Unlock()

	return s.fetch(filter)
}

// Read loads the store from file.
// It returns any errors.
func (s *Store) Read() error {
	b, err := ioutil.ReadFile(s.Filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

// Write writes the store to a JSON file.
// It returns any errors.
func (s *Store) Write() error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Filename, b, 0600)
}

func (s *Store) create(user Credentials) error {
	s.Credentials[user.Id] = user
	return nil
}

func (s *Store) delete(id string) error {
	delete(s.Credentials, id)
	return nil
}

func (s *Store) fetch(filter func(Credentials) bool) []Credentials {
	var set []Credentials
	for _, user := range s.Credentials {
		if filter(user) {
			set = append(set, user)
		}
	}
	return set
}
