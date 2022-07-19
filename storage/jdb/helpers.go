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

package jdb

import "fmt"

func unitAttributes(name string, techLevel int) (mets, nmts, totalMassUnits, fuelPerTurn, fuelPerCombatRound float64) {
	tl := float64(techLevel)
	switch name {
	case "anti-missile":
		return 2 * tl, 2 * tl, 4 * tl, 0, 0
	case "assault-craft":
		return 3 * tl, 2 * tl, 5 * tl, 0, 0.1
	case "assault-weapon":
		return 1 * tl, 1 * tl, 2 * tl, 2 * tl * tl, 0
	case "automation":
		return 2 * tl, 2 * tl, 4 * tl, 0, 0
	case "consumer-goods":
		return 0.2, 0.4, 0.6, 0, 0
	case "energy-shield":
		return 25 * tl, 25 * tl, 50 * tl, 0, 10 * tl
	case "energy-weapon":
		return 5 * tl, 5 * tl, 10 * tl, 0, 4 * tl
	case "factory":
		return 8 * tl, 4 * tl, 12 + 2*tl, 0.5 * tl, 4 * tl
	case "farm":
		if techLevel == 1 {
			return 4 + tl, 2 + tl, 6 + 2*tl, 0.5 * tl, 0
		} else if techLevel < 6 {
			return 4 + tl, 4 + tl, 6 + 2*tl, 0.5 * tl, 0
		}
		return 4 + tl, 2 + tl, 6 + 2*tl, tl, 0
	case "food":
		return 0, 0, 6, 0, 0
	case "fuel":
		return 0, 0, 1, 0, 0
	case "gold":
		return 0, 0, 1, 0, 0
	case "hyper-drive":
		return 25 * tl, 20 * tl, 45 * tl, 0, 0
	case "life-support":
		return 3 * tl, 5 * tl, 8 * tl, 1 * tl, 0
	case "light-structural":
		return 0.01, 0.04, 0.05, 0, 0
	case "metallics":
		return 0, 0, 1, 0, 0
	case "military-robots":
		return 10 * tl, 10 * tl, 20 + 2*tl, 0, 0
	case "military-supplies":
		return 0.02, 0.02, 0.04, 0, 0
	case "mine":
		return 5 + tl, 5 + tl, 10 + (2 * tl), 0.5 * tl, 0
	case "missile":
		return 2 * tl, 2 * tl, 4 * tl, 0, 0
	case "missile-launcher":
		return 15 * tl, 10 * tl, 25 * tl, 0, 0
	case "non-metallics":
		return 0, 0, 1, 0, 0
	case "sensor":
		return 10 * tl, 20 * tl, 40 * tl, tl / 20, 0
	case "space-drive":
		return 15 * tl, 10 * tl, 25 * tl, 0, tl * tl
	case "structural":
		return 0.1, 0.4, 0.5, 0, 0
	case "super-light-structural":
		return 0.001, 0.004, 0.005, 0, 0
	case "transport":
		return 3 * tl, tl, 4 * tl, 0.1 * tl * tl, 0.01 * tl * tl
	}
	panic(fmt.Sprintf("assert(unit.name != %q)", name))
}
