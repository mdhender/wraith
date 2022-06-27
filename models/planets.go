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

func (s *Store) genAsteroidBelt(star *Star, orbit int) *Planet {
	efftn, endtn := &Turn{}, &Turn{Year: 9999, Quarter: 4}

	planet := &Planet{Star: star, Kind: "asteroid-belt", OrbitNo: orbit}
	planet.Details = []*PlanetDetail{{Planet: planet, EffTurn: efftn, EndTurn: endtn}}

	for r := 0; r <= rand.Intn(40); r++ {
		nr := &NaturalResource{Planet: planet}
		nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn}}

		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 11, 12, 13, 14, 15, 16, 17:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 18, 19:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
		case 20:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "gold", 0.01+float64(rand.Intn(5))/100, rand.Intn(30)*100_000
		}
		if nr.QtyInitial < 100_000 {
			nr.QtyInitial = 100_000
		}
		nr.Details[0].QtyRemaining = nr.QtyInitial

		planet.Deposits = append(planet.Deposits, nr)
	}

	return planet
}

func (s *Store) genEmpty(star *Star, orbit int) *Planet {
	planet := &Planet{Star: star, Kind: "empty", OrbitNo: orbit}

	return planet
}

func (s *Store) genGasGiant(star *Star, orbit int) *Planet {
	efftn, endtn := &Turn{}, &Turn{Year: 9999, Quarter: 4}

	planet := &Planet{Star: star, Kind: "gas-giant", OrbitNo: orbit}
	planet.Details = []*PlanetDetail{{Planet: planet, EffTurn: efftn, EndTurn: endtn}}

	if 3 <= orbit && orbit <= 5 {
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5:
			planet.HabitabilityNo = rand.Intn(1)
		case 6, 7, 8, 9, 10:
			planet.HabitabilityNo = rand.Intn(1) + rand.Intn(1)
		case 11, 12, 13, 14:
			planet.HabitabilityNo = rand.Intn(2) + rand.Intn(1) + rand.Intn(1)
		case 15, 16, 17:
			planet.HabitabilityNo = rand.Intn(2) + rand.Intn(2) + rand.Intn(1) + rand.Intn(1)
		case 18, 19:
			planet.HabitabilityNo = rand.Intn(3) + rand.Intn(2) + rand.Intn(2) + rand.Intn(1) + rand.Intn(1)
		case 20:
			planet.HabitabilityNo = rand.Intn(3) + rand.Intn(3) + rand.Intn(2) + rand.Intn(2) + rand.Intn(1) + rand.Intn(1)
		}
	}

	for r := 0; r <= rand.Intn(40); r++ {
		nr := &NaturalResource{Planet: planet}
		nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn}}

		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 16, 17, 18, 19:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 20:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
		}
		if nr.QtyInitial < 100_000 {
			nr.QtyInitial = 100_000
		}
		nr.Details[0].QtyRemaining = nr.QtyInitial

		planet.Deposits = append(planet.Deposits, nr)
	}

	return planet
}

func (s *Store) genHomeTerrestrial(star *Star, orbit int) *Planet {
	efftn, endtn := &Turn{}, &Turn{Year: 9999, Quarter: 4}

	planet := &Planet{Star: star, Kind: "terrestrial", OrbitNo: orbit, HomePlanet: true, HabitabilityNo: 25}
	planet.Details = []*PlanetDetail{{Planet: planet, EffTurn: efftn, EndTurn: endtn}}

	nr := &NaturalResource{Planet: planet, Kind: "gold", YieldPct: 0.07, QtyInitial: 300_000}
	nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn, QtyRemaining: nr.QtyInitial}}
	planet.Deposits = append(planet.Deposits, nr)

	nr = &NaturalResource{Planet: planet, Kind: "fuel", YieldPct: 0.25, QtyInitial: 99_999_999}
	nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn, QtyRemaining: nr.QtyInitial}}
	planet.Deposits = append(planet.Deposits, nr)

	nr = &NaturalResource{Planet: planet, Kind: "non-metallic", YieldPct: 0.25, QtyInitial: 99_999_999}
	nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn, QtyRemaining: nr.QtyInitial}}
	planet.Deposits = append(planet.Deposits, nr)

	nr = &NaturalResource{Planet: planet, Kind: "metallic", YieldPct: 0.25, QtyInitial: 99_999_999}
	nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn, QtyRemaining: nr.QtyInitial}}
	planet.Deposits = append(planet.Deposits, nr)

	yield, qty := 0.45, 90_000_000
	for i := len(planet.Deposits); i < 8; i++ {
		yield, qty = yield*0.9, qty*8/10

		nr := &NaturalResource{Planet: planet, Kind: "fuel", YieldPct: 1 - yield, QtyInitial: qty}
		nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn, QtyRemaining: nr.QtyInitial}}

		planet.Deposits = append(planet.Deposits, nr)
	}

	yield, qty = 0.95, 90_000_000
	for i := len(planet.Deposits); i < 22; i++ {
		yield, qty = yield*0.9, qty*8/10

		nr := &NaturalResource{Planet: planet, Kind: "non-metallic", YieldPct: 1 - yield, QtyInitial: qty}
		nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn, QtyRemaining: nr.QtyInitial}}

		planet.Deposits = append(planet.Deposits, nr)
	}

	yield, qty = 0.95, 90_000_000
	for i := len(planet.Deposits); i < 40; i++ {
		yield, qty = yield*0.9, qty*9/10

		nr := &NaturalResource{Planet: planet, Kind: "metallic", YieldPct: 1 - yield, QtyInitial: qty}
		nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn, QtyRemaining: nr.QtyInitial}}

		planet.Deposits = append(planet.Deposits, nr)
	}

	return planet
}

func (s *Store) genTerrestrial(star *Star, orbit int) *Planet {
	efftn, endtn := &Turn{}, &Turn{Year: 9999, Quarter: 4}

	planet := &Planet{Star: star, Kind: "terrestrial", OrbitNo: orbit}
	planet.Details = []*PlanetDetail{{Planet: planet, EffTurn: efftn, EndTurn: endtn}}

	if orbit <= 5 {
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5:
			planet.HabitabilityNo = rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 6, 7, 8, 9, 10:
			planet.HabitabilityNo = rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 11, 12, 13, 14:
			planet.HabitabilityNo = rand.Intn(4) + rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 15, 16, 17:
			planet.HabitabilityNo = rand.Intn(5) + rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 18, 19:
			planet.HabitabilityNo = rand.Intn(6) + rand.Intn(5) + rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 20:
			planet.HabitabilityNo = rand.Intn(7) + rand.Intn(6) + rand.Intn(5) + rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		}
	}

	for r := 0; r <= rand.Intn(40); r++ {
		nr := &NaturalResource{Planet: planet}
		nr.Details = []*NaturalResourceDetail{{NaturalResource: nr, EffTurn: efftn, EndTurn: endtn}}
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 11, 12, 13, 14, 15, 16, 17:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 18, 19:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
		case 20:
			nr.Kind, nr.YieldPct, nr.QtyInitial = "gold", 0.01+float64(rand.Intn(5))/100, rand.Intn(30)*100_000
		}
		if nr.QtyInitial < 100_000 {
			nr.QtyInitial = 100_000
		}
		nr.Details[0].QtyRemaining = nr.QtyInitial
		planet.Deposits = append(planet.Deposits, nr)
	}

	return planet
}
