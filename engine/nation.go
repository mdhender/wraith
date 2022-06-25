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
	Id          int // this is the PK in the database
	No          int
	Name        string
	Description string
	Speciality  string
	Government  struct {
		Kind string
		Name string
	}
	HomePlanet struct {
		Name     string
		Location *Planet
	}
	TechLevel    int
	ResearchPool int
	Skills       Skills
	Colonies     []*XColony
	Ships        []*XShip
}

func (e *Engine) createNation(id int, planet *Planet) *Nation {
	n := &Nation{No: id, Speciality: "exploration"}
	n.Name = fmt.Sprintf("SP%d", n.No)
	n.Government.Name = fmt.Sprintf("GOV%d", n.No)
	n.Government.Kind = "monarchy"

	n.HomePlanet.Location = planet
	n.HomePlanet.Name = "Not Named"

	n.TechLevel, n.ResearchPool = 1, 0
	n.Skills.Biology = n.TechLevel
	n.Skills.Bureaucracy = n.TechLevel
	n.Skills.Gravitics = n.TechLevel
	n.Skills.LifeSupport = n.TechLevel
	n.Skills.Manufacturing = n.TechLevel
	n.Skills.Military = n.TechLevel
	n.Skills.Mining = n.TechLevel
	n.Skills.Shields = n.TechLevel

	colony := e.genHomeOpenColony(planet)
	n.Colonies = append(n.Colonies, colony)
	colony = e.genHomeOrbitalColony(planet)
	n.Colonies = append(n.Colonies, colony)

	return n
}
