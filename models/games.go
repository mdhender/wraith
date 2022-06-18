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

type Game struct {
	Id        int
	EffDt     time.Time
	EndDt     time.Time
	ShortName string
	Name      string
}

type GameTurn struct {
	GameId int
	TurnNo int
	EffDt  time.Time
	EndDt  time.Time
}

// CreateGame adds a new game to the store if it passes validation
func (s *Store) CreateGame(name, shortName string, startDt time.Time) (Game, error) {
	if s.db == nil {
		return Game{}, ErrNoConnection
	}

	if shortName = strings.ToUpper(strings.TrimSpace(shortName)); shortName == "" {
		return Game{}, errors.Wrap(ErrMissingField, "short name")
	}
	for _, r := range shortName { // check for invalid runes in the field
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return Game{}, errors.Wrap(ErrInvalidField, "short name: invalid rune")
		}
	}
	if name = strings.TrimSpace(name); name == "" {
		name = shortName
	}

	// check for duplicate keys
	stmtDup, err := s.db.Prepare("select ifnull(count(id), 0) from game where short_name = ?")
	if err != nil {
		return Game{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmtDup)
	var count int
	err = stmtDup.QueryRow(shortName).Scan(&count)
	if err != nil {
		return Game{}, err
	}
	if count != 0 {
		return Game{}, ErrDuplicateKey
	}

	g := Game{
		Id:        s.nextGameId(),
		EffDt:     startDt,
		EndDt:     s.endOfTime,
		ShortName: shortName,
		Name:      name,
	}

	createGame, err := s.db.Prepare("insert into game (id, effdt, enddt, short_name, name) values(?, ?, ?, ?, ?)")
	if err != nil {
		return Game{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(createGame)
	_, err = createGame.Exec(g.Id, g.EffDt, g.EndDt, g.ShortName, g.Name)
	if err != nil {
		return Game{}, err
	}

	createTurn, err := s.db.Prepare("insert into game_turn (game_id, turn_no, effdt, enddt) values(?, ?, ?, ?)")
	if err != nil {
		return Game{}, err
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(createTurn)
	_, err = createTurn.Exec(g.Id, 0, startDt, s.endOfTime)
	if err != nil {
		return Game{}, err
	}

	return g, nil
}

func (s *Store) nextGameId() (id int) {
	stmt, err := s.db.Prepare("select ifnull(max(id), 0) from game")
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
