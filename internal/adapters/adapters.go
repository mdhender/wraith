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

// Package adapters implements functions to convert between data types
package adapters

import (
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/models"
)

func ModelsPlayerToEnginePlayer(mp *models.Player) *engine.Player {
	var ep engine.Player
	ep.Id = mp.Id
	ep.Nation.Name = mp.MemberOf.Details[0].Name
	ep.Nation.Speciality = mp.MemberOf.Speciality
	ep.Nation.HomeWorld = mp.MemberOf.HomePlanet.String()
	ep.Nation.GovtName = mp.MemberOf.Details[0].GovtName
	ep.Nation.GovtKind = mp.MemberOf.Details[0].GovtKind
	return &ep
}
