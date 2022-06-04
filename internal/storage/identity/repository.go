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
	"encoding/hex"
	isvc "github.com/mdhender/wraith/internal/services/identity"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// AuthenticateUser implements the identity service UserRepository interface
func (s *Store) AuthenticateUser(r isvc.AuthenticateUserRequest) (*isvc.AuthenticatedUserResponse, error) {
	// must supply either e-mail or handle, not both.
	if r.Email == nil && r.Handle == nil {
		return nil, errors.New("invalid credentials")
	} else if r.Email != nil && r.Handle != nil {
		return nil, errors.New("invalid credentials")
	}

	secret, err := hex.DecodeString(r.Secret)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	var users []Identity
	if r.Handle != nil {
		users = s.Fetch(func(i Identity) bool {
			return *r.Handle == i.Handle
		})
	} else {
		users = s.Fetch(func(i Identity) bool {
			return *r.Email == i.Email
		})
	}

	if len(users) != 0 {
		return nil, errors.New("invalid credentials")
	}
	user := users[0]

	// bcrypt.CompareHashAndPassword returns an error if the credentials don't match the stored values.
	err = bcrypt.CompareHashAndPassword([]byte(users[0].HashedSecret), secret)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	var claim struct {
		Id    string   `json:"id"`
		Roles []string `json:"roles"`
	}
	claim.Id = user.Id
	claim.Roles = []string{"authenticated"}

	t, err := s.authSigner.Token(91*time.Second, claim)
	if err != nil {
		return nil, err
	}

	return &isvc.AuthenticatedUserResponse{
		Id:     user.Id,
		Handle: user.Handle,
		Token:  t.String(),
	}, nil
}

// CreateUser implements the identity service UserRepository interface
func (s *Store) CreateUser(r isvc.AuthenticateUserRequest) error {
	return errors.New("not implemented")
}

// DeleteUser implements the identity service UserRepository interface
func (s *Store) DeleteUser(r isvc.AuthenticateUserRequest) error {
	return errors.New("not implemented")
}

// FetchUser implements the identity service UserRepository interface
func (s *Store) FetchUser(r isvc.AuthenticateUserRequest) error {
	return errors.New("not implemented")
}

// UpdateUser implements the identity service UserRepository interface
func (s *Store) UpdateUser(r isvc.AuthenticateUserRequest) error {
	return errors.New("not implemented")
}
