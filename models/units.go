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
	"strings"
)

func (s *Store) CreateUnit(code, name, descr string, usesTech bool) error {
	// get a transaction with a deferred rollback in case things fail
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("createUnit: beginTx: %w", err)
	}
	defer tx.Rollback()

	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return fmt.Errorf("code: %w", ErrMissingField)
	}
	descr = strings.ToUpper(strings.TrimSpace(descr))
	if descr == "" {
		return fmt.Errorf("descr: %w", ErrMissingField)
	}
	name = strings.ToUpper(strings.TrimSpace(name))
	if name == "" {
		return fmt.Errorf("name: %w", ErrMissingField)
	}

	_, err = tx.ExecContext(s.ctx, "insert into units (code, descr, name) values (?, ?, ?)", code, descr, name)
	if err != nil {
		return fmt.Errorf("createUnit: insert: %w", err)
	}

	return tx.Commit()
}

func (s *Store) lookupUnitIdByCode(code string) int {
	if s.units == nil {
		units := make(map[string]*Unit)
		rows, err := s.db.Query("select id, code, name, descr from units")
		if err != nil {
			return 0
		}
		for rows.Next() {
			unit := &Unit{}
			err := rows.Scan(&unit.Id, &unit.Code, &unit.Name, &unit.Description)
			if err != nil {
				break
			}
			units[unit.Code] = unit
		}
		s.units = units
	}
	if u, ok := s.units[code]; ok {
		return u.Id
	}
	return 0
}
