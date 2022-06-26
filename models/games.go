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
	"github.com/pkg/errors"
	"strings"
	"unicode"
)

// CreateGame adds a new game to the store if it passes validation
func (s *Store) CreateGame(g *Game) error {
	if s.db == nil {
		return ErrNoConnection
	}

	g.ShortName = strings.ToUpper(strings.TrimSpace(g.ShortName))
	if g.ShortName == "" {
		return errors.Wrap(ErrMissingField, "short name")
	}
	for _, r := range g.ShortName { // check for invalid runes in the field
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return errors.Wrap(ErrInvalidField, "short name: invalid rune")
		}
	}

	// check for duplicate keys
	var count int
	row := s.db.QueryRow("select ifnull(count(id), 0) from games where short_name = ?", g.ShortName)
	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("createGame: count: %w", err)
	} else if count != 0 {
		return ErrDuplicateKey
	}

	g.Name = strings.TrimSpace(g.Name)
	if g.Name == "" {
		g.Name = g.ShortName
	}
	g.Description = strings.TrimSpace(g.Description)
	if g.Description == "" {
		g.Description = g.Name
	}

	return nil
}

// DeleteGame removes a game from the store.
func (s *Store) DeleteGame(g *Game) error {
	if g.Id != 0 {
		return s.deleteGame(g.Id)
	} else if g.ShortName != "" {
		return s.deleteGameByName(g.ShortName)
	}
	return ErrMissingField
}

func (s *Store) deleteGame(id int) error {
	if s.db == nil {
		return ErrNoConnection
	}
	_, err := s.db.Exec("delete from games where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) deleteGameByName(shortName string) error {
	if s.db == nil {
		return ErrNoConnection
	}
	_, err := s.db.Exec("delete from games where short_name = ?", shortName)
	if err != nil {
		return err
	}
	return nil
}
