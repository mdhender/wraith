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

package adapters

import (
	"github.com/mdhender/wraith/storage/jdb"
	"github.com/mdhender/wraith/wraith"
	"sort"
	"time"
)

// WraithEngineToJdbGame converts an Engine to a Game.
func WraithEngineToJdbGame(e *wraith.Engine) *jdb.Game {
	jg := &jdb.Game{
		Id:        e.Game.Id,
		ShortName: e.Game.Code,
		Name:      e.Game.Name,
	}
	jg.Turn.Year = e.Game.Turn.Year
	jg.Turn.Quarter = e.Game.Turn.Quarter
	jg.Turn.StartDt = e.Game.Turn.StartDt.Format(time.RFC3339)
	jg.Turn.EndDt = e.Game.Turn.EndDt.Format(time.RFC3339)

	for _, colony := range e.Colonies {
		switch colony.Kind {
		case "enclosed":
			jg.EnclosedColonies = append(jg.EnclosedColonies, wraithColonyToJdbEnclosedColony(colony, jg))
		case "open", "surface":
			jg.SurfaceColonies = append(jg.SurfaceColonies, wraithColonyToJdbSurfaceColony(colony, jg))
		case "orbital":
			jg.OrbitalColonies = append(jg.OrbitalColonies, wraithColonyToJdbOrbitalColony(colony, jg))
		}
	}
	sort.Sort(jg.EnclosedColonies)
	sort.Sort(jg.FactoryGroups)
	sort.Sort(jg.FarmGroups)
	sort.Sort(jg.MineGroups)
	sort.Sort(jg.OrbitalColonies)
	sort.Sort(jg.SurfaceColonies)

	for _, deposit := range e.Deposits {
		jg.Deposits = append(jg.Deposits, wraithDepositToJdbDeposit(deposit))
	}
	sort.Sort(jg.Deposits)

	for _, nation := range e.Nations {
		n := &jdb.Nation{
			Id:                 nation.Id,
			No:                 nation.No,
			Name:               nation.Name,
			GovtName:           nation.GovtName,
			GovtKind:           nation.GovtKind,
			Speciality:         nation.Speciality,
			TechLevel:          nation.TechLevel,
			ResearchPointsPool: nation.ResearchPointsPool,
		}
		if nation.ControlledBy != nil {
			n.ControlledByPlayerId = nation.ControlledBy.Id
		}
		if nation.HomePlanet != nil {
			n.HomePlanetId = nation.HomePlanet.Id
		}
		n.Skills.Biology = nation.Skills.Biology
		n.Skills.Bureaucracy = nation.Skills.Bureaucracy
		n.Skills.Gravitics = nation.Skills.Gravitics
		n.Skills.LifeSupport = nation.Skills.LifeSupport
		n.Skills.Manufacturing = nation.Skills.Manufacturing
		n.Skills.Military = nation.Skills.Military
		n.Skills.Mining = nation.Skills.Mining
		n.Skills.Shields = nation.Skills.Shields

		jg.Nations = append(jg.Nations, n)
	}
	sort.Sort(jg.Nations)

	for _, planet := range e.Planets {
		p := &jdb.Planet{
			Id:             planet.Id,
			SystemId:       planet.System.Id,
			StarId:         planet.Star.Id,
			OrbitNo:        planet.OrbitNo,
			Kind:           planet.Kind,
			HabitabilityNo: planet.HabitabilityNo,
		}

		for _, colony := range planet.Colonies {
			switch colony.Kind {
			case "enclosed":
				p.EnclosedColonyIds = append(p.EnclosedColonyIds, colony.Id)
			case "open", "surface":
				p.SurfaceColonyIds = append(p.SurfaceColonyIds, colony.Id)
			case "orbital":
				p.OrbitalColonyIds = append(p.OrbitalColonyIds, colony.Id)
			}
		}
		sort.Ints(p.EnclosedColonyIds)
		sort.Ints(p.OrbitalColonyIds)
		sort.Ints(p.SurfaceColonyIds)

		for _, deposit := range planet.Deposits {
			p.DepositIds = append(p.DepositIds, deposit.Id)
		}
		sort.Ints(p.DepositIds)

		for _, ship := range planet.Ships {
			p.ShipIds = append(p.ShipIds, ship.Id)
		}
		sort.Ints(p.ShipIds)

		jg.Planets = append(jg.Planets, p)
	}
	sort.Sort(jg.Planets)

	for _, player := range e.Players {
		p := &jdb.Player{
			Id:       player.Id,
			UserId:   player.UserId,
			Name:     player.Name,
			MemberOf: player.MemberOf.Id,
		}
		if player.ReportsTo != nil {
			p.ReportsToPlayerId = player.ReportsTo.Id
		}
		jg.Players = append(jg.Players, p)
	}
	sort.Sort(jg.Players)

	for _, ship := range e.Ships {
		jg.Ships = append(jg.Ships, wraithShipToJdbShip(ship, jg))
	}
	sort.Sort(jg.Ships)

	for _, star := range e.Stars {
		s := &jdb.Star{
			Id:        star.Id,
			SystemId:  star.System.Id,
			Sequence:  star.Sequence,
			Kind:      star.Kind,
			PlanetIds: nil,
		}
		for _, planet := range star.Planets {
			s.PlanetIds = append(s.PlanetIds, planet.Id)
		}
		jg.Stars = append(jg.Stars, s)
	}
	sort.Sort(jg.Stars)

	for _, system := range e.Systems {
		s := &jdb.System{
			Id: system.Id,
			Coords: jdb.Coordinates{
				X: system.Coords.X,
				Y: system.Coords.Y,
				Z: system.Coords.Z,
			},
		}
		for _, star := range system.Stars {
			s.StarIds = append(s.StarIds, star.Id)
		}
		jg.Systems = append(jg.Systems, s)
	}
	sort.Sort(jg.Systems)

	for _, unit := range e.Units {
		u := &jdb.Unit{
			Id:                  unit.Id,
			Kind:                unit.Kind,
			Code:                unit.Code,
			TechLevel:           unit.TechLevel,
			Name:                unit.Name,
			Description:         unit.Description,
			MassPerUnit:         unit.MassPerUnit,
			VolumePerUnit:       unit.VolumePerUnit,
			Hudnut:              unit.Hudnut,
			StowedVolumePerUnit: unit.StowedVolumePerUnit,
			FuelPerUnitPerTurn:  unit.FuelPerUnitPerTurn,
			MetsPerUnit:         unit.MetsPerUnitPerTurn,
			NonMetsPerUnit:      unit.NonMetsPerUnitPerTurn,
		}
		jg.Units = append(jg.Units, u)
	}
	sort.Sort(jg.Units)

	//jg.Deposits = nil
	//jg.FactoryGroups = nil
	//jg.FarmGroups = nil
	//jg.MineGroups = nil
	//jg.Players = nil
	//jg.Units = nil

	return jg
}

func wraithDepositToJdbDeposit(deposit *wraith.Deposit) *jdb.Deposit {
	d := &jdb.Deposit{
		Id:           deposit.Id,
		PlanetId:     deposit.Planet.Id,
		No:           deposit.No,
		UnitId:       deposit.Product.Id,
		InitialQty:   deposit.InitialQty,
		RemainingQty: deposit.RemainingQty,
		YieldPct:     deposit.YieldPct,
	}
	if deposit.ControlledBy != nil {
		d.ControlledByColonyId = deposit.ControlledBy.Id
	}
	return d
}

func wraithColonyToJdbEnclosedColony(colony *wraith.CorS, game *jdb.Game) *jdb.EnclosedColony {
	c := &jdb.EnclosedColony{
		Id:              colony.Id,
		MSN:             colony.MSN,
		BuiltByNationId: colony.BuiltBy.Id,
		Name:            colony.Name,
		TechLevel:       colony.TechLevel,
		PlanetId:        colony.Planet.Id,
	}

	if colony.ControlledBy != nil {
		c.ControlledByPlayerId = colony.ControlledBy.Id
	}

	for _, group := range colony.FactoryGroups {
		c.FactoryGroupIds = append(c.FactoryGroupIds, group.Id)
		fg := &jdb.FactoryGroup{
			Id:        group.Id,
			CorSId:    group.CorS.Id,
			No:        group.No,
			Product:   group.Product.Id,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		for _, unit := range group.Units {
			fg.Units = append(fg.Units, &jdb.FactoryGroupUnits{
				UnitId:   unit.Unit.Id,
				TotalQty: unit.TotalQty,
			})
		}
		game.FactoryGroups = append(game.FactoryGroups, fg)
	}

	for _, group := range colony.FarmGroups {
		c.FarmGroupIds = append(c.FarmGroupIds, group.Id)
		fg := &jdb.FarmGroup{
			Id:        group.Id,
			CorSId:    group.CorS.Id,
			No:        group.No,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		for _, unit := range group.Units {
			fg.Units = append(fg.Units, &jdb.FarmGroupUnits{
				UnitId:   unit.Unit.Id,
				TotalQty: unit.TotalQty,
			})
		}
		game.FarmGroups = append(game.FarmGroups, fg)
	}

	for _, unit := range colony.Hull {
		c.Hull = append(c.Hull, &jdb.HullUnit{
			UnitId:   unit.Unit.Id,
			TotalQty: unit.TotalQty,
		})
	}
	sort.Sort(c.Hull)

	for _, unit := range colony.Inventory {
		c.Inventory = append(c.Inventory, &jdb.InventoryUnit{
			UnitId:    unit.Unit.Id,
			TotalQty:  unit.TotalQty,
			StowedQty: unit.StowedQty,
		})
	}
	sort.Sort(c.Inventory)

	for _, group := range colony.MineGroups {
		c.MineGroupIds = append(c.MineGroupIds, group.Id)
		fg := &jdb.MineGroup{
			Id:        group.Id,
			ColonyId:  group.CorS.Id,
			No:        group.No,
			DepositId: group.Deposit.Id,
			UnitId:    group.Unit.Id,
			TotalQty:  group.TotalQty,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		game.MineGroups = append(game.MineGroups, fg)
	}

	c.Pay.ProfessionalPct = colony.Pay.ProfessionalPct
	c.Pay.SoldierPct = colony.Pay.SoldierPct
	c.Pay.UnskilledPct = colony.Pay.UnskilledPct

	c.Population.ProfessionalQty = colony.Population.ProfessionalQty
	c.Population.SoldierQty = colony.Population.SoldierQty
	c.Population.UnskilledQty = colony.Population.UnskilledQty
	c.Population.UnemployedQty = colony.Population.UnemployedQty
	c.Population.ConstructionCrewQty = colony.Population.ConstructionCrewQty
	c.Population.SpyTeamQty = colony.Population.SpyTeamQty
	c.Population.RebelPct = colony.Population.RebelPct
	c.Population.BirthsPriorTurn = colony.Population.BirthsPriorTurn
	c.Population.NaturalDeathsPriorTurn = colony.Population.NaturalDeathsPriorTurn

	c.Rations.ProfessionalPct = colony.Rations.ProfessionalPct
	c.Rations.SoldierPct = colony.Rations.SoldierPct
	c.Rations.UnskilledPct = colony.Rations.UnskilledPct
	c.Rations.UnemployedPct = colony.Rations.UnemployedPct
	return c
}

func wraithColonyToJdbOrbitalColony(colony *wraith.CorS, game *jdb.Game) *jdb.OrbitalColony {
	c := &jdb.OrbitalColony{
		Id:              colony.Id,
		MSN:             colony.MSN,
		BuiltByNationId: colony.BuiltBy.Id,
		Name:            colony.Name,
		TechLevel:       colony.TechLevel,
		PlanetId:        colony.Planet.Id,
	}

	if colony.ControlledBy != nil {
		c.ControlledByPlayerId = colony.ControlledBy.Id
	}

	for _, group := range colony.FactoryGroups {
		c.FactoryGroupIds = append(c.FactoryGroupIds, group.Id)
		fg := &jdb.FactoryGroup{
			Id:        group.Id,
			CorSId:    group.CorS.Id,
			No:        group.No,
			Product:   group.Product.Id,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		for _, unit := range group.Units {
			fg.Units = append(fg.Units, &jdb.FactoryGroupUnits{
				UnitId:   unit.Unit.Id,
				TotalQty: unit.TotalQty,
			})
		}
		game.FactoryGroups = append(game.FactoryGroups, fg)
	}

	for _, group := range colony.FarmGroups {
		c.FarmGroupIds = append(c.FarmGroupIds, group.Id)
		fg := &jdb.FarmGroup{
			Id:        group.Id,
			CorSId:    group.CorS.Id,
			No:        group.No,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		for _, unit := range group.Units {
			fg.Units = append(fg.Units, &jdb.FarmGroupUnits{
				UnitId:   unit.Unit.Id,
				TotalQty: unit.TotalQty,
			})
		}
		game.FarmGroups = append(game.FarmGroups, fg)
	}

	for _, unit := range colony.Hull {
		c.Hull = append(c.Hull, &jdb.HullUnit{
			UnitId:   unit.Unit.Id,
			TotalQty: unit.TotalQty,
		})
	}
	sort.Sort(c.Hull)

	for _, unit := range colony.Inventory {
		c.Inventory = append(c.Inventory, &jdb.InventoryUnit{
			UnitId:    unit.Unit.Id,
			TotalQty:  unit.TotalQty,
			StowedQty: unit.StowedQty,
		})
	}
	sort.Sort(c.Inventory)

	c.Pay.ProfessionalPct = colony.Pay.ProfessionalPct
	c.Pay.SoldierPct = colony.Pay.SoldierPct
	c.Pay.UnskilledPct = colony.Pay.UnskilledPct

	c.Population.ProfessionalQty = colony.Population.ProfessionalQty
	c.Population.SoldierQty = colony.Population.SoldierQty
	c.Population.UnskilledQty = colony.Population.UnskilledQty
	c.Population.UnemployedQty = colony.Population.UnemployedQty
	c.Population.ConstructionCrewQty = colony.Population.ConstructionCrewQty
	c.Population.SpyTeamQty = colony.Population.SpyTeamQty
	c.Population.RebelPct = colony.Population.RebelPct
	c.Population.BirthsPriorTurn = colony.Population.BirthsPriorTurn
	c.Population.NaturalDeathsPriorTurn = colony.Population.NaturalDeathsPriorTurn

	c.Rations.ProfessionalPct = colony.Rations.ProfessionalPct
	c.Rations.SoldierPct = colony.Rations.SoldierPct
	c.Rations.UnskilledPct = colony.Rations.UnskilledPct
	c.Rations.UnemployedPct = colony.Rations.UnemployedPct

	return c
}

func wraithColonyToJdbSurfaceColony(colony *wraith.CorS, game *jdb.Game) *jdb.SurfaceColony {
	c := &jdb.SurfaceColony{
		Id:              colony.Id,
		MSN:             colony.MSN,
		BuiltByNationId: colony.BuiltBy.Id,
		Name:            colony.Name,
		TechLevel:       colony.TechLevel,
		PlanetId:        colony.Planet.Id,
	}

	if colony.ControlledBy != nil {
		c.ControlledByPlayerId = colony.ControlledBy.Id
	}

	for _, group := range colony.FactoryGroups {
		c.FactoryGroupIds = append(c.FactoryGroupIds, group.Id)
		fg := &jdb.FactoryGroup{
			Id:        group.Id,
			CorSId:    group.CorS.Id,
			No:        group.No,
			Product:   group.Product.Id,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		for _, unit := range group.Units {
			fg.Units = append(fg.Units, &jdb.FactoryGroupUnits{
				UnitId:   unit.Unit.Id,
				TotalQty: unit.TotalQty,
			})
		}
		game.FactoryGroups = append(game.FactoryGroups, fg)
	}

	for _, group := range colony.FarmGroups {
		c.FarmGroupIds = append(c.FarmGroupIds, group.Id)
		fg := &jdb.FarmGroup{
			Id:        group.Id,
			CorSId:    group.CorS.Id,
			No:        group.No,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		for _, unit := range group.Units {
			fg.Units = append(fg.Units, &jdb.FarmGroupUnits{
				UnitId:   unit.Unit.Id,
				TotalQty: unit.TotalQty,
			})
		}
		game.FarmGroups = append(game.FarmGroups, fg)
	}

	for _, unit := range colony.Hull {
		c.Hull = append(c.Hull, &jdb.HullUnit{
			UnitId:   unit.Unit.Id,
			TotalQty: unit.TotalQty,
		})
	}
	sort.Sort(c.Hull)

	for _, unit := range colony.Inventory {
		c.Inventory = append(c.Inventory, &jdb.InventoryUnit{
			UnitId:    unit.Unit.Id,
			TotalQty:  unit.TotalQty,
			StowedQty: unit.StowedQty,
		})
	}
	sort.Sort(c.Inventory)

	for _, group := range colony.MineGroups {
		c.MineGroupIds = append(c.MineGroupIds, group.Id)
		fg := &jdb.MineGroup{
			Id:        group.Id,
			ColonyId:  group.CorS.Id,
			No:        group.No,
			DepositId: group.Deposit.Id,
			UnitId:    group.Unit.Id,
			TotalQty:  group.TotalQty,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		game.MineGroups = append(game.MineGroups, fg)
	}

	c.Pay.ProfessionalPct = colony.Pay.ProfessionalPct
	c.Pay.SoldierPct = colony.Pay.SoldierPct
	c.Pay.UnskilledPct = colony.Pay.UnskilledPct

	c.Population.ProfessionalQty = colony.Population.ProfessionalQty
	c.Population.SoldierQty = colony.Population.SoldierQty
	c.Population.UnskilledQty = colony.Population.UnskilledQty
	c.Population.UnemployedQty = colony.Population.UnemployedQty
	c.Population.ConstructionCrewQty = colony.Population.ConstructionCrewQty
	c.Population.SpyTeamQty = colony.Population.SpyTeamQty
	c.Population.RebelPct = colony.Population.RebelPct
	c.Population.BirthsPriorTurn = colony.Population.BirthsPriorTurn
	c.Population.NaturalDeathsPriorTurn = colony.Population.NaturalDeathsPriorTurn

	c.Rations.ProfessionalPct = colony.Rations.ProfessionalPct
	c.Rations.SoldierPct = colony.Rations.SoldierPct
	c.Rations.UnskilledPct = colony.Rations.UnskilledPct
	c.Rations.UnemployedPct = colony.Rations.UnemployedPct

	return c
}

func wraithShipToJdbShip(ship *wraith.CorS, game *jdb.Game) *jdb.Ship {
	s := &jdb.Ship{
		Id:              ship.Id,
		MSN:             ship.MSN,
		BuiltByNationId: ship.BuiltBy.Id,
		Name:            ship.Name,
		TechLevel:       ship.TechLevel,
		PlanetId:        ship.Planet.Id,
	}

	if ship.ControlledBy != nil {
		s.ControlledByPlayerId = ship.ControlledBy.Id
	}

	for _, group := range ship.FactoryGroups {
		s.FactoryGroupIds = append(s.FactoryGroupIds, group.Id)
		fg := &jdb.FactoryGroup{
			Id:        group.Id,
			CorSId:    group.CorS.Id,
			No:        group.No,
			Product:   group.Product.Id,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		for _, unit := range group.Units {
			fg.Units = append(fg.Units, &jdb.FactoryGroupUnits{
				UnitId:   unit.Unit.Id,
				TotalQty: unit.TotalQty,
			})
		}
		game.FactoryGroups = append(game.FactoryGroups, fg)
	}

	for _, group := range ship.FarmGroups {
		s.FarmGroupIds = append(s.FarmGroupIds, group.Id)
		fg := &jdb.FarmGroup{
			Id:        group.Id,
			CorSId:    group.CorS.Id,
			No:        group.No,
			Stage1Qty: group.StageQty[0],
			Stage2Qty: group.StageQty[1],
			Stage3Qty: group.StageQty[2],
			Stage4Qty: group.StageQty[3],
		}
		for _, unit := range group.Units {
			fg.Units = append(fg.Units, &jdb.FarmGroupUnits{
				UnitId:   unit.Unit.Id,
				TotalQty: unit.TotalQty,
			})
		}
		game.FarmGroups = append(game.FarmGroups, fg)
	}

	for _, unit := range ship.Hull {
		s.Hull = append(s.Hull, &jdb.HullUnit{
			UnitId:   unit.Unit.Id,
			TotalQty: unit.TotalQty,
		})
	}
	sort.Sort(s.Hull)

	for _, unit := range ship.Inventory {
		s.Inventory = append(s.Inventory, &jdb.InventoryUnit{
			UnitId:    unit.Unit.Id,
			TotalQty:  unit.TotalQty,
			StowedQty: unit.StowedQty,
		})
	}
	sort.Sort(s.Inventory)

	s.Pay.ProfessionalPct = ship.Pay.ProfessionalPct
	s.Pay.SoldierPct = ship.Pay.SoldierPct
	s.Pay.UnskilledPct = ship.Pay.UnskilledPct

	s.Population.ProfessionalQty = ship.Population.ProfessionalQty
	s.Population.SoldierQty = ship.Population.SoldierQty
	s.Population.UnskilledQty = ship.Population.UnskilledQty
	s.Population.UnemployedQty = ship.Population.UnemployedQty
	s.Population.ConstructionCrewQty = ship.Population.ConstructionCrewQty
	s.Population.SpyTeamQty = ship.Population.SpyTeamQty
	s.Population.RebelPct = ship.Population.RebelPct

	s.Rations.ProfessionalPct = ship.Rations.ProfessionalPct
	s.Rations.SoldierPct = ship.Rations.SoldierPct
	s.Rations.UnskilledPct = ship.Rations.UnskilledPct
	s.Rations.UnemployedPct = ship.Rations.UnemployedPct

	return s
}
