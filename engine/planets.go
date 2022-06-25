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

type Planet struct {
	Id                 int
	Star               *Star // star the planet orbits
	Orbit              int
	Kind               string
	HomePlanet         bool
	HabitabilityNumber int
	Resources          []*NaturalResource
}

func (e *Engine) genAsteroidBelt(star *Star, orbit int) *Planet {
	planet := &Planet{Star: star, Kind: "asteroid-belt", Orbit: orbit}

	for r := 0; r <= rand.Intn(40); r++ {
		nr := &NaturalResource{}
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
			nr.Kind, nr.YieldPct, nr.InitialQuantity = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 11, 12, 13, 14, 15, 16, 17:
			nr.Kind, nr.YieldPct, nr.InitialQuantity = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 18, 19:
			nr.Kind, nr.YieldPct, nr.InitialQuantity = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
		case 20:
			nr.Kind, nr.YieldPct, nr.InitialQuantity = "gold", 0.01+float64(rand.Intn(5))/100, rand.Intn(30)*100_000
		}
		if nr.InitialQuantity < 100_000 {
			nr.InitialQuantity = 100_000
		}
		nr.QuantityRemaining = nr.InitialQuantity
		planet.Resources = append(planet.Resources, nr)
	}

	return planet
}

func (e *Engine) genEmpty(star *Star, orbit int) *Planet {
	planet := &Planet{Star: star, Kind: "empty", Orbit: orbit}
	return planet
}

func (e *Engine) genGasGiant(star *Star, orbit int) *Planet {
	planet := &Planet{Star: star, Kind: "gas-giant", Orbit: orbit}
	if 3 <= orbit && orbit <= 5 {
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5:
			planet.HabitabilityNumber = rand.Intn(1)
		case 6, 7, 8, 9, 10:
			planet.HabitabilityNumber = rand.Intn(1) + rand.Intn(1)
		case 11, 12, 13, 14:
			planet.HabitabilityNumber = rand.Intn(2) + rand.Intn(1) + rand.Intn(1)
		case 15, 16, 17:
			planet.HabitabilityNumber = rand.Intn(2) + rand.Intn(2) + rand.Intn(1) + rand.Intn(1)
		case 18, 19:
			planet.HabitabilityNumber = rand.Intn(3) + rand.Intn(2) + rand.Intn(2) + rand.Intn(1) + rand.Intn(1)
		case 20:
			planet.HabitabilityNumber = rand.Intn(3) + rand.Intn(3) + rand.Intn(2) + rand.Intn(2) + rand.Intn(1) + rand.Intn(1)
		}
	}

	nr := &NaturalResource{}
	switch rand.Intn(21) {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15:
		nr.Kind, nr.YieldPct, nr.InitialQuantity = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
	case 16, 17, 18, 19:
		nr.Kind, nr.YieldPct, nr.InitialQuantity = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
	case 20:
		nr.Kind, nr.YieldPct, nr.InitialQuantity = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
	}
	if nr.InitialQuantity < 100_000 {
		nr.InitialQuantity = 100_000
	}
	nr.QuantityRemaining = nr.InitialQuantity
	planet.Resources = append(planet.Resources, nr)

	return planet
}

func (e *Engine) genHomeTerrestrial(star *Star, orbit int) *Planet {
	planet := &Planet{Star: star, Kind: "terrestrial", Orbit: orbit, HomePlanet: true}

	planet.HabitabilityNumber = 25

	planet.Resources = append(planet.Resources, &NaturalResource{Kind: "gold", YieldPct: 0.07, InitialQuantity: 300_000, QuantityRemaining: 300_000})
	planet.Resources = append(planet.Resources, &NaturalResource{Kind: "fuel", YieldPct: 0.25, InitialQuantity: 99_999_999, QuantityRemaining: 99_999_999})
	planet.Resources = append(planet.Resources, &NaturalResource{Kind: "non-metallic", YieldPct: 0.25, InitialQuantity: 99_999_999, QuantityRemaining: 99_999_999})
	planet.Resources = append(planet.Resources, &NaturalResource{Kind: "metallic", YieldPct: 0.25, InitialQuantity: 99_999_999, QuantityRemaining: 99_999_999})

	yield, qty := 0.45, 90_000_000
	for i := len(planet.Resources); i < 8; i++ {
		yield, qty = yield*0.9, qty*8/10
		planet.Resources = append(planet.Resources, &NaturalResource{Kind: "fuel", YieldPct: 1 - yield, InitialQuantity: qty, QuantityRemaining: qty})
	}

	yield, qty = 0.95, 90_000_000
	for i := len(planet.Resources); i < 22; i++ {
		yield, qty = yield*0.9, qty*8/10
		planet.Resources = append(planet.Resources, &NaturalResource{Kind: "non-metallic", YieldPct: 1 - yield, InitialQuantity: qty, QuantityRemaining: qty})
	}

	yield, qty = 0.95, 90_000_000
	for i := len(planet.Resources); i < 40; i++ {
		yield, qty = yield*0.9, qty*9/10
		planet.Resources = append(planet.Resources, &NaturalResource{Kind: "metallic", YieldPct: 1 - yield, InitialQuantity: qty, QuantityRemaining: qty})
	}

	return planet
}

func (e *Engine) genTerrestrial(star *Star, orbit int) *Planet {
	planet := &Planet{Star: star, Kind: "terrestrial", Orbit: orbit}

	if orbit <= 5 {
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5:
			planet.HabitabilityNumber = rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 6, 7, 8, 9, 10:
			planet.HabitabilityNumber = rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 11, 12, 13, 14:
			planet.HabitabilityNumber = rand.Intn(4) + rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 15, 16, 17:
			planet.HabitabilityNumber = rand.Intn(5) + rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 18, 19:
			planet.HabitabilityNumber = rand.Intn(6) + rand.Intn(5) + rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		case 20:
			planet.HabitabilityNumber = rand.Intn(7) + rand.Intn(6) + rand.Intn(5) + rand.Intn(4) + rand.Intn(3) + rand.Intn(2) + rand.Intn(1)
		}
	}

	for r := 0; r <= rand.Intn(40); r++ {
		nr := &NaturalResource{}
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
			nr.Kind, nr.YieldPct, nr.InitialQuantity = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 11, 12, 13, 14, 15, 16, 17:
			nr.Kind, nr.YieldPct, nr.InitialQuantity = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 18, 19:
			nr.Kind, nr.YieldPct, nr.InitialQuantity = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
		case 20:
			nr.Kind, nr.YieldPct, nr.InitialQuantity = "gold", 0.01+float64(rand.Intn(5))/100, rand.Intn(30)*100_000
		}
		if nr.InitialQuantity < 100_000 {
			nr.InitialQuantity = 100_000
		}
		nr.QuantityRemaining = nr.InitialQuantity
		planet.Resources = append(planet.Resources, nr)
	}

	return planet
}
