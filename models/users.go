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

package models

import (
	"database/sql"
	"github.com/pkg/errors"
	"log"
	"strings"
	"unicode"
)

// User data
type User struct {
	Id           int
	Email        string
	Handle       string
	HashedSecret string
}

// CreateUser adds a new user to the store if it passes validation
func (s *Store) CreateUser(u User) (User, error) {
	if s.db == nil {
		return User{}, ErrNoConnection
	}

	if u.Email = strings.ToLower(strings.TrimSpace(u.Email)); u.Email == "" {
		return User{}, errors.Wrap(ErrMissingField, "email")
	}
	if u.Handle = strings.ToLower(strings.TrimSpace(u.Handle)); u.Handle == "" {
		return User{}, errors.Wrap(ErrMissingField, "handle")
	}
	for _, r := range u.Handle { // check for invalid runes in the field
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return User{}, errors.Wrap(ErrInvalidField, "handle: invalid rune")
		}
	}
	if u.HashedSecret == "" {
		return User{}, errors.Wrap(ErrMissingField, "secret")
	}

	stmt, err := s.db.Prepare("insert into users (email, handle, hashed_secret) values(?, ?, ?)")
	if err != nil {
		return User{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)

	r, err := stmt.Exec(u.Email, u.Handle, u.HashedSecret)
	if err != nil {
		return User{}, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return User{}, err
	}
	u.Id = int(id)

	return u, nil
}

// SelectUserByEmail returns the user that matches the email
func (s *Store) SelectUserByEmail(email string) (User, error) {
	if s.db == nil {
		return User{}, ErrNoConnection
	}
	email = strings.ToLower(email)

	stmt, err := s.db.Prepare("select id, email, handle, hashed_secret from users where email = ?")
	if err != nil {
		return User{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)

	var u User
	err = stmt.QueryRow(strings.ToLower(email)).Scan(&u.Id, &u.Email, &u.Handle, &u.HashedSecret)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

// SelectUserByHandle returns the user that matches the handle
func (s *Store) SelectUserByHandle(handle string) (User, error) {
	if s.db == nil {
		return User{}, ErrNoConnection
	}

	stmt, err := s.db.Prepare("select id, email, handle, hashed_secret from users where handle = ?")
	if err != nil {
		return User{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)

	var u User
	err = stmt.QueryRow(strings.ToLower(handle)).Scan(&u.Id, &u.Email, &u.Handle, &u.HashedSecret)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

// SelectUserById returns the user that matches the id
func (s *Store) SelectUserById(id int) (User, error) {
	if s.db == nil {
		return User{}, ErrNoConnection
	}

	stmt, err := s.db.Prepare("select id, email, handle, hashed_secret from users where id = ?")
	if err != nil {
		return User{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)

	var u User
	err = stmt.QueryRow(id).Scan(&u.Id, &u.Email, &u.Handle, &u.HashedSecret)
	if err != nil {
		return User{}, err
	}

	return u, nil
}
