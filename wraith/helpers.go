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

package wraith

import (
	"log"
	"math"
)

type ClusterList []*ClusterListItem

type ClusterListItem struct {
	X, Y, Z  int
	QtyStars int
	Distance float64
}

// Len implements the sort.Sort interface
func (u ClusterList) Len() int {
	return len(u)
}

// Less implements the sort.Sort interface
func (u ClusterList) Less(i, j int) bool {
	if u[i].Distance < u[j].Distance {
		return true
	} else if u[i].Distance > u[j].Distance {
		return false
	}

	if u[i].X < u[j].X {
		return true
	} else if u[i].X > u[j].X {
		return false
	}

	if u[i].Y < u[j].Y {
		return true
	} else if u[i].Y > u[j].Y {
		return false
	}

	return u[i].Z < u[j].Z
}

// Swap implements the sort.Sort interface
func (u ClusterList) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (e *Engine) ClusterScan(origin Coordinates) ClusterList {
	var cl ClusterList
	for _, system := range e.Systems {
		dx, dy, dz := system.Coords.X-origin.X, system.Coords.Y-origin.Y, system.Coords.Z-origin.Z
		cl = append(cl, &ClusterListItem{
			X:        system.Coords.X,
			Y:        system.Coords.Y,
			Z:        system.Coords.Z,
			QtyStars: len(system.Stars),
			Distance: math.Sqrt(float64(dx*dx + dy*dy + dz*dz)),
		})
	}
	return cl
}

func (e *Engine) findColony(id string) (*CorS, bool) {
	c, ok := e.Colonies[id]
	return c, ok
}

func (e *Engine) findShip(id string) (*CorS, bool) {
	s, ok := e.Ships[id]
	return s, ok
}

func indexOf(s string, sl []string) int {
	for i, p := range sl {
		if s == p {
			return i
		}
	}
	return -1
}

func isSolarPowered(u *Unit, cs *CorS) bool {
	switch u.Kind {
	case "farm":
		return cs.Planet.OrbitNo <= 5 && cs.Kind == "orbital" && ("FRM-2" <= u.Code && u.Code <= "FRM-5")
	}
	return false
}

func isZero(n float64) bool {
	return n < 0.001
}

// TODO: logic for factory groups efficiency, automation units, and construction crews
func maxCapacity(cs *CorS, u *InventoryUnit) int {
	// assume maximum capacity
	maxUnits := u.ActiveQty

	// limit capacity based on available fuel
	if !(isSolarPowered(u.Unit, cs) || isZero(u.Unit.FuelPerUnitPerTurn)) {
		fuelCapacity := int(float64(cs.fuel.available) / u.Unit.FuelPerUnitPerTurn)
		if maxUnits > fuelCapacity {
			maxUnits = fuelCapacity
		}
	}

	// limit capacity based on available professionals
	if maxUnits > cs.population.professional {
		maxUnits = cs.population.professional
	}

	// limit capacity based on available unskilled workers
	if maxUnits > cs.population.unskilled/3 {
		maxUnits = cs.population.unskilled / 3
	}

	return maxUnits
}

func mineProduction(cs *CorS, pos []*PhaseOrders) {
	playerId := 0
	var playerName string
	if cs.ControlledBy != nil {
		playerId, playerName = cs.ControlledBy.Id, cs.ControlledBy.Name
	} else {
		playerName = "nobody"
	}

	for _, group := range cs.MineGroups {
		unitsProduced := 0
		moe := group.Unit
		unitsActive := maxCapacity(cs, moe)

		// allocate fuel
		moe.fuel.needed = moe.Unit.fuelUsed(moe.ActiveQty)
		moe.fuel.allocated = moe.Unit.fuelUsed(unitsActive)
		cs.fuel.available -= moe.fuel.allocated

		// allocate professional labor
		moe.pro.needed = moe.ActiveQty
		moe.pro.allocated = unitsActive
		cs.population.professional -= moe.pro.allocated

		// allocate unskilled labor
		// TODO: allow automation units to replace unskilled labor
		moe.uns.needed = 3 * moe.pro.needed
		moe.uns.allocated = 3 * unitsActive
		cs.population.unskilled -= moe.uns.allocated

		log.Printf("cors: mineProduction: %2d: %-20s: %-6s: group %2d: fuel %8d / %8d: pro %8d / %8d: uns %8d / %8d\n",
			playerId, playerName, cs.HullId, group.No, moe.fuel.allocated, moe.fuel.needed, moe.pro.needed, moe.pro.allocated, moe.uns.allocated, moe.uns.allocated)

		// determine number of units produced
		unitsProduced = unitsActive * 100 * moe.Unit.TechLevel
		// convert from units per year to units per turn
		unitsProduced = unitsProduced / 4

		// push the newly produced units through the pipeline
		if group.StageQty[2] > unitsProduced {
			group.StageQty[3] = unitsProduced
			group.StageQty[2] -= unitsProduced
		} else {
			group.StageQty[3] = group.StageQty[2]
			group.StageQty[2] = 0
		}
		if group.StageQty[1] > unitsProduced {
			group.StageQty[2] += unitsProduced
			group.StageQty[1] -= unitsProduced
		} else {
			group.StageQty[2] += group.StageQty[1]
			group.StageQty[1] = 0
		}
		if group.StageQty[0] > unitsProduced {
			group.StageQty[1] += unitsProduced
			group.StageQty[0] -= unitsProduced
		} else {
			group.StageQty[1] += group.StageQty[0]
			group.StageQty[0] = 0
		}
		group.Deposit.RemainingQty -= unitsProduced
		group.StageQty[0] += int(math.Ceil(float64(unitsProduced) * group.Deposit.YieldPct))
	}
}
