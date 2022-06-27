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

import (
	"math"
	"math/rand"
)

// mkrings creates the rings used to generate systems.
// it limits the number of systems in each ring based on the systemsPerRing value.
// it removes the first three rings before returning the results.
func mkrings(radius, systemsPerRing int) [][]Coordinates {
	rings := make([][]Coordinates, radius+1, radius+1)

	// define the cartesian boundaries of the cluster
	minX, minY, minZ := -radius, -radius, -radius
	maxX, maxY, maxZ := radius, radius, radius

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				if d := int(math.Round(math.Sqrt(float64(x*x + y*y + z*z)))); d <= radius {
					// only allow stars in rings 3+
					if d > 2 {
						rings[d] = append(rings[d], Coordinates{X: x, Y: y, Z: z})
					}
				}
			}
		}
	}

	// shuffle the stars in each ring
	for d := 0; d <= radius; d++ {
		rand.Shuffle(len(rings[d]), func(i, j int) {
			rings[d][i], rings[d][j] = rings[d][j], rings[d][i]
		})
	}

	// limit the number of systems in each ring based on number of nations in game
	for d := 0; d <= radius; d++ {
		tots := systemsPerRing
		if d < 5 {
			tots += rand.Intn(10-d) + 1
		} else {
			tots += rand.Intn(d) + 1
		}
		if len(rings[d]) > tots {
			rings[d] = rings[d][0:tots]
		}
	}

	return rings

	//// return only the non-empty rings
	//return rings[3:]
}
