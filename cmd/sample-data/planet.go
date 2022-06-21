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

import "math/rand"

type Planet struct {
	Id                 int                `json:"planet-id,omitempty"`
	Orbit              int                `json:"orbit,omitempty"`
	Kind               string             `json:"kind"`
	HomePlanet         bool               `json:"home-planet,omitempty"`
	HabitabilityNumber int                `json:"habitability-number,omitempty"`
	Deposits           []*NaturalResource `json:"deposits,omitempty"`
	// From Far Horizons
	//Diameter   int     `json:"diameter,omitempty"`
	//Gravity    float64 `json:"gravity,omitempty"`
	//TemperatureClass   int        `json:"temperature-class"`
	//PressureClass      int        `json:"pressure-class"`
	//MiningDifficulty   float64 `json:"mining-difficulty,omitempty"`
	//Gases              []*GAS     `json:"gases"`
}

//type GAS struct {
//	Code       string `json:"code"`
//	Percentage int    `json:"pct"`
//}

type NaturalResource struct {
	Id                int     `json:"natural-resource-id,omitempty"`
	Kind              string  `json:"kind,omitempty"`
	Yield             float64 `json:"yield,omitempty"`
	InitialQuantity   int     `json:"initial-quantity,omitempty"`
	QuantityRemaining int     `json:"quantity-remaining,omitempty"`
}

var numAsteroidBelts int
var numGasGiants int
var numNaturalResources int
var numTerrestrials int

func GenAsteroidBelt(id, orbit int) *Planet {
	numAsteroidBelts++
	planet := &Planet{Id: id, Kind: "asteroid-belt", Orbit: orbit}

	for r := 0; r <= rand.Intn(40); r++ {
		numNaturalResources++
		nr := &NaturalResource{Id: numNaturalResources}
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
			nr.Kind, nr.Yield, nr.InitialQuantity = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 11, 12, 13, 14, 15, 16, 17:
			nr.Kind, nr.Yield, nr.InitialQuantity = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 18, 19:
			nr.Kind, nr.Yield, nr.InitialQuantity = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
		case 20:
			nr.Kind, nr.Yield, nr.InitialQuantity = "gold", 0.01+float64(rand.Intn(5))/100, rand.Intn(30)*100_000
		}
		if nr.InitialQuantity < 100_000 {
			nr.InitialQuantity = 100_000
		}
		nr.QuantityRemaining = nr.InitialQuantity
		planet.Deposits = append(planet.Deposits, nr)
	}

	return planet
}

func GenEmpty(id, orbit int) *Planet {
	planet := &Planet{Id: id, Kind: "empty", Orbit: orbit}
	return planet
}

func GenGasGiant(id, orbit int) *Planet {
	numGasGiants++
	planet := &Planet{Id: id, Kind: "gas-giant", Orbit: orbit}
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

	numNaturalResources++
	nr := &NaturalResource{Id: numNaturalResources}
	switch rand.Intn(21) {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15:
		nr.Kind, nr.Yield, nr.InitialQuantity = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
	case 16, 17, 18, 19:
		nr.Kind, nr.Yield, nr.InitialQuantity = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
	case 20:
		nr.Kind, nr.Yield, nr.InitialQuantity = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
	}
	nr.QuantityRemaining = nr.InitialQuantity
	planet.Deposits = append(planet.Deposits, nr)

	return planet
}

func GenHomeTerrestrial(id, orbit int) *Planet {
	numTerrestrials++
	planet := &Planet{Id: id, Kind: "terrestrial", Orbit: orbit, HomePlanet: true}

	planet.HabitabilityNumber = 25

	numNaturalResources++
	planet.Deposits = append(planet.Deposits, &NaturalResource{Id: numNaturalResources, Kind: "gold", Yield: 0.07, InitialQuantity: 300_000, QuantityRemaining: 300_000})
	numNaturalResources++
	planet.Deposits = append(planet.Deposits, &NaturalResource{Id: numNaturalResources, Kind: "fuel", Yield: 0.25, InitialQuantity: 99_999_999, QuantityRemaining: 99_999_999})

	yield, qty := 0.45, 90_000_000
	for i := len(planet.Deposits); i < 8; i++ {
		numNaturalResources++
		yield, qty = yield*0.9, qty*8/10
		planet.Deposits = append(planet.Deposits, &NaturalResource{Id: numNaturalResources, Kind: "fuel", Yield: 1 - yield, InitialQuantity: qty, QuantityRemaining: qty})
	}

	numNaturalResources++
	planet.Deposits = append(planet.Deposits, &NaturalResource{Id: numNaturalResources, Kind: "non-metallic", Yield: 0.25, InitialQuantity: 99_999_999, QuantityRemaining: 99_999_999})

	yield, qty = 0.95, 90_000_000
	for i := len(planet.Deposits); i < 22; i++ {
		numNaturalResources++
		yield, qty = yield*0.9, qty*8/10
		planet.Deposits = append(planet.Deposits, &NaturalResource{Id: numNaturalResources, Kind: "non-metallic", Yield: 1 - yield, InitialQuantity: qty, QuantityRemaining: qty})
	}

	numNaturalResources++
	planet.Deposits = append(planet.Deposits, &NaturalResource{Id: numNaturalResources, Kind: "metallic", Yield: 0.25, InitialQuantity: 99_999_999, QuantityRemaining: 99_999_999})

	yield, qty = 0.95, 90_000_000
	for i := len(planet.Deposits); i < 40; i++ {
		numNaturalResources++
		yield, qty = yield*0.9, qty*9/10
		planet.Deposits = append(planet.Deposits, &NaturalResource{Id: numNaturalResources, Kind: "metallic", Yield: 1 - yield, InitialQuantity: qty, QuantityRemaining: qty})
	}

	return planet
}

func GenTerrestrial(id, orbit int) *Planet {
	numTerrestrials++
	planet := &Planet{Id: id, Kind: "terrestrial", Orbit: orbit}

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
		numNaturalResources++
		nr := &NaturalResource{Id: numNaturalResources}
		switch rand.Intn(21) {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10:
			nr.Kind, nr.Yield, nr.InitialQuantity = "metallic", 0.75+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 11, 12, 13, 14, 15, 16, 17:
			nr.Kind, nr.Yield, nr.InitialQuantity = "non-metallic", 0.50+float64(rand.Intn(25))/100, rand.Intn(100)*1_000_000
		case 18, 19:
			nr.Kind, nr.Yield, nr.InitialQuantity = "fuel", 0.10+float64(rand.Intn(35))/100, rand.Intn(100)*1_000_000
		case 20:
			nr.Kind, nr.Yield, nr.InitialQuantity = "gold", 0.01+float64(rand.Intn(5))/100, rand.Intn(30)*100_000
		}
		nr.QuantityRemaining = nr.InitialQuantity
		planet.Deposits = append(planet.Deposits, nr)
	}

	return planet
}
