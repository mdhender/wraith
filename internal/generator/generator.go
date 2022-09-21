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

// Package generator implements a cluster generator.
package generator

import (
	"fmt"
	"math"
	"math/rand"
)

type System struct{}

func Generator(radius int, density int) (s []System) {
	for r := 1; r <= radius; r++ {
		R := float64(r)
		// z is in range -R ... R
		z := rand.Float64()
		// phi is in range 0 ... 2pi
		phi := rand.Float64() * 2 * math.Pi
		// theta is sin-1(z/R)
		theta := math.Asin(z / R)
		// x is R cos(theta) cos(phi)
		x := R * math.Cos(theta) * math.Cos(phi)
		// y is R cos(theta) sin(phi)
		y := R * math.Cos(theta) * math.Sin(phi)

		fmt.Println(x, y, z)
	}
	return s
}
