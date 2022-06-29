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

type unitValues struct {
	code, name string
	tl, oq, sq int
}

func (s *Store) genHomeOpenColony(no int, planet *Planet, player *Player) *ColonyOrShip {
	effTurn, endTurn := &Turn{}, &Turn{Year: 9999, Quarter: 4}

	c := &ColonyOrShip{MSN: no, Kind: "open", HomeColony: true}
	c.Details = []*CSDetail{{
		CS:           c,
		EffTurn:      effTurn,
		EndTurn:      endTurn,
		TechLevel:    1,
		Name:         "Not Named",
		ControlledBy: player,
	}}
	c.Locations = []*CSLocation{{
		CS:       c,
		EffTurn:  effTurn,
		EndTurn:  endTurn,
		Location: planet,
	}}

	// create hull
	for _, unit := range []unitValues{
		{code: "STUN", name: "structural", tl: 0, oq: 87_500_000},
		{code: "ANM", name: "anti-missile", tl: 1, oq: 25_000},
		{code: "MSL", name: "missile-launcher", tl: 1, oq: 8_000},
		{code: "MSS", name: "missile", tl: 1, oq: 240_000},
		{code: "SNR", name: "sensor", tl: 1, oq: 50},
	} {
		c.Hull = append(c.Hull, &CSHull{
			CS:      c,
			EffTurn: effTurn,
			EndTurn: endTurn,
			Unit: &Unit{
				Code:      unit.code,
				TechLevel: unit.tl,
				Name:      unit.name,
			},
			QtyOperational: unit.oq,
		})
	}

	// add cargo
	for _, unit := range []unitValues{
		{code: "ASC", name: "assault-craft", tl: 1, oq: 6_750, sq: 0},
		{code: "ASW", name: "assault-weapon", tl: 1, oq: 10_000, sq: 0},
		{code: "CNGD", name: "consumer-goods", tl: 0, oq: 0, sq: 2_000_000},
		{code: "FCT", name: "factory", tl: 1, oq: 275_000, sq: 3_750_000},
		{code: "FOOD", name: "food", tl: 0, oq: 0, sq: 7_500_000},
		{code: "FRM", name: "farm", tl: 1, oq: 170_000, sq: 0},
		{code: "FUEL", name: "fuel", tl: 0, oq: 0, sq: 5_000_000},
		{code: "MIN", name: "mine", tl: 1, oq: 100_000, sq: 30_000},
		{code: "MTLS", name: "metallics", tl: 0, oq: 100_000, sq: 0},
		{code: "MLSP", name: "military-supplies", tl: 0, oq: 2_000_000},
		{code: "NMTS", name: "non-metallics", tl: 0, oq: 100_000, sq: 0},
		{code: "STUN", name: "structural", tl: 0, oq: 0, sq: 150_000},
		{code: "TPT", name: "transport", tl: 1, oq: 5_000, sq: 0},
	} {
		c.Inventory = append(c.Inventory, &CSInventory{
			CS:      c,
			EffTurn: effTurn,
			EndTurn: endTurn,
			Unit: &Unit{
				Code:      unit.code,
				TechLevel: unit.tl,
				Name:      unit.name,
			},
			QtyOperational: unit.oq,
			QtyStowed:      unit.sq,
		})

	}

	c.Pay = []*CSPay{{
		CS:              c,
		EffTurn:         effTurn,
		EndTurn:         endTurn,
		ProfessionalPct: 1.0,
		SoldierPct:      1.0,
		UnskilledPct:    1.0,
		UnemployedPct:   1.0,
	}}

	c.Population = []*CSPopulation{{
		CS:                  c,
		EffTurn:             effTurn,
		EndTurn:             endTurn,
		QtyProfessional:     2_000_000,
		QtySoldier:          2_500_000,
		QtyUnskilled:        6_000_000,
		QtyUnemployed:       5_900_000,
		QtyConstructionCrew: 2_000,
		QtySpyTeam:          3,
		RebelPct:            0.0125,
	}}

	c.Rations = []*CSRations{{
		CS:              c,
		EffTurn:         effTurn,
		EndTurn:         endTurn,
		ProfessionalPct: 1.0,
		SoldierPct:      1.0,
		UnskilledPct:    1.0,
		UnemployedPct:   1.0,
	}}

	factoryGroup := &FactoryGroup{
		CS:      c,
		No:      1,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Unit:    &Unit{Code: "CNGD"},
	}
	factoryGroup.Units = []*FactoryGroupUnits{{
		Group:   factoryGroup,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Unit: &Unit{
			Code:      "FCT",
			TechLevel: 1,
		},
		QtyOperational: 275_000,
	}}
	factoryGroup.Stages = []*FactoryGroupStages{{
		Group:     factoryGroup,
		Turn:      effTurn,
		QtyStage1: 2_291_666,
		QtyStage2: 2_291_666,
		QtyStage3: 2_291_666,
		QtyStage4: 0,
	}}
	c.Factories = append(c.Factories, factoryGroup)

	miningGroup := &MiningGroup{
		CS:      c,
		No:      1,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Deposit: planet.Deposits[0],
	}
	miningGroup.Units = []*MiningGroupUnits{{
		Group:   miningGroup,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Unit: &Unit{
			Code:      "MIN",
			TechLevel: 1,
		},
		QtyOperational: 1_000,
	}}
	miningGroup.Stages = []*MiningGroupStages{{
		Group:     miningGroup,
		Turn:      effTurn,
		QtyStage1: 1_000,
		QtyStage2: 1_000,
		QtyStage3: 1_000,
		QtyStage4: 0,
	}}
	c.Mines = append(c.Mines, miningGroup)

	miningGroup = &MiningGroup{
		CS:      c,
		No:      2,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Deposit: planet.Deposits[1],
	}
	miningGroup.Units = []*MiningGroupUnits{{
		Group:   miningGroup,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Unit: &Unit{
			Code:      "MIN",
			TechLevel: 1,
		},
		QtyOperational: 50_000,
	}}
	miningGroup.Stages = []*MiningGroupStages{{
		Group:     miningGroup,
		Turn:      effTurn,
		QtyStage1: 1_250_000,
		QtyStage2: 1_250_000,
		QtyStage3: 1_250_000,
		QtyStage4: 0,
	}}
	c.Mines = append(c.Mines, miningGroup)

	miningGroup = &MiningGroup{
		CS:      c,
		No:      3,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Deposit: planet.Deposits[2],
	}
	miningGroup.Units = []*MiningGroupUnits{{
		Group:   miningGroup,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Unit: &Unit{
			Code:      "MIN",
			TechLevel: 1,
		},
		QtyOperational: 100_000,
	}}
	miningGroup.Stages = []*MiningGroupStages{{
		Group:     miningGroup,
		Turn:      effTurn,
		QtyStage1: 2_500_000,
		QtyStage2: 2_500_000,
		QtyStage3: 2_500_000,
		QtyStage4: 0,
	}}
	c.Mines = append(c.Mines, miningGroup)

	miningGroup = &MiningGroup{
		CS:      c,
		No:      4,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Deposit: planet.Deposits[3],
	}
	miningGroup.Units = []*MiningGroupUnits{{
		Group:   miningGroup,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Unit: &Unit{
			Code:      "MIN",
			TechLevel: 1,
		},
		QtyOperational: 100_000,
	}}
	miningGroup.Stages = []*MiningGroupStages{{
		Group:     miningGroup,
		Turn:      effTurn,
		QtyStage1: 2_500_000,
		QtyStage2: 2_500_000,
		QtyStage3: 2_500_000,
		QtyStage4: 0,
	}}
	c.Mines = append(c.Mines, miningGroup)

	return c
}

func (s *Store) genHomeOrbitalColony(no int, planet *Planet, player *Player) *ColonyOrShip {
	effTurn, endTurn := &Turn{}, &Turn{Year: 9999, Quarter: 4}

	c := &ColonyOrShip{MSN: no, Kind: "orbital", HomeColony: true}
	c.Details = []*CSDetail{{
		CS:           c,
		EffTurn:      effTurn,
		EndTurn:      endTurn,
		TechLevel:    1,
		Name:         "Not Named",
		ControlledBy: player,
	}}
	c.Locations = []*CSLocation{{
		CS:       c,
		EffTurn:  effTurn,
		EndTurn:  endTurn,
		Location: planet,
	}}

	// create hull
	for _, unit := range []unitValues{
		{code: "LSP", name: "life-support", tl: 1, oq: 2_000},
		{code: "SNR", name: "sensor", tl: 1, oq: 5_000},
		{code: "STUN", name: "structural", tl: 0, oq: 45_000_000},
	} {
		c.Hull = append(c.Hull, &CSHull{
			CS:      c,
			EffTurn: effTurn,
			EndTurn: endTurn,
			Unit: &Unit{
				Code:      unit.code,
				TechLevel: unit.tl,
				Name:      unit.name,
			},
			QtyOperational: unit.oq,
		})
	}

	// add cargo
	for _, unit := range []unitValues{
		{code: "CNGD", name: "consumer-goods", oq: 0, sq: 2_000},
		{code: "FOOD", name: "food", tl: 0, oq: 0, sq: 500_000},
		{code: "FUEL", name: "fuel", tl: 0, oq: 0, sq: 500_000},
		{code: "HDR", name: "hyper-drive", tl: 1, oq: 0, sq: 500},
		{code: "LTSU", name: "light-structural", tl: 0, oq: 45_000_000, sq: 5_000},
		{code: "MTLS", name: "metallics", tl: 0, oq: 0, sq: 100_000},
		{code: "NMTS", name: "non-metallics", tl: 0, oq: 0, sq: 100_000},
		{code: "SDR", name: "star-drive", tl: 1, oq: 0, sq: 250},
	} {
		c.Inventory = append(c.Inventory, &CSInventory{
			CS:      c,
			EffTurn: effTurn,
			EndTurn: endTurn,
			Unit: &Unit{
				Code:      unit.code,
				TechLevel: unit.tl,
				Name:      unit.name,
			},
			QtyOperational: unit.oq,
			QtyStowed:      unit.sq,
		})

	}

	c.Pay = []*CSPay{{
		CS:              c,
		EffTurn:         effTurn,
		EndTurn:         endTurn,
		ProfessionalPct: 1.0,
		SoldierPct:      1.0,
		UnskilledPct:    1.0,
		UnemployedPct:   1.0,
	}}

	c.Population = []*CSPopulation{{
		CS:                  c,
		EffTurn:             effTurn,
		EndTurn:             endTurn,
		QtyProfessional:     10_000,
		QtySoldier:          20,
		QtyUnskilled:        30_000,
		QtyUnemployed:       500,
		QtyConstructionCrew: 100,
	}}

	c.Rations = []*CSRations{{
		CS:              c,
		EffTurn:         effTurn,
		EndTurn:         endTurn,
		ProfessionalPct: 1.0,
		SoldierPct:      1.0,
		UnskilledPct:    1.0,
		UnemployedPct:   1.0,
	}}

	factoryGroup := &FactoryGroup{
		CS:      c,
		No:      1,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Unit:    &Unit{Code: "LTSU"},
	}
	factoryGroup.Units = []*FactoryGroupUnits{{
		Group:   factoryGroup,
		EffTurn: effTurn,
		EndTurn: endTurn,
		Unit: &Unit{
			Code:      "FCT",
			TechLevel: 1,
		},
		QtyOperational: 5_000,
	}}
	factoryGroup.Stages = []*FactoryGroupStages{{
		Group:     factoryGroup,
		Turn:      effTurn,
		QtyStage1: 500_000,
		QtyStage2: 500_000,
		QtyStage3: 500_000,
		QtyStage4: 0,
	}}
	c.Factories = append(c.Factories, factoryGroup)

	return c
}
