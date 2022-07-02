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
	"bytes"
	"fmt"
)

func (s *Store) FetchClusterByGame(gameId, x, y, z int) ([]byte, error) {
	rows, err := s.db.Query(`
		select s.x, s.y, s.z,
			   round(sqrt((o.x - s.x)*(o.x - s.x)+(o.y - s.y)*(o.y - s.y)+(o.z - s.z)*(o.z - s.z)), 3) as distance,
			   s.qty_stars
		from systems o, systems s
		where (o.game_id = ? and o.x = ? and o.y = ? and o.z = ?)
		  and s.game_id = o.game_id
		order by distance, s.x, s.y, s.z`, gameId, x, y, z)
	if err != nil {
		return nil, fmt.Errorf("fetchClusterByGame: %d: %w", gameId, err)
	}
	bb := bytes.NewBuffer(nil)
	_, _ = bb.WriteString("<table border=\"1\">")
	_, _ = bb.WriteString("<thead><tr><td>Distance from Origin</td><td>X</td><td>Y</td><td>Z</td><td>Number of Stars</td></thead>")
	for rows.Next() {
		var xx, yy, zz, ns int
		var distance float64
		err := rows.Scan(&xx, &yy, &zz, &distance, &ns)
		if err != nil {
			return nil, fmt.Errorf("fetchClusterByGame: %d: %w", gameId, err)
		}
		_, _ = bb.WriteString(fmt.Sprintf("<tr><td>%12.3f</td><td align=\"right\">&nbsp;&nbsp;%d&nbsp;</td><td align=\"right\">&nbsp;&nbsp;%d&nbsp;</td><td align=\"right\">&nbsp;&nbsp;%d&nbsp;</td><td align=\"right\">&nbsp;%d&nbsp;</td></tr>",
			distance, xx, yy, zz, ns))
	}
	_, _ = bb.WriteString("</table>")
	return bb.Bytes(), nil
}

func (s *Store) FetchClusterByGameOrigin(gameId int) ([]byte, error) {
	rows, err := s.db.Query(`
		select x, y, z,
			   sqrt(x*x+y*y+z*z) as distance,
			   qty_stars
		from systems
		where game_id = ?
		order by distance, x, y, z`, gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchClusterByGameOrigin: %d: %w", gameId, err)
	}
	bb := bytes.NewBuffer(nil)
	_, _ = bb.WriteString("<table border=\"1\">")
	_, _ = bb.WriteString("<thead><tr><td>Distance from Origin</td><td>X</td><td>Y</td><td>Z</td><td>Number of Stars</td></thead>")
	for rows.Next() {
		var xx, yy, zz, ns int
		var distance float64
		err := rows.Scan(&xx, &yy, &zz, &distance, &ns)
		if err != nil {
			return nil, fmt.Errorf("fetchClusterByGameOrigin: %d: %w", gameId, err)
		}
		_, _ = bb.WriteString(fmt.Sprintf("<tr><td>%12.3f</td><td align=\"right\">&nbsp;&nbsp;%d&nbsp;</td><td align=\"right\">&nbsp;&nbsp;%d&nbsp;</td><td align=\"right\">&nbsp;&nbsp;%d&nbsp;</td><td align=\"right\">&nbsp;%d&nbsp;</td></tr>",
			distance, xx, yy, zz, ns))
	}
	_, _ = bb.WriteString("</table>")
	return bb.Bytes(), nil
}
