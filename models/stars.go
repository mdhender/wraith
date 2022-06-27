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

import "math/rand"

func (s *Store) genHomeStar(system *System) *Star {
	star := &Star{System: system, HomeStar: true, Orbits: make([]*Planet, 11, 11)}
	star.Kind = "A"
	star.Orbits[1] = s.genTerrestrial(star, 1)
	star.Orbits[2] = s.genTerrestrial(star, 2)
	star.Orbits[3] = s.genHomeTerrestrial(star, 3)
	star.Orbits[4] = s.genTerrestrial(star, 4)
	star.Orbits[5] = s.genAsteroidBelt(star, 5)
	star.Orbits[6] = s.genTerrestrial(star, 6)
	star.Orbits[7] = s.genGasGiant(star, 7)
	star.Orbits[8] = s.genGasGiant(star, 8)
	star.Orbits[9] = s.genTerrestrial(star, 9)
	star.Orbits[10] = s.genAsteroidBelt(star, 10)
	return star
}

func (s *Store) genStar(system *System) *Star {
	star := &Star{System: system, Orbits: make([]*Planet, 11, 11)}
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
	for orbit := range star.Orbits {
		if orbit == 0 {
			continue
		}
		switch rand.Intn(10) {
		case 0, 1, 2, 3:
			star.Orbits[orbit] = s.genTerrestrial(star, orbit)
		case 4, 5, 6:
			star.Orbits[orbit] = s.genGasGiant(star, orbit)
		case 7, 8:
			star.Orbits[orbit] = s.genAsteroidBelt(star, orbit)
		case 9:
			star.Orbits[orbit] = s.genEmpty(star, orbit)
		}
	}
	return star
}
