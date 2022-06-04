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
	"context"
	"encoding/hex"
	"github.com/mdhender/jsonwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Service struct {
	st *store
	sf *jsonwt.Factory
}

//// bootstrap the identity signing factory
//hs256, err := signers.NewHS256([]byte("long though the bough bows before the foremast"))
//if err != nil {
//log.Fatal(err)
//}
//s.sf = jsonwt.NewFactory("る", hs256)

func NewService(sysopSecret []byte, signer jsonwt.Signer) (Service, error) {
	s := Service{sf: jsonwt.NewFactory("る", signer)}

	// create an in-memory identity store
	s.st, _ = bootstrap(bcrypt.DefaultCost)

	// create the sysop account used by the command line interface
	if sysopSecret == nil {
		sysopSecret = []byte("shout bell bottoms canyon")
	}
	hashedSecretBytes := string(s.HashSecret(sysopSecret))
	if hashedSecretBytes == "" {
		return Service{}, errors.New("internal bcrypt error")
	}
	if err := s.st.create(identityValue{
		Email:        "sysop",
		Handle:       "sysop",
		HashedSecret: hashedSecretBytes,
		Roles:        []string{"sysop"},
	}); err != nil {
		return Service{}, err
	}

	return s, nil
}

// HashSecret returns the hashed value of the secret as a slice of byte.
// It uses bcrypt internally, so it will return a nil slice if bcrypt throws an error.
// The slice is allocated and the caller assumes ownership of it.
func (s Service) HashSecret(secret []byte) []byte {
	hashedSecretBytes, err := bcrypt.GenerateFromPassword(secret, s.st.BCrypt.Cost)
	if err != nil {
		return nil
	}
	return append(make([]byte, 0, len(hashedSecretBytes)), hashedSecretBytes...)
}

// AuthenticateUser is a query that returns a valid response only if the credentials are valid.
func (s Service) AuthenticateUser(ctx context.Context, request AuthenticateUserRequest) (*UserResponse, error) {
	// verify the request before processing it.
	if len(request.Secret) == 0 {
		log.Printf("[identity] AuthenticateUser: secret is empty string\n")
		return nil, errors.New("bad request")
	} else if len(request.Secret)%2 != 0 {
		log.Printf("[identity] AuthenticateUser: secret is not valid base16 string\n")
		return nil, errors.New("bad request")
	} else if request.Email == nil && request.Handle == nil {
		log.Printf("[identity] AuthenticateUser: missing email and handle\n")
		return nil, errors.New("bad request")
	} else if request.Email != nil && request.Handle != nil {
		log.Printf("[identity] AuthenticateUser: both email and handle\n")
		return nil, errors.New("bad request")
	} else if request.Email != nil && *request.Email == "" {
		log.Printf("[identity] AuthenticateUser: email is empty string\n")
		return nil, errors.New("bad request")
	} else if request.Handle != nil && *request.Handle == "" {
		log.Printf("[identity] AuthenticateUser: handle is empty string\n")
		return nil, errors.New("bad request")
	}

	// decode base16 string into slice of bytes
	secretByes, err := hex.DecodeString(request.Secret)
	if err != nil {
		log.Printf("[identity] AuthenticateUser: secret %q: %+v\n", request.Secret, err)
		return nil, errors.New("bad request")
	}
	log.Printf("[identity] AuthenticateUser: secret %q: %q\n", request.Secret, secretByes)

	var ids []identityValue
	if request.Email != nil {
		ids = s.st.filterByEmail(*request.Email)
		if len(ids) != 1 {
			log.Printf("[identity] AuthenticateUser: email %q: not found\n", *request.Email)
		}
	} else if request.Handle != nil {
		ids = s.st.filterByHandle(*request.Handle)
		if len(ids) != 1 {
			log.Printf("[identity] AuthenticateUser: handle %q: not found\n", *request.Handle)
		}
	}
	if len(ids) != 1 {
		log.Printf("[identity] AuthenticateUser: found %d matches\n", len(ids))
		return nil, errors.New("invalid credentials")
	}
	id := ids[0]

	// bcrypt.CompareHashAndPassword returns an error if the credentials don't match the stored values.
	err = bcrypt.CompareHashAndPassword([]byte(id.HashedSecret), secretByes)
	if err != nil {
		log.Printf("[identity] AuthenticateUser: id %q: failed authentication\n", id.Id)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("[identity] AuthenticateUser: id %q: authenticated successfully\n", id.Id)

	return &UserResponse{
		Id:     id.Id,
		Email:  id.Email,
		Handle: id.Handle,
		Roles:  id.Roles,
	}, nil
}

func (s Service) CreateUser(ctx context.Context, request CreateUserRequest) (*UserResponse, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (s Service) DeleteUser(ctx context.Context, request DeleteUserRequest) (*UserResponse, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (s Service) FetchUser(ctx context.Context, request FetchUserRequest) (*UserResponse, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (s Service) UpdateUser(ctx context.Context, request UpdateUserRequest) (*UserResponse, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}
