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

type System struct {
	Id          int
	X           int
	Y           int
	Z           int
	HomeSystem  bool
	Singularity bool
	Ring        int
	Stars       []*Star
}

func (e *Engine) genSystem(id int) *System {
	system := &System{Id: id}
	switch rand.Intn(21) {
	case 0, 1, 2, 3, 4, 5:
		system.Stars = make([]*Star, 1, 1)
	case 6, 7, 8, 9, 10:
		system.Stars = make([]*Star, 2, 2)
	case 11, 12, 13, 14:
		system.Stars = make([]*Star, 3, 3)
	case 15, 16, 17:
		system.Stars = make([]*Star, 4, 4)
	case 18, 19:
		system.Stars = make([]*Star, 5, 5)
	case 20:
		system.Stars = make([]*Star, 6, 6)
	}
	for i := range system.Stars {
		star := e.genStar(system)
		star.Sequence = string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i])
		system.Stars[i] = star
	}

	return system
}

func (e *Engine) genHomeSystem(id int) *System {
	system := &System{Id: id, HomeSystem: true}
	system.Stars = make([]*Star, 1, 1)
	for i := range system.Stars {
		var star *Star
		if i == 0 {
			star = e.genHomeStar(system)
		} else {
			star = e.genStar(system)
		}
		star.Sequence = string("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i])
		system.Stars[i] = star
	}
	return system
}
