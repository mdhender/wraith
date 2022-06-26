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

import "math/rand"

type Star struct {
	Id       int
	System   *System // system the star is in
	Sequence string
	Kind     string
	HomeStar bool
	Orbits   []*Planet
}

func (e *Engine) genHomeStar(system *System) *Star {
	star := &Star{System: system, HomeStar: true}
	star.Kind = "A"
	star.Orbits = make([]*Planet, 11, 11)
	star.Orbits[1] = e.genTerrestrial(star, 1)
	star.Orbits[2] = e.genTerrestrial(star, 2)
	star.Orbits[3] = e.genHomeTerrestrial(star, 3)
	star.Orbits[4] = e.genTerrestrial(star, 4)
	star.Orbits[5] = e.genAsteroidBelt(star, 5)
	star.Orbits[6] = e.genTerrestrial(star, 6)
	star.Orbits[7] = e.genGasGiant(star, 7)
	star.Orbits[8] = e.genGasGiant(star, 8)
	star.Orbits[9] = e.genTerrestrial(star, 9)
	star.Orbits[10] = e.genAsteroidBelt(star, 10)
	return star
}

func (e *Engine) genStar(system *System) *Star {
	star := &Star{System: system}
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
	for orbit := range star.Orbits {
		if orbit == 0 {
			continue
		}
		switch rand.Intn(10) {
		case 0, 1, 2, 3:
			star.Orbits[orbit] = e.genTerrestrial(star, orbit)
		case 4, 5, 6:
			star.Orbits[orbit] = e.genGasGiant(star, orbit)
		case 7, 8:
			star.Orbits[orbit] = e.genAsteroidBelt(star, orbit)
		case 9:
			star.Orbits[orbit] = e.genEmpty(star, orbit)
		}
	}
	return star
}
