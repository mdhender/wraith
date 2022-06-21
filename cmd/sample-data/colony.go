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

type Colony struct {
	Id       int    `json:"colony-id"`
	Kind     string `json:"kind"`
	Location struct {
		X     int `json:"x"`
		Y     int `json:"y"`
		Z     int `json:"z"`
		Star  int `json:"star,omitempty"`
		Orbit int `json:"orbit"`
	} `json:"location"`
	TechLevel  int `json:"tech-level"`
	Population struct {
		Professional      Population `json:"professional"`
		Soldier           Population `json:"soldier"`
		Unskilled         Population `json:"unskilled"`
		Unemployed        Population `json:"unemployed"`
		ConstructionCrews int        `json:"construction-crews,omitempty"`
		SpyTeams          int        `json:"spy-teams,omitempty"`
		Births            int        `json:"births,omitempty"`
		Deaths            int        `json:"deaths,omitempty"`
	} `json:"population"`
	Inventory     []*Inventory `json:"inventory,omitempty"`
	FactoryGroups []*Group     `json:"factory-groups,omitempty"`
	FarmGroups    []*Group     `json:"farm-groups,omitempty"`
	MiningGroups  []*Group     `json:"mining-groups,omitempty"`
}

type Inventory struct {
	Name           string `json:"name"`
	Code           string `json:"code,omitempty"`
	TechLevel      int    `json:"tech-level,omitempty"`
	OperationalQty int    `json:"operational-qty,omitempty"`
	StowedQty      int    `json:"stowed-qty,omitempty"`
	MassUnits      int    `json:"mass-units,omitempty"`
	EnclosedUnits  int    `json:"enclosed-units,omitempty"`
}

type Population struct {
	Code   string  `json:"code,omitempty"`
	Qty    int     `json:"qty,omitempty"`
	Pay    float64 `json:"pay,omitempty"`
	Ration float64 `json:"ration,omitempty"`
}

type Group struct {
	Id    int           `json:"group-id,omitempty"`
	Name  string        `json:"name,omitempty"`
	Units []*GroupUnits `json:"units,omitempty"`
}

type GroupUnits struct {
	TechLevel int   `json:"tech-level,omitempty"`
	Qty       int   `json:"qty,omitempty"`
	Stages    []int `json:"stages,omitempty"`
}

var numColonies int

func GenHomeOpenColony(id int) *Colony {
	numColonies++
	c := &Colony{Id: numColonies, Kind: "open", TechLevel: 1}

	c.Population.Professional = Population{Code: "PRO", Qty: 2_000_000, Pay: 1.0, Ration: 1.0}
	c.Population.Soldier = Population{Code: "SLD", Qty: 2_500_000, Pay: 1.0, Ration: 1.0}
	c.Population.Unskilled = Population{Code: "USK", Qty: 6_000_000, Pay: 1.0, Ration: 1.0}
	c.Population.Unemployed = Population{Code: "UEM", Qty: 5_900_000, Pay: 1.0, Ration: 1.0}
	c.Population.ConstructionCrews = 2_000
	c.Population.SpyTeams = 3

	// add operational inventory
	c.Inventory = append(c.Inventory, &Inventory{Code: "ANM", Name: "anti-missile", TechLevel: 1, OperationalQty: 25_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "ASC", Name: "assault-craft", TechLevel: 1, OperationalQty: 6_750})
	c.Inventory = append(c.Inventory, &Inventory{Code: "ASW", Name: "assault-weapon", TechLevel: 1, OperationalQty: 10_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FCT", Name: "factory", TechLevel: 1, OperationalQty: 275_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FRM", Name: "farm", TechLevel: 1, OperationalQty: 170_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MTLS", Name: "metallics", OperationalQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MLSP", Name: "military-supplies", OperationalQty: 2_000_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MIN", Name: "mine", TechLevel: 1, OperationalQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MSS", Name: "missile", TechLevel: 1, OperationalQty: 240_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MSL", Name: "missile-launcher", TechLevel: 1, OperationalQty: 8_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "NMTS", Name: "non-metallics", OperationalQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "STUN", Name: "structural", OperationalQty: 87_500_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "TPT", Name: "transport", TechLevel: 1, OperationalQty: 5_000})

	// add disassembled inventory
	c.Inventory = append(c.Inventory, &Inventory{Code: "CNGD", Name: "consumer-goods", StowedQty: 2_000_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FCT", Name: "factory", TechLevel: 1, StowedQty: 3_750_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FOOD", Name: "food", StowedQty: 7_500_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FUEL", Name: "fuel", StowedQty: 5_000_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MIN", Name: "mine", TechLevel: 1, StowedQty: 30_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "STUN", Name: "structural", StowedQty: 150_000})

	return c
}

func GenHomeOrbitalColony(id int) *Colony {
	numColonies++
	c := &Colony{Id: numColonies, Kind: "orbital", TechLevel: 1}

	c.Population.Professional = Population{Code: "PRO", Qty: 10_000, Pay: 1.0, Ration: 1.0}
	c.Population.Soldier = Population{Code: "SLD", Qty: 20, Pay: 1.0, Ration: 1.0}
	c.Population.Unskilled = Population{Code: "USK", Qty: 30_000, Pay: 1.0, Ration: 1.0}
	c.Population.Unemployed = Population{Code: "UEM", Qty: 500, Pay: 1.0, Ration: 1.0}
	c.Population.ConstructionCrews = 100
	c.Population.SpyTeams = 0

	// add operational inventory
	c.Inventory = append(c.Inventory, &Inventory{Code: "LSP", Name: "life-support", TechLevel: 1, OperationalQty: 2_500})
	c.Inventory = append(c.Inventory, &Inventory{Code: "LTSU", Name: "light-structural", OperationalQty: 45_000_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "SNR", Name: "sensors", TechLevel: 1, OperationalQty: 5_000})

	// add disassembled inventory
	c.Inventory = append(c.Inventory, &Inventory{Code: "CNGD", Name: "consumer-goods", StowedQty: 2_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FOOD", Name: "food", StowedQty: 500_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FUEL", Name: "fuel", StowedQty: 500_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "HDR", Name: "hyper-drive", TechLevel: 1, StowedQty: 500})
	c.Inventory = append(c.Inventory, &Inventory{Code: "LTSU", Name: "light-structural", StowedQty: 5_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MTLS", Name: "metallics", StowedQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "NMTS", Name: "non-metallics", StowedQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "SDR", Name: "star-drive", TechLevel: 1, StowedQty: 250})
	return c
}
