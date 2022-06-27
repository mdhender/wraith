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
	"fmt"
	"log"
	"math/rand"
)

func (s *Store) AddSystem(g Game, x, y, z int) (System, error) {
	if s.db == nil {
		return System{}, ErrNoConnection
	}

	stmt, err := s.db.Prepare("insert into systems (game_id, x, y, z) values(?, ?, ?, ?)")
	if err != nil {
		return System{}, fmt.Errorf("prepare insert new systems: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Printf("%+v\n", err)
		}
	}(stmt)
	r, err := stmt.Exec(g.Id, x, y, z)
	if err != nil {
		return System{}, fmt.Errorf("exec insert new systems: %w", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		return System{}, fmt.Errorf("fetch id new systems: %w", err)
	}

	return System{Id: int(id), Coords: Coordinates{X: x, Y: y, Z: z}}, nil
}

func (s *Store) genHomeSystem(id int) *System {
	system := &System{Id: id, HomeSystem: true}

	system.Stars = make([]*Star, 1, 1)
	for i := range system.Stars {
		var star *Star
		if i == 0 {
			star = s.genHomeStar(system)
		} else {
			star = s.genStar(system)
		}
		star.Sequence = string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i])
		system.Stars[i] = star
	}
	return system
}

func (s *Store) genSystem(id int) *System {
	system := &System{Id: id}

	switch rand.Intn(21) {
	case 0, 1, 2, 3, 4, 5:
		system.Stars = make([]*Star, 1, 1)
	case 6, 7, 8, 9, 10:
		system.Stars = make([]*Star, 2, 2)
	case 11, 12, 13, 14:
		system.Stars = make([]*Star, 3, 3)
	case 15, 16, 17:
		system.Stars = make([]*Star, 4, 4)
	case 18, 19:
		system.Stars = make([]*Star, 5, 5)
	case 20:
		system.Stars = make([]*Star, 6, 6)
	}
	for i := range system.Stars {
		star := s.genStar(system)
		star.Sequence = string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i])
		system.Stars[i] = star
	}

	return system
}
