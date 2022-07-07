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

func (t *Turn) String() string {
	if t.Year == 0 && t.Quarter == 0 {
		return "0000/0"
	}
	return fmt.Sprintf("%04d/%d", t.Year, t.Quarter)
}

func (s *Store) FetchCurrentTurn(userHandle, gameName string) (*Turn, error) {
	row := s.db.QueryRow(`
		select no, year, quarter, start_dt, end_dt
		from games
			join turns on games.id = turns.game_id and turn = current_turn
			join users on users.handle = ?
			join player_dtl pd on
				users.id = pd.controlled_by
					and (pd.efftn <= games.current_turn and games.current_turn < pd.endtn)
		where games.short_name = ?`, userHandle, gameName)
	var t Turn
	err := row.Scan(&t.No, &t.Year, &t.Quarter, &t.StartDt, &t.EndDt)
	if err != nil {
		return nil, fmt.Errorf("fetchCurrentTurn: %q %q: %w", gameName, userHandle, err)
	}

	return &t, nil
}
