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

type XColony struct {
	Id         int // key for database
	No         int
	Name       string
	Kind       string
	Location   *Planet
	TechLevel  int
	Population struct {
		Professional      XPopulation
		Soldier           XPopulation
		Unskilled         XPopulation
		Unemployed        XPopulation
		ConstructionCrews int
		SpyTeams          int
		RebelPct          float64
		Births            int
		Deaths            int
	}
	Hull          []*Inventory // units used to build colony
	Inventory     []*Inventory // units stored in colony
	FactoryGroups []*XGroup
	FarmGroups    []*XGroup
	MiningGroups  []*MiningGroup
}

type XPopulation struct {
	Code   string
	Qty    int
	Pay    float64
	Ration float64
}

type MiningGroup struct {
	Id      int // key in database
	No      int // mining group number
	Deposit *NaturalResource
	Units   []*XGroupUnits
}

type XGroup struct {
	Id    int
	Name  string
	Units []*XGroupUnits
}

type XGroupUnits struct {
	TechLevel int
	Qty       int
	Stages    []int
}

func (e *Engine) genHomeOpenColony(planet *Planet) *XColony {
	c := &XColony{Kind: "open", TechLevel: 1, Name: "Not Named"}

	c.Population.Professional = XPopulation{Code: "PRO", Qty: 2_000_000, Pay: 1.0, Ration: 1.0}
	c.Population.Soldier = XPopulation{Code: "SLD", Qty: 2_500_000, Pay: 1.0, Ration: 1.0}
	c.Population.Unskilled = XPopulation{Code: "USK", Qty: 6_000_000, Pay: 1.0, Ration: 1.0}
	c.Population.Unemployed = XPopulation{Code: "UEM", Qty: 5_900_000, Pay: 1.0, Ration: 1.0}
	c.Population.ConstructionCrews = 2_000
	c.Population.SpyTeams = 3
	c.Population.RebelPct = 0.0125

	// create hull
	c.Hull = append(c.Hull, &Inventory{Code: "STUN", Name: "structural", OperationalQty: 87_500_000})
	c.Hull = append(c.Hull, &Inventory{Code: "ANM", Name: "anti-missile", TechLevel: 1, OperationalQty: 25_000})
	c.Hull = append(c.Hull, &Inventory{Code: "MSL", TechLevel: 1, Name: "missile-launcher", OperationalQty: 8_000})
	c.Hull = append(c.Hull, &Inventory{Code: "MSS", TechLevel: 1, Name: "missile", OperationalQty: 240_000})
	c.Hull = append(c.Hull, &Inventory{Code: "SNR", TechLevel: 1, Name: "sensor", OperationalQty: 50})

	// add cargo
	c.Inventory = append(c.Inventory, &Inventory{Code: "ASC", Name: "assault-craft", TechLevel: 1, OperationalQty: 6_750})
	c.Inventory = append(c.Inventory, &Inventory{Code: "ASW", Name: "assault-weapon", TechLevel: 1, OperationalQty: 10_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "CNGD", Name: "consumer-goods", StowedQty: 2_000_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FCT", Name: "factory", TechLevel: 1, OperationalQty: 275_000, StowedQty: 3_750_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FOOD", Name: "food", StowedQty: 7_500_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FRM", Name: "farm", TechLevel: 1, OperationalQty: 170_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FUEL", Name: "fuel", StowedQty: 5_000_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MTLS", Name: "metallics", OperationalQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MLSP", Name: "military-supplies", OperationalQty: 2_000_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MIN", Name: "mine", TechLevel: 1, OperationalQty: 100_000, StowedQty: 30_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "NMTS", Name: "non-metallics", OperationalQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "STUN", Name: "structural", StowedQty: 150_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "TPT", Name: "transport", TechLevel: 1, OperationalQty: 5_000})

	miningGroup := &MiningGroup{No: 1, Deposit: planet.Resources[0], Units: []*XGroupUnits{{TechLevel: 1, Qty: 1_000, Stages: []int{1_000, 1_000, 1_000}}}}
	c.MiningGroups = append(c.MiningGroups, miningGroup)
	miningGroup = &MiningGroup{No: 2, Deposit: planet.Resources[1], Units: []*XGroupUnits{{TechLevel: 1, Qty: 50_000, Stages: []int{1_250_000, 1_250_000, 1_250_000}}}}
	c.MiningGroups = append(c.MiningGroups, miningGroup)
	miningGroup = &MiningGroup{No: 3, Deposit: planet.Resources[2], Units: []*XGroupUnits{{TechLevel: 1, Qty: 100_000, Stages: []int{2_500_000, 2_500_000, 2_500_000}}}}
	c.MiningGroups = append(c.MiningGroups, miningGroup)
	miningGroup = &MiningGroup{No: 4, Deposit: planet.Resources[4], Units: []*XGroupUnits{{TechLevel: 1, Qty: 100_000, Stages: []int{2_500_000, 2_500_000, 2_500_000}}}}
	c.MiningGroups = append(c.MiningGroups, miningGroup)

	return c
}

func (e *Engine) genHomeOrbitalColony(planet *Planet) *XColony {
	c := &XColony{Kind: "orbital", TechLevel: 1}

	c.Population.Professional = XPopulation{Code: "PRO", Qty: 10_000, Pay: 1.0, Ration: 1.0}
	c.Population.Soldier = XPopulation{Code: "SLD", Qty: 20, Pay: 1.0, Ration: 1.0}
	c.Population.Unskilled = XPopulation{Code: "USK", Qty: 30_000, Pay: 1.0, Ration: 1.0}
	c.Population.Unemployed = XPopulation{Code: "UEM", Qty: 500, Pay: 1.0, Ration: 1.0}
	c.Population.ConstructionCrews = 100
	c.Population.SpyTeams = 0

	// create hull
	c.Hull = append(c.Hull, &Inventory{Code: "STUN", Name: "structural", OperationalQty: 45_000_000})
	c.Hull = append(c.Hull, &Inventory{Code: "LSP", TechLevel: 1, Name: "life-support", OperationalQty: 2_500})
	c.Hull = append(c.Hull, &Inventory{Code: "SNR", TechLevel: 1, Name: "sensor", OperationalQty: 5_000})

	// add cargo
	c.Inventory = append(c.Inventory, &Inventory{Code: "CNGD", Name: "consumer-goods", StowedQty: 2_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FOOD", Name: "food", StowedQty: 500_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "FUEL", Name: "fuel", StowedQty: 500_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "HDR", Name: "hyper-drive", TechLevel: 1, StowedQty: 500})
	c.Inventory = append(c.Inventory, &Inventory{Code: "LTSU", Name: "light-structural", OperationalQty: 45_000_000, StowedQty: 5_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "MTLS", Name: "metallics", StowedQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "NMTS", Name: "non-metallics", StowedQty: 100_000})
	c.Inventory = append(c.Inventory, &Inventory{Code: "SDR", Name: "star-drive", TechLevel: 1, StowedQty: 250})
	return c
}
