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

// Package identity implements a store for identity data.
package identity

import (
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mdhender/jsonwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Store implements a file based data store for identity data.
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
	Users      map[string]Identity `json:"users"`
	Keys       map[string]*Key     `json:"keys"`
	signers    map[string]*jsonwt.Factory
	authSigner *jsonwt.Factory
}

// Identity defines information for a user.
type Identity struct {
	Id           string   `json:"id"`
	Handle       string   `json:"handle"`
	Email        string   `json:"email"`
	Secret       string   `json:"-"`
	HashedSecret string   `json:"hashed_secret"`
	Roles        []string `json:"roles"`
}

type UserValue struct {
	Id     string
	Handle string
	Roles  []string
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
		Version:  "0.1.0",
		Filename: filepath.Clean(filename),
		Users:    make(map[string]Identity),
		Keys:     make(map[string]*Key),
		signers:  make(map[string]*jsonwt.Factory),
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
		Filename: filename,
		Users:    make(map[string]Identity),
		Keys:     make(map[string]*Key),
		signers:  make(map[string]*jsonwt.Factory),
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
	if len(credentials.Handle) != 0 || len(credentials.Secret) == 0 {
		return "", errors.New("invalid credentials")
	}

	secret, err := hex.DecodeString(credentials.Secret)
	if err != nil {
		log.Printf("[identity] authCred: secret %q: %+v\n", credentials.Secret, err)
		return "", errors.New("invalid credentials")
	}
	log.Printf("[identity] authCred: secret %q: %q\n", credentials.Secret, secret)

	s.Lock()
	users := s.fetch(func(u Identity) bool {
		return credentials.Handle == u.Handle
	})
	s.Unlock()

	if len(users) != 1 {
		return "", errors.New("invalid credentials")
	}
	user := users[0]

	// bcrypt.CompareHashAndPassword returns an error if the credentials don't match the stored values.
	err = bcrypt.CompareHashAndPassword([]byte(users[0].HashedSecret), secret)
	if err != nil {
		log.Printf("[identity] authCred: secret %q: %+v\n", credentials.Secret, err)
		return "", errors.New("invalid credentials")
	}

	var claim struct {
		Id    string   `json:"id"`
		Roles []string `json:"roles"`
	}
	claim.Id = user.Id
	claim.Roles = []string{"authenticated"}

	t, err := s.authSigner.Token(91*time.Second, claim)
	if err != nil {
		log.Printf("[identity] authCred: secret %q: %+v\n", credentials.Secret, err)
		return "", err
	}

	return t.String(), nil
}

func (s *Store) Create(user Identity) error {
	hashedSecretBytes, err := bcrypt.GenerateFromPassword([]byte(user.Secret), s.BCrypt.Cost)
	if err != nil {
		return err
	}
	user.Id = uuid.New().String()
	user.HashedSecret = string(hashedSecretBytes)

	s.Lock()
	defer s.Unlock()

	if len(s.fetch(func(u Identity) bool {
		return u.Handle == user.Handle || u.Email == user.Email
	})) != 0 {
		return errors.New("duplicate handle")
	}

	if err := s.create(user); err != nil {
		return err
	}

	return s.Write()
}

func (s *Store) Delete(user Identity) error {
	s.Lock()
	defer s.Unlock()

	if err := s.delete(user.Id); err != nil {
		return err
	}

	return s.Write()
}

func (s *Store) Fetch(filter func(Identity) bool) []Identity {
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

func (s *Store) create(user Identity) error {
	s.Users[user.Id] = user
	return nil
}

func (s *Store) delete(id string) error {
	delete(s.Users, id)
	return nil
}

func (s *Store) fetch(filter func(Identity) bool) []Identity {
	var set []Identity
	for _, user := range s.Users {
		if filter(user) {
			set = append(set, user)
		}
	}
	return set
}
