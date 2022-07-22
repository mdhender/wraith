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
	"math"
	"strings"
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

func factoryProduction(cs *CorS, pos []*PhaseOrders) {
	cs.Log("Colony: %-10s   Kind: %-10s  Name: %s\n", cs.HullId, cs.Kind, cs.Name)
	cs.Log("  PRO %13d  SOL %13d  UNS %13d  FUEL %13d\n  UEM %13d  CON %13d  SPY %13d\n",
		availablePro(cs), availableSol(cs), availableUns(cs), availableFuel(cs),
		availableUem(cs), availableCon(cs), availableSpy(cs))

	for _, group := range cs.FactoryGroups {
		unitsProduced := 0
		for _, moe := range group.Units {
			unitsActive := maxCapacity(cs, moe)

			// allocate fuel
			moe.fuel.needed = moe.Unit.fuelUsed(moe.ActiveQty)
			moe.fuel.allocated = moe.Unit.fuelUsed(unitsActive)
			cs.fuel.available -= moe.fuel.allocated

			// allocate professional labor
			// TODO: add in factory efficiency
			moe.pro.needed = moe.ActiveQty
			moe.pro.allocated = unitsActive
			cs.pro.allocated += moe.pro.allocated

			// allocate unskilled labor
			// TODO: allow automation units to replace unskilled labor
			moe.uns.needed = 3 * moe.pro.needed
			moe.uns.allocated = 3 * unitsActive
			cs.uns.allocated += moe.uns.allocated

			cs.Log("  Group %2d: fuel %8d / %8d: pro %8d / %8d: uns %8d / %8d\n",
				group.No, moe.fuel.allocated, moe.fuel.needed, moe.pro.needed, moe.pro.allocated, moe.uns.allocated, moe.uns.allocated)

			// determine number of units produced
			if moe.Unit.TechLevel == 1 {
				unitsProduced = unitsActive * 100
			} else {
				unitsProduced = unitsActive * 20 * moe.Unit.TechLevel
			}
			// convert from units per year to units per turn
			unitsProduced = unitsProduced / 4
		}

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
		group.StageQty[0] += unitsProduced
		cs.Log("            25%%: %13d\n", group.StageQty[0])
		cs.Log("            50%%: %13d\n", group.StageQty[1])
		cs.Log("            75%%: %13d  finished: %13d %s\n", group.StageQty[2], group.StageQty[3], group.Product.Code)

		cs.Log("  PRO %13d  SOL %13d  UNS %13d  FUEL %13d\n  UEM %13d  CON %13d  SPY %13d\n",
			availablePro(cs), availableSol(cs), availableUns(cs), availableFuel(cs),
			availableUem(cs), availableCon(cs), availableSpy(cs))
	}
}

func farmProduction(cs *CorS, pos []*PhaseOrders) {
	cs.Log("Colony: %-10s   Kind: %-10s  Name: %s\n", cs.HullId, cs.Kind, cs.Name)
	cs.Log("  PRO %13d  SOL %13d  UNS %13d  FUEL %13d\n  UEM %13d  CON %13d  SPY %13d\n",
		availablePro(cs), availableSol(cs), availableUns(cs), availableFuel(cs),
		availableUem(cs), availableCon(cs), availableSpy(cs))

	for _, group := range cs.FarmGroups {
		unitsProduced := 0
		for _, moe := range group.Units {
			unitsActive := maxCapacity(cs, moe)

			// allocate fuel
			moe.fuel.needed = moe.Unit.fuelUsed(moe.ActiveQty)
			moe.fuel.allocated = moe.Unit.fuelUsed(unitsActive)
			cs.fuel.available -= moe.fuel.allocated

			// allocate professional labor
			moe.pro.needed = moe.ActiveQty
			moe.pro.allocated = unitsActive
			cs.pro.allocated += moe.pro.allocated

			// allocate unskilled labor
			// TODO: allow automation units to replace unskilled labor
			moe.uns.needed = 3 * moe.pro.needed
			moe.uns.allocated = 3 * unitsActive
			cs.uns.allocated += moe.uns.allocated

			cs.Log("  Group %2d: fuel %8d / %8d: pro %8d / %8d: uns %8d / %8d\n",
				group.No, moe.fuel.allocated, moe.fuel.needed, moe.pro.needed, moe.pro.allocated, moe.uns.allocated, moe.uns.allocated)

			// determine number of units produced
			if moe.Unit.TechLevel == 1 {
				unitsProduced = unitsActive * 100
			} else {
				unitsProduced = unitsActive * 20 * moe.Unit.TechLevel
			}
			// convert from units per year to units per turn
			unitsProduced = unitsProduced / 4
		}

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
		group.StageQty[0] += unitsProduced
		cs.Log("            25%%: %13d\n", group.StageQty[0])
		cs.Log("            50%%: %13d\n", group.StageQty[1])
		cs.Log("            75%%: %13d  finished: %13d %s\n", group.StageQty[2], group.StageQty[3], group.Product.Code)

		cs.Log("  PRO %13d  SOL %13d  UNS %13d  FUEL %13d\n  UEM %13d  CON %13d  SPY %13d\n",
			availablePro(cs), availableSol(cs), availableUns(cs), availableFuel(cs),
			availableUem(cs), availableCon(cs), availableSpy(cs))
	}
}

func (e *Engine) findColony(id string) (*CorS, bool) {
	c, ok := e.Colonies[id]
	return c, ok
}

func (e *Engine) findShip(id string) (*CorS, bool) {
	s, ok := e.Ships[id]
	return s, ok
}

func fuelInitialization(cs *CorS, pos []*PhaseOrders) {
	cs.Log("Colony: %-10s   Kind: %-10s  Name: %s\n", cs.HullId, cs.Kind, cs.Name)
	for _, u := range cs.Inventory {
		if u.Unit.Kind != "fuel" {
			continue
		}
		cs.fuel.available += u.ActiveQty + u.StowedQty
	}
	cs.Log("  %13d FUEL available for use\n\n", cs.fuel.available)
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

// killProportionally kills by proportion.
// deaths are booked as a reduction in the number of units available.
// tricky bit is dealing with crews and teams since they each have 2 population units
func killProportionally(cs *CorS, n int) {
	for n > 0 && totalPop(cs) > 0 {
		pct := float64(n) / float64(totalPop(cs))
		if cs.pro.available > 0 {
			k := int(pct * float64(cs.pro.available))
			n -= k
			cs.pro.available -= k
		}
		if cs.sol.available > 0 {
			k := int(pct * float64(cs.sol.available))
			if k < 1 {
				k = 1
			}
			n -= k
			cs.sol.available -= k
		}
		if cs.uns.available > 0 {
			k := int(pct * float64(cs.uns.available))
			if k < 1 {
				k = 1
			}
			n -= k
			cs.uns.available -= k
		}
		if cs.uem.available > 0 {
			k := int(pct * float64(cs.uem.available))
			if k < 1 {
				k = 1
			}
			n -= k
			cs.uem.available -= k
		}
		if cs.cons.available > 0 {
			k := int(pct * float64(cs.cons.available))
			if k < 1 {
				k = 1
			}
			n -= 2 * k
			cs.cons.available -= k
		}
		if cs.spy.available > 0 {
			k := int(pct * float64(cs.spy.available))
			if k < 1 {
				k = 1
			}
			n -= 2 * k
			cs.spy.available -= k
		}
	}
}

// laborInitialization updates the available labor pool.
// it allocates first to construction crews and then to spy teams.
// if there are not enough people to fill those crews or teams,
// we allocate as many as we can.
// the remaining population are added to the appropriate pools.
func laborInitialization(cs *CorS, pos []*PhaseOrders) {
	cs.Log("Colony: %-10s   Kind: %-10s  Name: %s\n", cs.HullId, cs.Kind, cs.Name)
	cs.Log("  PRO %13d  SOL %13d  UNS %13d  FUEL %13d\n  UEM %13d  CON %13d  SPY %13d\n",
		cs.Population.ProfessionalQty, cs.Population.SoldierQty, cs.Population.UnskilledQty, cs.fuel.available,
		cs.Population.UnemployedQty, cs.Population.ConstructionCrewQty, cs.Population.SpyTeamQty)

	cs.pro.available = cs.Population.ProfessionalQty
	cs.sol.available = cs.Population.SoldierQty
	cs.uns.available = cs.Population.UnskilledQty
	cs.uem.available = cs.Population.UnemployedQty

	// TODO: let automation units replace unskilled workers
	cs.cons.available = cs.Population.ConstructionCrewQty
	if cs.cons.available > 0 {
		if availablePro(cs) < cs.cons.available {
			cs.cons.available = availablePro(cs)
		}
		if availableUns(cs) < cs.cons.available {
			cs.cons.available = availableUns(cs)
		}
		cs.pro.allocated += cs.cons.available
		cs.uns.allocated += cs.cons.available
	}

	cs.spy.available = cs.Population.SpyTeamQty
	if cs.spy.available > 0 {
		if availablePro(cs) < cs.spy.available {
			cs.spy.available = availablePro(cs)
		}
		if availableSol(cs) < cs.spy.available {
			cs.spy.available = availableSol(cs)
		}
		cs.pro.allocated += cs.spy.available
		cs.sol.allocated += cs.spy.available
	}

	cs.Log("  PRO %13d  SOL %13d  UNS %13d  FUEL %13d\n  UEM %13d  CON %13d  SPY %13d\n",
		availablePro(cs), availableSol(cs), availableUns(cs), availableFuel(cs),
		availableUem(cs), availableCon(cs), availableSpy(cs))

	return
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
	if maxUnits > availablePro(cs) {
		maxUnits = availablePro(cs)
	}

	// limit capacity based on available unskilled workers
	if maxUnits > availableUns(cs)/3 {
		maxUnits = availableUns(cs) / 3
	}

	return maxUnits
}

func mineProduction(cs *CorS, pos []*PhaseOrders) {
	cs.Log("Colony: %-10s   Kind: %-10s  Name: %s\n", cs.HullId, cs.Kind, cs.Name)
	cs.Log("  PRO %13d  SOL %13d  UNS %13d  FUEL %13d\n  UEM %13d  CON %13d  SPY %13d\n",
		availablePro(cs), availableSol(cs), availableUns(cs), availableFuel(cs),
		availableUem(cs), availableCon(cs), availableSpy(cs))

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
		cs.pro.allocated += moe.pro.allocated

		// allocate unskilled labor
		// TODO: allow automation units to replace unskilled labor
		moe.uns.needed = 3 * moe.pro.needed
		moe.uns.allocated = 3 * unitsActive
		cs.uns.allocated += moe.uns.allocated

		cs.Log("  Group %2d: %-6s      yield: %7.3f%%     reserves: %13d tonnes\n",
			group.No, group.Deposit.Product.Code, 100*group.Deposit.YieldPct, group.Deposit.RemainingQty)
		cs.Log("            fuel %8d / %8d: pro %8d / %8d: uns %8d / %8d\n",
			moe.fuel.allocated, moe.fuel.needed, moe.pro.needed, moe.pro.allocated, moe.uns.allocated, moe.uns.allocated)

		// determine number of units produced
		unitsProduced = unitsActive * 100 * moe.Unit.TechLevel
		// convert from units per year to units per turn
		unitsProduced = unitsProduced / 4

		// push the newly produced units through the pipeline
		if group.StageQty[2] > unitsProduced {
			group.StageQty[3] = int(math.Ceil(float64(unitsProduced) * group.Deposit.YieldPct))
			group.StageQty[2] -= unitsProduced
		} else {
			group.StageQty[3] = int(math.Ceil(float64(group.StageQty[2]) * group.Deposit.YieldPct))
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
		group.StageQty[0] += unitsProduced
		cs.Log("            25%%: %13d       50%%: %13d\n", group.StageQty[0], group.StageQty[1])
		cs.Log("            75%%: %13d  finished: %13d %s\n", group.StageQty[2], group.StageQty[3], group.Deposit.Product.Code)

		cs.Log("  PRO %13d  SOL %13d  UNS %13d  FUEL %13d\n  UEM %13d  CON %13d  SPY %13d\n",
			availablePro(cs), availableSol(cs), availableUns(cs), availableFuel(cs),
			availableUem(cs), availableCon(cs), availableSpy(cs))
	}
}

func totalPop(cs *CorS) int {
	return cs.pro.available + cs.sol.available + cs.uns.available + cs.uem.available + 2*cs.cons.available + 2*cs.spy.available
}
func availableCon(cs *CorS) int {
	return cs.cons.available - cs.cons.allocated
}
func availableFuel(cs *CorS) int {
	return cs.fuel.available - cs.fuel.allocated
}
func availablePro(cs *CorS) int {
	return cs.pro.available - cs.pro.allocated
}
func availableSol(cs *CorS) int {
	return cs.sol.available - cs.sol.allocated
}
func availableSpy(cs *CorS) int {
	return cs.spy.available - cs.spy.allocated
}
func availableUem(cs *CorS) int {
	return cs.uem.available - cs.uem.allocated
}
func availableUns(cs *CorS) int {
	return cs.uns.available - cs.uns.allocated
}

func unitFromString(e *Engine, s string) (*Unit, bool) {
	// try code first
	u, ok := e.UnitsFromString[strings.ToUpper(s)]
	if !ok {
		// try long name
		u, ok = e.UnitsFromString[strings.ToLower(s)]
	}
	return u, ok
}
