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

package adapters

import (
	"github.com/mdhender/wraith/storage/jdb"
	"github.com/mdhender/wraith/wraith"
	"sort"
)

// WraithEngineToJdbGame converts an Engine to a Game.
func WraithEngineToJdbGame(e *wraith.Engine) *jdb.Game {
	jg := &jdb.Game{}
	jg.ShortName = e.Game.Code
	jg.Turn.Year = e.Game.Turn.Year
	jg.Turn.Quarter = e.Game.Turn.Quarter

	for _, player := range e.Players {
		p := &jdb.Player{
			Id:       player.Id,
			UserId:   player.UserId,
			Name:     player.Name,
			MemberOf: player.MemberOf.Id,
		}
		if player.ReportsTo != nil {
			p.ReportsToPlayerId = player.ReportsTo.Id
		}
		jg.Players = append(jg.Players, p)
	}
	sort.Sort(jg.Players)

	for _, nation := range e.Nations {
		n := &jdb.Nation{
			Id:                 nation.Id,
			No:                 nation.No,
			Name:               nation.Name,
			GovtName:           nation.GovtName,
			GovtKind:           nation.GovtKind,
			Speciality:         nation.Speciality,
			TechLevel:          nation.TechLevel,
			ResearchPointsPool: nation.ResearchPointsPool,
		}
		if nation.ControlledBy != nil {
			n.ControlledByPlayerId = nation.ControlledBy.Id
		}
		if nation.HomePlanet != nil {
			n.HomePlanetId = nation.HomePlanet.Id
		}
	}
	sort.Sort(jg.Nations)

	return jg
}
