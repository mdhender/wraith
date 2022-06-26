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
	ControlledBy *Player
}

func (e *Engine) createNation(id int, planet *Planet, player *Player) *Nation {
	n := &Nation{No: id}
	n.Name = player.Nation.Name
	n.Government.Kind = player.Nation.GovtKind
	n.Government.Name = player.Nation.GovtName
	n.HomePlanet.Location = planet
	n.HomePlanet.Name = player.Nation.HomeWorld
	n.ResearchPool = 0
	n.Skills.Biology = 1
	n.Skills.Bureaucracy = 1
	n.Skills.Gravitics = 1
	n.Skills.LifeSupport = 1
	n.Skills.Manufacturing = 1
	n.Skills.Military = 1
	n.Skills.Mining = 1
	n.Skills.Shields = 1
	n.Speciality = player.Nation.Speciality
	n.TechLevel = 1
	n.ControlledBy = player

	colony := e.genHomeOpenColony(planet)
	n.Colonies = append(n.Colonies, colony)
	colony = e.genHomeOrbitalColony(planet)
	n.Colonies = append(n.Colonies, colony)

	return n
}
