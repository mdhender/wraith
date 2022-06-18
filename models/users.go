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
	"time"
	"unicode"
)

type User struct {
	Id     int
	EffDt  time.Time
	EndDt  time.Time
	Email  string
	Handle string
}

type UserSecret struct {
	Id           int
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

	// check for duplicate email
	stmtDup, err := s.db.Prepare("select ifnull(count(id), 0) from user where email = ?")
	if err != nil {
		return User{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmtDup)
	var count int
	err = stmtDup.QueryRow(u.Email).Scan(&count)
	if err != nil {
		return User{}, err
	}
	if count != 0 {
		return User{}, errors.Wrap(ErrDuplicateKey, "email")
	}

	// get sequence and add effective dates
	u.Id, u.EffDt, u.EndDt = s.nextUserId(), time.Now().UTC(), s.endOfTime

	stmt, err := s.db.Prepare("insert into user (id, effdt, enddt, email, handle) values(?, ?, ?, ?, ?)")
	if err != nil {
		return User{}, errors.Wrap(err, "prepare insert new user")
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)
	_, err = stmt.Exec(u.Id, s.timeToDate(u.EffDt), s.timeToDate(u.EndDt), u.Email, u.Handle)
	if err != nil {
		return User{}, errors.Wrap(err, "exec insert new user")
	}

	stmtSecret, err := s.db.Prepare("insert into user_secret (id, hashed_secret) values(?, ?)")
	if err != nil {
		return User{}, errors.Wrap(err, "prepare insert new secret")
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmtSecret)
	_, err = stmtSecret.Exec(u.Id, "*login-not-permitted*")
	if err != nil {
		return User{}, errors.Wrap(err, "exec insert new secret")
	}

	return u, nil
}

func (s *Store) nextUserId() (id int) {
	stmt, err := s.db.Prepare("select ifnull(max(id), 0) from user")
	if err != nil {
		return 0
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)
	_ = stmt.QueryRow().Scan(&id)
	return id + 1
}

// SelectUserByEmail returns the user that matches the email
func (s *Store) SelectUserByEmail(email string) (User, error) {
	if s.db == nil {
		return User{}, ErrNoConnection
	}
	email = strings.ToLower(email)

	stmt, err := s.db.Prepare("select id, effdt, enddt, email, handle from user where email = ?")
	if err != nil {
		return User{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)

	var u User
	var effDt, endDt string
	err = stmt.QueryRow(strings.ToLower(email)).Scan(&u.Id, &effDt, &endDt, &u.Email, &u.Handle)
	if err != nil {
		return User{}, err
	}
	u.EffDt, u.EndDt = s.dateToTime(effDt), s.dateToTime(endDt)

	return u, nil
}

// SelectUserByHandle returns the user that matches the handle
func (s *Store) SelectUserByHandle(handle string) (User, error) {
	if s.db == nil {
		return User{}, ErrNoConnection
	}

	stmt, err := s.db.Prepare("select id, effdt, enddt, email, handle from user where handle = ?")
	if err != nil {
		return User{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)

	var u User
	var effDt, endDt string
	err = stmt.QueryRow(strings.ToLower(handle)).Scan(&u.Id, &effDt, &endDt, &u.Email, &u.Handle)
	if err != nil {
		return User{}, err
	}
	u.EffDt, u.EndDt = s.dateToTime(effDt), s.dateToTime(endDt)

	return u, nil
}

// SelectUserById returns the user that matches the id
func (s *Store) SelectUserById(id int) (User, error) {
	if s.db == nil {
		return User{}, ErrNoConnection
	}

	stmt, err := s.db.Prepare("select id, effdt, enddt, email, handle from user where id = ?")
	if err != nil {
		return User{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)

	var u User
	var effDt, endDt string
	err = stmt.QueryRow(id).Scan(&u.Id, &u.EffDt, &u.EndDt, &u.Email, &u.Handle)
	if err != nil {
		return User{}, err
	}
	u.EffDt, u.EndDt = s.dateToTime(effDt), s.dateToTime(endDt)

	return u, nil
}
