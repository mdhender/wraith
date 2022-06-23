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

package engine

import (
	"fmt"
)

// Nation configuration
type Nation struct {
	Id          int
	Name        string
	Description string
	Speciality  string
	Government  struct {
		Kind string
		Name string
	}
	HomePlanet struct {
		Name     string
		Location struct {
			X     int
			Y     int
			Z     int
			Orbit int
		}
	}
	Skills   Skills
	Colonies []*XColony
	Ships    []*XShip
}

func (e *Engine) createNation(id int) *Nation {
	n := &Nation{Id: id, Name: fmt.Sprintf("SP%d", id), Speciality: "exploration"}
	n.Government.Name = fmt.Sprintf("GOV%d", id)
	n.Government.Kind = "monarchy"

	n.HomePlanet.Name = "Home Planet"

	n.Skills.Biology = 1
	n.Skills.Bureaucracy = 1
	n.Skills.Gravitics = 1
	n.Skills.LifeSupport = 1
	n.Skills.Manufacturing = 1
	n.Skills.Military = 1
	n.Skills.Mining = 1
	n.Skills.Shields = 1

	return n
}
