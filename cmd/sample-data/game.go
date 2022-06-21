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

import (
	"fmt"
	"math"
	"math/rand"
)

type Game struct {
	Id      string    `json:"id"`
	Turn    int       `json:"turn"`
	Nations []*Nation `json:"nations,omitempty"`
	Cluster *Cluster  `json:"cluster,omitempty"`
}

type Coordinates struct {
	X, Y, Z int
}

func GenGame(numberOfNations, radius int) *Game {
	systemsPerRing := numberOfNations
	totalSystems := radius * systemsPerRing
	fmt.Printf("totalSystems %3d %6d\n", systemsPerRing, totalSystems)

	g := &Game{Id: "PT-1", Turn: 0}

	for i := 0; i < numberOfNations; i++ {
		g.Nations = append(g.Nations, GenNation(i+1))
	}

	g.Cluster = GenCluster(radius, mkrings(radius, systemsPerRing), g.Nations)

	return g
}

func mkrings(radius, systemsPerRing int) [][]Coordinates {
	// generate rings to use for distributing stars in a much better version of this program
	minX, minY, minZ := -radius, -radius, -radius
	maxX, maxY, maxZ := radius, radius, radius
	rings := make([][]Coordinates, radius+1, radius+1)
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				if d := int(math.Round(math.Sqrt(float64(x*x + y*y + z*z)))); d <= radius {
					rings[d] = append(rings[d], Coordinates{X: x, Y: y, Z: z})
				}
			}
		}
	}

	//numPoints := 0
	//for d := 0; d <= radius; d++ {
	//	numPoints += len(rings[d])
	//	fmt.Printf("ring %2d: %5d\n", d, len(rings[d]))
	//}
	//fmt.Printf("  total: %5d\n", numPoints)

	// prevent systems from being created in the first three rings
	rings[0], rings[1], rings[2] = nil, nil, nil

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

	numPoints := 0
	for d := 0; d <= radius; d++ {
		numPoints += len(rings[d])
		fmt.Printf("ring %2d: %5d\n", d, len(rings[d]))
	}
	fmt.Printf("  total: %5d\n", numPoints)

	return rings
}
