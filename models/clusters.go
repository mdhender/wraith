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
)

func (s *Store) FetchClusterListByGame(gameId int) ([]*SystemScan, error) {
	var scans []*SystemScan
	rows, err := s.db.Query(`select x, y, z, qty_stars from systems where game_id = ?`, gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchClusterListByGame: %d: %w", gameId, err)
	}
	for rows.Next() {
		var x, y, z, n int
		err := rows.Scan(&x, &y, &z, &n)
		if err != nil {
			return nil, fmt.Errorf("fetchClusterListByGame: %d: %w", gameId, err)
		}
		scans = append(scans, &SystemScan{X: x, Y: y, Z: z, QtyStars: n})
	}
	return scans, nil
}
