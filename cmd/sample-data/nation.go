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

package main

import "fmt"

type Nation struct {
	Id         int    `json:"nation-id"`
	Name       string `json:"name"`
	Speciality string `json:"speciality"`
	Government struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
	} `json:"government"`
	HomePlanet struct {
		Name     string `json:"name"`
		Location struct {
			X     int `json:"x"`
			Y     int `json:"y"`
			Z     int `json:"z"`
			Star  int `json:"star,omitempty"`
			Orbit int `json:"orbit"`
		} `json:"location"`
	} `json:"home-planet"`
	Skills   Skills    `json:"skills"`
	Colonies []*Colony `json:"colonies,omitempty"`
}

type Skills struct {
	Biology       int `json:"biology"`
	Bureaucracy   int `json:"bureaucracy"`
	Gravitics     int `json:"gravitics"`
	LifeSupport   int `json:"life-support"`
	Manufacturing int `json:"manufacturing"`
	Military      int `json:"military"`
	Mining        int `json:"mining"`
	Shields       int `json:"shields"`
}

func GenNation(id int) *Nation {
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

	n.Colonies = append(n.Colonies, GenHomeOpenColony(id))
	n.Colonies = append(n.Colonies, GenHomeOrbitalColony(id))

	return n
}
