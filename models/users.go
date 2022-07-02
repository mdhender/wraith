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
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"unicode"
)

// CreateUser adds a new user to the store if it passes validation
func (s *Store) CreateUser(displayHandle, handle, email, secret string) error {
	if s.db == nil {
		return ErrNoConnection
	}
	return s.createUser(displayHandle, handle, email, secret)
}

func (s *Store) FetchUser(id int) (*User, error) {
	if s.db == nil {
		return nil, ErrNoConnection
	}
	return s.fetchUser(id)
}

// FetchUserByEmail returns the user that matches the email
func (s *Store) FetchUserByEmail(email string) (*User, error) {
	if s.db == nil {
		return nil, ErrNoConnection
	}
	return s.fetchUserByEmail(email)
}

func (s *Store) FetchUserByHandle(handle string) (*User, error) {
	if s.db == nil {
		return nil, ErrNoConnection
	}
	return s.fetchUserByHandle(handle)
}

func (s *Store) UpdateUserSecret(id int, secret string) error {
	if s.db == nil {
		return ErrNoConnection
	}
	return s.updateUserSecret(id, secret)
}

func (s *Store) FetchUserByCredentials(handle, secret string) (*User, error) {
	if s.db == nil {
		return nil, ErrNoConnection
	}
	u, err := s.fetchUserByCredentials(handle, secret)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) createUser(displayHandle, handle, email, secret string) error {
	now := time.Now()

	displayHandle = strings.TrimSpace(displayHandle)
	if displayHandle == "" {
		return fmt.Errorf("display-handle: %w", ErrMissingField)
	}
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return fmt.Errorf("email: %w", ErrMissingField)
	}
	handle = strings.ToLower(strings.TrimSpace(handle))
	if handle == "" {
		return fmt.Errorf("handle: %w", ErrMissingField)
	}
	for _, r := range handle { // check for invalid runes in the field
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' || r == ' ') {
			return fmt.Errorf("handle: invalid rune: %w", ErrInvalidField)
		}
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("createUser: hash secret: %w", err)
	}

	// get a transaction with a deferred rollback in case things fail
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("createUser: beginTx: %w", err)
	}
	defer tx.Rollback()

	// check for duplicate handle
	var matches int
	row := tx.QueryRow("select ifnull(count(id), 0) from users where handle = ?", handle)
	err = row.Scan(&matches)
	if err != nil {
		return err
	} else if matches != 0 {
		return fmt.Errorf("createUser: %w", ErrDuplicateKey)
	}

	// check for duplicate email or handle
	row = tx.QueryRow("select ifnull(count(user_id), 0) from user_profile where (effdt <= ? and ? < enddt) and email = ?", now, now, email)
	err = row.Scan(&matches)
	if err != nil {
		return err
	} else if matches != 0 {
		return fmt.Errorf("createUser: %w", ErrDuplicateKey)
	}

	r, err := tx.ExecContext(s.ctx, "insert into users (handle, hashed_secret) values (?, ?)", handle, string(hashedPasswordBytes))
	if err != nil {
		return fmt.Errorf("createUser: insert: %w", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("createUser: fetchId: %w", err)
	}
	uid := int(id)

	_, err = tx.ExecContext(s.ctx, "insert into user_profile (user_id, effdt, enddt, handle, email) values (?, ?, str_to_date('2099/12/31', '%Y/%m/%d'), ?, ?)",
		uid, now, displayHandle, email)
	if err != nil {
		return fmt.Errorf("createUser: insert: %w", err)
	}

	return tx.Commit()
}

func (s *Store) fetchUser(id int) (*User, error) {
	now := time.Now()
	u := User{Profiles: []*UserProfile{&UserProfile{}}}

	row := s.db.QueryRow("select id, users.handle, user_profile.handle, user_profile.email from users, user_profile where id = ? and user_profile.user_id = users.id and (effdt <= ? and ? < enddt)", id, now, now)
	err := row.Scan(&u.Id, &u.Handle, &u.Profiles[0].Handle, &u.Profiles[0].Email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Store) fetchUserByCredentials(handle, secret string) (*User, error) {
	now := time.Now()
	u, up := User{}, UserProfile{}

	row := s.db.QueryRow(`
		select id, hashed_secret, users.handle, user_profile.handle, user_profile.email
		from users, user_profile
		where users.handle = ?
		  and user_profile.user_id = users.id
		  and (user_profile.effdt <= ? and ? < user_profile.enddt)`, strings.ToLower(handle), now, now)
	var hashedSecret string
	err := row.Scan(&u.Id, &hashedSecret, &u.Handle, &up.Handle, &up.Email)
	if err != nil {
		return nil, err
	} else if bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(secret)) != nil {
		return nil, ErrUnauthorized
	}

	u.Profiles = []*UserProfile{&up}

	return &u, nil
}

func (s *Store) fetchUserByEmail(email string) (*User, error) {
	now := time.Now()

	var id int
	row := s.db.QueryRow("select user_id from user_profile where email = ? and (effdt <= ? and ? < enddt)", strings.ToLower(email), now, now)
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return s.fetchUser(id)
}

func (s *Store) fetchUserByHandle(handle string) (*User, error) {
	var id int

	row := s.db.QueryRow("select id from users where handle = ?", strings.ToLower(handle))
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return s.fetchUser(id)
}

func (s *Store) updateUserSecret(id int, secret string) error {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("updateUserSecret: hash secret: %w", err)
	}
	_, err = s.db.ExecContext(s.ctx, "update users set hashed_secret = ? where id = ?", string(hashedPasswordBytes), id)
	if err != nil {
		return fmt.Errorf("updateUserSecret: %w", err)
	}
	return nil
}
