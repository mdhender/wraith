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

import "math/rand"

type Star struct {
	Id       int       `json:"star-id,omitempty"`
	Kind     string    `json:"kind,omitempty"`
	HomeStar bool      `json:"home-star,omitempty"`
	Orbits   []*Planet `json:"orbits,omitempty"`
}

var numPlanets int

func GenHomeStar(systemId, id int) *Star {
	star := &Star{Id: id, HomeStar: true}
	star.Kind = "A"
	star.Orbits = make([]*Planet, 11, 11)
	numPlanets++
	star.Orbits[1] = GenTerrestrial(numPlanets, 1)
	numPlanets++
	star.Orbits[2] = GenTerrestrial(numPlanets, 2)
	numPlanets++
	star.Orbits[3] = GenHomeTerrestrial(numPlanets, 3)
	numPlanets++
	star.Orbits[4] = GenTerrestrial(numPlanets, 4)
	numPlanets++
	star.Orbits[5] = GenAsteroidBelt(numPlanets, 5)
	numPlanets++
	star.Orbits[6] = GenTerrestrial(numPlanets, 6)
	numPlanets++
	star.Orbits[7] = GenGasGiant(numPlanets, 7)
	numPlanets++
	star.Orbits[8] = GenGasGiant(numPlanets, 8)
	numPlanets++
	star.Orbits[9] = GenTerrestrial(numPlanets, 9)
	numPlanets++
	star.Orbits[10] = GenAsteroidBelt(numPlanets, 10)
	return star
}

func GenStar(systemId, id int) *Star {
	star := &Star{Id: id}
	switch rand.Intn(14) {
	case 0, 1, 2, 3, 4, 5, 6:
		star.Kind = "A"
	case 7, 8, 9, 10:
		star.Kind = "B"
	case 11, 12:
		star.Kind = "C"
	case 13:
		star.Kind = "D"
	}
	star.Orbits = make([]*Planet, 11, 11)
	for i := range star.Orbits {
		if i == 0 {
			continue
		}
		numPlanets++
		switch rand.Intn(10) {
		case 0, 1, 2, 3:
			star.Orbits[i] = GenTerrestrial(numPlanets, i)
		case 4, 5, 6:
			star.Orbits[i] = GenGasGiant(numPlanets, i)
		case 7, 8:
			star.Orbits[i] = GenAsteroidBelt(numPlanets, i)
		case 9:
			star.Orbits[i] = GenEmpty(numPlanets, i)
		}
	}
	return star
}
