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
	"math"
	"math/rand"
)

type Cluster struct {
	Radius  int       `json:"radius"`
	Systems []*System `json:"systems"`
}

func GenCluster(radius int, rings [][]Coordinates, nations []*Nation) *Cluster {
	cluster := Cluster{Radius: radius}

	systemId := 0
	for i := 0; i < len(nations); i++ {
		systemId++
		r := 5
		coords := rings[r][0]
		rings[r] = rings[r][1:]
		system := GenHomeSystem(systemId)
		system.Ring, system.X, system.Y, system.Z = r, coords.X, coords.Y, coords.Z
		cluster.Systems = append(cluster.Systems, system)

		nations[i].HomePlanet.Location.X = coords.X
		nations[i].HomePlanet.Location.Y = coords.Y
		nations[i].HomePlanet.Location.Z = coords.Z
		nations[i].HomePlanet.Location.Star = 0
		nations[i].HomePlanet.Location.Orbit = 3

		nations[i].Colonies[0].Location.X = coords.X
		nations[i].Colonies[0].Location.Y = coords.Y
		nations[i].Colonies[0].Location.Z = coords.Z
		nations[i].Colonies[0].Location.Star = 0
		nations[i].Colonies[0].Location.Orbit = 3

		nations[i].Colonies[1].Location.X = coords.X
		nations[i].Colonies[1].Location.Y = coords.Y
		nations[i].Colonies[1].Location.Z = coords.Z
		nations[i].Colonies[1].Location.Star = 0
		nations[i].Colonies[1].Location.Orbit = 3
	}

	for r := 0; r < len(rings); r++ {
		for _, coords := range rings[r] {
			systemId++
			system := GenSystem(systemId)
			system.Ring, system.X, system.Y, system.Z = r, coords.X, coords.Y, coords.Z
			cluster.Systems = append(cluster.Systems, system)
		}
	}

	return &cluster
}

func (c *Cluster) randomXYZ() (int, int, int) {
	radius, points := float64(c.Radius), 2*c.Radius+1
	for {
		x, y, z := rand.Intn(points)-c.Radius, rand.Intn(points)-c.Radius, rand.Intn(points)-c.Radius
		d := math.Sqrt(float64(x*x + y*y + z*z))
		if 2 <= d && d <= radius {
			dup := false
			for _, s := range c.Systems {
				if x == s.X && y == s.Y && z == s.Z {
					dup = true
					break
				}
			}
			if !dup {
				return x, y, z
			}
		}
	}
}
