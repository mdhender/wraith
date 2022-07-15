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
	"time"
)

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

type UserClaim struct {
	Id     int       // user id
	Handle string    // user handle
	AsOf   time.Time // as of date for pulling claims
	Games  []*GameClaim
}

type GameClaim struct {
	Id           int    // game id
	ShortName    string // game code
	EffTurn      *Turn  // game's turn that is effective as of the claim date
	PlayerId     int    // player id
	PlayerHandle string // player handle in game
	NationId     int    // nation id
	NationNo     int    // nation player is aligned with
	NationName   string // nation government name
}

func (s *Store) FetchUserClaimAsOf(id int, asOf time.Time) (*UserClaim, error) {
	if s.db == nil {
		return nil, ErrNoConnection
	}

	claim := &UserClaim{AsOf: asOf}

	rows, err := s.db.Query(`
		select u.id, up.handle,
			   g.id, g.short_name,
			   t.turn,
			   p.id, pd.handle,
			   n.id, n.nation_no,
			   nd.govt_name
		from users u
			inner join user_profile up
				on u.id = up.user_id
					   and (up.effdt <= ? and ? < up.enddt)
			inner join player_dtl pd
				on u.id = pd.controlled_by
			inner join players p
				on p.id = pd.player_id
			inner join games g
				on p.game_id = g.id
			inner join turns t
				on g.id = t.game_id
					   and g.current_turn = t.turn
			inner join nations n
				on g.id = n.game_id
			inner join nation_dtl nd
				on n.id = nd.nation_id
					   and p.id = nd.controlled_by
					   and (nd.efftn <= t.turn and t.turn < nd.endtn)
		where u.id = ?`, asOf, asOf, id)
	if err != nil {
		return nil, fmt.Errorf("fetchUserClaimAsOf: %q: %w", asOf, err)
	}
	for rows.Next() {
		var currentTurn string
		gc := &GameClaim{}
		err := rows.Scan(&claim.Id, &claim.Handle, &gc.Id, &gc.ShortName, &currentTurn, &gc.PlayerId, &gc.PlayerHandle, &gc.NationId, &gc.NationNo, &gc.NationName)
		if err != nil {
			return nil, fmt.Errorf("fetchUserClaimAsOf: %q: %w", asOf, err)
		}
		gc.EffTurn, err = s.fetchTurn(gc.Id, currentTurn)
		if err != nil {
			return nil, fmt.Errorf("fetchUserClaimAsOf: %q: %d %q: %w", asOf, gc.Id, gc.ShortName, err)
		}
		claim.Games = append(claim.Games, gc)
	}

	return claim, nil
}

func (s *Store) FetchUserClaimsFromGameAsOf(id int, asOf time.Time) ([]*UserClaim, error) {
	if s.db == nil {
		return nil, ErrNoConnection
	}

	var claims []*UserClaim

	rows, err := s.db.Query(`
		select u.id,
			   up.handle,
			   g.id, g.short_name,
			   t.turn,
			   p.id, pd.handle,
			   n.id, n.nation_no,
			   nd.govt_name
		from games g
			inner join turns t
				on g.id = t.game_id
					   and g.current_turn = t.turn
			inner join players p
				on g.id = p.game_id
			inner join player_dtl pd
				on p.id = pd.player_id
					   and (pd.efftn <= t.turn and t.turn < pd.endtn)
			inner join users u
				on pd.controlled_by = u.id
			inner join user_profile up
				on u.id = up.user_id
					   and (up.effdt <= ? and ? < up.enddt)
			inner join nations n
				on g.id = n.game_id
			inner join nation_dtl nd
				on n.id = nd.nation_id
					   and p.id = nd.controlled_by
					   and (nd.efftn <= t.turn and t.turn < nd.endtn)
		where g.id = ?
		order by n.nation_no, pd.handle`, asOf, asOf, id)
	if err != nil {
		return nil, fmt.Errorf("fetchUserClaimsFromGameAsOf: %d: %q: %w", id, asOf, err)
	}
	var turn *Turn
	for rows.Next() {
		user := &UserClaim{AsOf: asOf}
		var currentTurn string
		gc := &GameClaim{}
		err := rows.Scan(&user.Id, &user.Handle, &gc.Id, &gc.ShortName, &currentTurn, &gc.PlayerId, &gc.PlayerHandle, &gc.NationId, &gc.NationNo, &gc.NationName)
		if err != nil {
			return nil, fmt.Errorf("fetchUserClaimsFromGameAsOf: %d: %q: %w", id, asOf, err)
		}
		if turn == nil {
			turn, err = s.fetchTurn(id, currentTurn)
			if err != nil {
				return nil, fmt.Errorf("fetchUserClaimsFromGameAsOf: %d: %q: %w", id, asOf, err)
			}
		}
		gc.EffTurn = turn
		user.Games = append(user.Games, gc)
		claims = append(claims, user)
	}

	return claims, nil
}
