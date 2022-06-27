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

func (s *Store) genNation(no int, planet *Planet, player *Player, position *PlayerPosition) *Nation {
	effTurn, endTurn := &Turn{}, &Turn{Year: 9999, Quarter: 4}

	n := &Nation{
		No:         no,
		HomePlanet: planet,
		Player:     player,
		Speciality: position.Nation.Speciality,
	}
	n.Details = []*NationDetail{{
		Nation:       n,
		EffTurn:      effTurn,
		EndTurn:      endTurn,
		Name:         position.Nation.Name,
		GovtKind:     position.Nation.GovtKind,
		GovtName:     position.Nation.GovtName,
		ControlledBy: player,
	}}
	n.Research = []*NationResearch{{
		Nation:             n,
		EffTurn:            effTurn,
		EndTurn:            endTurn,
		TechLevel:          1,
		ResearchPointsPool: 0,
	}}
	n.Skills = []*NationSkills{{
		Nation:        n,
		EffTurn:       effTurn,
		EndTurn:       endTurn,
		Biology:       1,
		Bureaucracy:   1,
		Gravitics:     1,
		LifeSupport:   1,
		Manufacturing: 1,
		Military:      1,
		Mining:        1,
		Shields:       1,
	}}

	n.Colonies = append(n.Colonies, s.genHomeOpenColony(1, planet, player))
	n.Colonies = append(n.Colonies, s.genHomeOrbitalColony(2, planet, player))

	return n
}
