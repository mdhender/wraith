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
)

type System struct {
	GameId int
	Id     int
	X      int
	Y      int
	Z      int
}

func (s *Store) AddSystem(g Game, x, y, z int) (System, error) {
	if s.db == nil {
		return System{}, ErrNoConnection
	}

	stmt, err := s.db.Prepare("insert into systems (game_id, x, y, z) values(?, ?, ?, ?)")
	if err != nil {
		return System{}, errors.Wrap(err, "prepare insert new systems")
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)
	r, err := stmt.Exec(g.Id, x, y, z)
	if err != nil {
		return System{}, errors.Wrap(err, "exec insert new systems")
	}
	id, err := r.LastInsertId()
	if err != nil {
		return System{}, errors.Wrap(err, "fetch id new systems")
	}

	return System{Id: int(id), X: x, Y: y, Z: z}, nil
}
