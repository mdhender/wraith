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

// CreateGame adds a new game to the store if it passes validation
func (s *Store) CreateGame(g Game) (Game, error) {
	if s.db == nil {
		return Game{}, ErrNoConnection
	}

	if g.Id = strings.ToUpper(strings.TrimSpace(g.Id)); g.Id == "" {
		return Game{}, errors.Wrap(ErrMissingField, "id")
	}
	for _, r := range g.Id { // check for invalid runes in the field
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return Game{}, errors.Wrap(ErrInvalidField, "id: invalid rune")
		}
	}
	if g.Name = strings.TrimSpace(g.Name); g.Name == "" {
		g.Name = g.Id
	}
	if g.TurnNumber < 0 {
		return Game{}, errors.Wrap(ErrInvalidField, "turn number")
	}

	stmt, err := s.db.Prepare("insert into games (id, name, turn_no) values(?, ?, ?)")
	if err != nil {
		return Game{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)

	_, err = stmt.Exec(g.Id, g.Name, g.TurnNumber)
	if err != nil {
		return Game{}, err
	}

	return g, nil
}
