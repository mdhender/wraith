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

import "fmt"

type Claim struct {
	User     string
	Player   string
	NationNo int
}

func (s *Store) FetchClaims(asOfTurn string) (map[string]*Claim, error) {
	if s.db == nil {
		return nil, ErrNoConnection
	}

	claims := make(map[string]*Claim)

	rows, err := s.db.Query(`
		select u.handle, pd.handle, n.nation_no
		from player_dtl pd
		left join users u on u.id = pd.controlled_by
		left join nation_player np on np.player_id = pd.player_id
		left join nations n on n.id = np.nation_id
		where (pd.efftn <= ? and ? < pd.endtn)`, asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchClaims: %q: %w", asOfTurn, err)
	}
	for rows.Next() {
		claim := &Claim{}
		err := rows.Scan(&claim.User, &claim.Player, &claim.NationNo)
		if err != nil {
			return nil, fmt.Errorf("fetchClaims: %q: %w", asOfTurn, err)
		}
		claims[claim.User] = claim
	}

	return claims, nil
}
