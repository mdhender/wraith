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
	"fmt"
	"github.com/mdhender/wraith/storage/jdb"
	"github.com/mdhender/wraith/wraith"
	"sort"
)

// JdbGameToWraithEngine converts a Game to an Engine.
func JdbGameToWraithEngine(jg *jdb.Game) *wraith.Engine {
	e := &wraith.Engine{
		Version:       "0.1.0",
		Colonies:      make(map[string]*wraith.CorS),
		CorSById:      make(map[int]*wraith.CorS),
		Deposits:      make(map[int]*wraith.Deposit),
		FactoryGroups: make(map[int]*wraith.FactoryGroup),
		FarmGroups:    make(map[int]*wraith.FarmGroup),
		MineGroups:    make(map[int]*wraith.MineGroup),
		Nations:       make(map[int]*wraith.Nation),
		Planets:       make(map[int]*wraith.Planet),
		Players:       make(map[int]*wraith.Player),
		Ships:         make(map[string]*wraith.CorS),
		Stars:         make(map[int]*wraith.Star),
		Systems:       make(map[int]*wraith.System),
		Units:         make(map[int]*wraith.Unit),
	}

	e.Game.Code = jg.ShortName
	e.Game.Turn.Year = jg.Turn.Year
	e.Game.Turn.Quarter = jg.Turn.Quarter

	// two loops to create players.
	// first loop creates the struct.
	for _, player := range jg.Players {
		e.Players[player.Id] = &wraith.Player{
			Id:     player.Id,
			UserId: player.UserId,
			Name:   player.Name,
		}
	}
	// second loop links players to rulers.
	for _, player := range jg.Players {
		if player.ReportsToPlayerId != 0 {
			e.Players[player.Id].ReportsTo = e.Players[player.ReportsToPlayerId]
		}
	}

	for _, system := range jg.Systems {
		s := jdbSystemToWraithSystem(system)
		e.Systems[s.Id] = s
	}

	for _, star := range jg.Stars {
		s := jdbStarToWraithStar(star, e.Systems)
		e.Stars[s.Id] = s
	}

	for _, planet := range jg.Planets {
		p := jdbPlanetToWraithPlanet(planet, e.Systems, e.Stars)
		e.Planets[p.Id] = p
	}

	for _, nation := range jg.Nations {
		n := jdbNationToWraithNation(nation, e.Players, e.Planets)
		e.Nations[n.Id] = n
		e.Players[nation.ControlledByPlayerId].MemberOf = n
	}

	for _, colony := range jg.SurfaceColonies {
		c := jdbSurfaceColonyToWraithColony(colony, e.FactoryGroups, e.FarmGroups, e.MineGroups, e.Nations, e.Planets, e.Players, e.Units)
		e.Colonies[c.HullId] = c
		e.CorSById[c.Id] = c
	}
	for _, colony := range jg.OrbitalColonies {
		c := jdbOrbitalColonyToWraithColony(colony, e.FactoryGroups, e.FarmGroups, e.Nations, e.Planets, e.Players, e.Units)
		e.Colonies[c.HullId] = c
		e.CorSById[c.Id] = c
	}
	for _, ship := range jg.Ships {
		s := jdbShipToWraithShip(ship, e.FactoryGroups, e.FarmGroups, e.Nations, e.Planets, e.Players, e.Units)
		e.Ships[s.HullId] = s
		e.CorSById[s.Id] = s
	}

	for _, deposit := range jg.Deposits {
		d := jdbDepositToWraithDeposit(deposit, e.Planets, e.CorSById, e.Units)
		e.Deposits[d.Id] = d
		d.Planet.Deposits = append(d.Planet.Deposits, d)
	}

	for _, group := range jg.FactoryGroups {
		g := jdbFactoryGroupToWraithFactoryGroup(group, e.CorSById, e.Units)
		e.FactoryGroups[g.Id] = g
	}

	for _, group := range jg.FarmGroups {
		g := jdbFarmGroupToWraithFarmGroup(group, e.CorSById, e.Units)
		e.FarmGroups[g.Id] = g
	}

	for _, group := range jg.MineGroups {
		g := jdbMineGroupToWraithMineGroup(group, e.CorSById, e.Deposits, e.Units)
		e.MineGroups[g.Id] = g
	}

	for _, planet := range e.Planets {
		sort.Sort(planet.Colonies)
		sort.Sort(planet.Deposits)
		sort.Sort(planet.Ships)
	}

	return e
}

func jdbDepositToWraithDeposit(deposit *jdb.Deposit, planets map[int]*wraith.Planet, cors map[int]*wraith.CorS, units map[int]*wraith.Unit) *wraith.Deposit {
	return &wraith.Deposit{
		Id:           deposit.Id,
		No:           deposit.No,
		Product:      units[deposit.UnitId],
		InitialQty:   deposit.InitialQty,
		RemainingQty: deposit.RemainingQty,
		YieldPct:     deposit.YieldPct,
		Planet:       planets[deposit.PlanetId],
		ControlledBy: cors[deposit.ControlledByColonyId],
	}
}

func jdbFactoryGroupToWraithFactoryGroup(group *jdb.FactoryGroup, cors map[int]*wraith.CorS, units map[int]*wraith.Unit) *wraith.FactoryGroup {
	g := &wraith.FactoryGroup{
		CorS:     cors[group.CorSId],
		Id:       group.Id,
		No:       group.No,
		Product:  units[group.Product],
		StageQty: [4]int{group.Stage1Qty, group.Stage2Qty, group.Stage3Qty, group.Stage4Qty},
	}
	for _, u := range group.Units {
		g.Units = append(g.Units, &wraith.FactoryGroupUnits{
			Unit:     units[u.UnitId],
			TotalQty: u.TotalQty,
		})
	}
	return g
}

func jdbFarmGroupToWraithFarmGroup(group *jdb.FarmGroup, cors map[int]*wraith.CorS, units map[int]*wraith.Unit) *wraith.FarmGroup {
	g := &wraith.FarmGroup{
		CorS:     cors[group.CorSId],
		Id:       group.Id,
		No:       group.No,
		StageQty: [4]int{group.Stage1Qty, group.Stage2Qty, group.Stage3Qty, group.Stage4Qty},
	}
	for _, u := range group.Units {
		g.Units = append(g.Units, &wraith.FarmGroupUnits{
			Unit:     units[u.UnitId],
			TotalQty: u.TotalQty,
		})
	}
	return g
}

func jdbMineGroupToWraithMineGroup(group *jdb.MineGroup, cors map[int]*wraith.CorS, deposits map[int]*wraith.Deposit, units map[int]*wraith.Unit) *wraith.MineGroup {
	g := &wraith.MineGroup{
		CorS:     cors[group.ColonyId],
		Id:       group.Id,
		No:       group.No,
		Deposit:  deposits[group.DepositId],
		Unit:     units[group.UnitId],
		TotalQty: group.TotalQty,
		StageQty: [4]int{group.Stage1Qty, group.Stage2Qty, group.Stage3Qty, group.Stage4Qty},
	}
	return g
}

func jdbNationToWraithNation(nation *jdb.Nation, players map[int]*wraith.Player, planets map[int]*wraith.Planet) *wraith.Nation {
	return &wraith.Nation{
		Id:                 nation.Id,
		No:                 nation.No,
		Name:               nation.Name,
		GovtName:           nation.GovtName,
		GovtKind:           nation.GovtKind,
		HomePlanet:         planets[nation.HomePlanetId],
		ControlledBy:       players[nation.ControlledByPlayerId],
		Speciality:         nation.Speciality,
		TechLevel:          nation.TechLevel,
		ResearchPointsPool: nation.ResearchPointsPool,
	}
}

func jdbOrbitalColonyToWraithColony(colony *jdb.OrbitalColony, factoryGroup map[int]*wraith.FactoryGroup, farmGroup map[int]*wraith.FarmGroup, nations map[int]*wraith.Nation, planets map[int]*wraith.Planet, players map[int]*wraith.Player, units map[int]*wraith.Unit) *wraith.CorS {
	cors := &wraith.CorS{
		Id:           colony.Id,
		Kind:         "orbital",
		HullId:       fmt.Sprintf("C%d", colony.MSN),
		MSN:          colony.MSN,
		BuiltBy:      nations[colony.BuiltByNationId],
		Name:         colony.Name,
		TechLevel:    colony.TechLevel,
		ControlledBy: players[colony.ControlledByPlayerId],
		Planet:       planets[colony.PlanetId],
		Population: wraith.Population{
			ProfessionalQty:     colony.Population.ProfessionalQty,
			SoldierQty:          colony.Population.SoldierQty,
			UnskilledQty:        colony.Population.UnskilledQty,
			UnemployedQty:       colony.Population.UnemployedQty,
			ConstructionCrewQty: colony.Population.ConstructionCrewQty,
			SpyTeamQty:          colony.Population.SpyTeamQty,
			RebelPct:            colony.Population.RebelPct,
		},
		Pay: wraith.Pay{
			ProfessionalPct: colony.Pay.ProfessionalPct,
			SoldierPct:      colony.Pay.SoldierPct,
			UnskilledPct:    colony.Pay.UnskilledPct,
		},
		Rations: wraith.Rations{
			ProfessionalPct: colony.Rations.ProfessionalPct,
			SoldierPct:      colony.Rations.SoldierPct,
			UnskilledPct:    colony.Rations.UnskilledPct,
			UnemployedPct:   colony.Rations.UnemployedPct,
		},
	}
	for _, group := range colony.FactoryGroupIds {
		cors.FactoryGroups = append(cors.FactoryGroups, factoryGroup[group])
	}
	for _, group := range colony.FarmGroupIds {
		cors.FarmGroups = append(cors.FarmGroups, farmGroup[group])
	}
	for _, unit := range colony.Hull {
		if u, ok := units[unit.UnitId]; ok {
			cors.Hull = append(cors.Hull, &wraith.HullUnit{
				Unit:     u,
				TotalQty: unit.TotalQty,
			})
		}
	}
	for _, unit := range colony.Inventory {
		if u, ok := units[unit.UnitId]; ok {
			cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{
				Unit:      u,
				TotalQty:  unit.TotalQty,
				StowedQty: unit.StowedQty,
			})
		}
	}
	return cors
}

func jdbPlanetToWraithPlanet(planet *jdb.Planet, systems map[int]*wraith.System, stars map[int]*wraith.Star) *wraith.Planet {
	return &wraith.Planet{
		Id:             planet.Id,
		System:         systems[planet.SystemId],
		Star:           stars[planet.StarId],
		OrbitNo:        planet.OrbitNo,
		Kind:           planet.Kind,
		HabitabilityNo: planet.HabitabilityNo,
	}
}

func jdbShipToWraithShip(ship *jdb.Ship, factoryGroup map[int]*wraith.FactoryGroup, farmGroup map[int]*wraith.FarmGroup, nations map[int]*wraith.Nation, planets map[int]*wraith.Planet, players map[int]*wraith.Player, units map[int]*wraith.Unit) *wraith.CorS {
	cors := &wraith.CorS{
		Id:           ship.Id,
		Kind:         "ship",
		HullId:       fmt.Sprintf("S%d", ship.MSN),
		MSN:          ship.MSN,
		BuiltBy:      nations[ship.BuiltByNationId],
		Name:         ship.Name,
		TechLevel:    ship.TechLevel,
		ControlledBy: players[ship.ControlledByPlayerId],
		Planet:       planets[ship.PlanetId],
		Population: wraith.Population{
			ProfessionalQty:     ship.Population.ProfessionalQty,
			SoldierQty:          ship.Population.SoldierQty,
			UnskilledQty:        ship.Population.UnskilledQty,
			UnemployedQty:       ship.Population.UnemployedQty,
			ConstructionCrewQty: ship.Population.ConstructionCrewQty,
			SpyTeamQty:          ship.Population.SpyTeamQty,
			RebelPct:            ship.Population.RebelPct,
		},
		Pay: wraith.Pay{
			ProfessionalPct: ship.Pay.ProfessionalPct,
			SoldierPct:      ship.Pay.SoldierPct,
			UnskilledPct:    ship.Pay.UnskilledPct,
		},
		Rations: wraith.Rations{
			ProfessionalPct: ship.Rations.ProfessionalPct,
			SoldierPct:      ship.Rations.SoldierPct,
			UnskilledPct:    ship.Rations.UnskilledPct,
			UnemployedPct:   ship.Rations.UnemployedPct,
		},
	}
	for _, group := range ship.FactoryGroupIds {
		cors.FactoryGroups = append(cors.FactoryGroups, factoryGroup[group])
	}
	for _, group := range ship.FarmGroupIds {
		cors.FarmGroups = append(cors.FarmGroups, farmGroup[group])
	}
	for _, unit := range ship.Hull {
		if u, ok := units[unit.UnitId]; ok {
			cors.Hull = append(cors.Hull, &wraith.HullUnit{
				Unit:     u,
				TotalQty: unit.TotalQty,
			})
		}
	}
	for _, unit := range ship.Inventory {
		if u, ok := units[unit.UnitId]; ok {
			cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{
				Unit:      u,
				TotalQty:  unit.TotalQty,
				StowedQty: unit.StowedQty,
			})
		}
	}
	return cors
}

func jdbSurfaceColonyToWraithColony(colony *jdb.SurfaceColony, factoryGroup map[int]*wraith.FactoryGroup, farmGroup map[int]*wraith.FarmGroup, mineGroup map[int]*wraith.MineGroup, nations map[int]*wraith.Nation, planets map[int]*wraith.Planet, players map[int]*wraith.Player, units map[int]*wraith.Unit) *wraith.CorS {
	cors := &wraith.CorS{
		Id:           colony.Id,
		Kind:         "surface",
		HullId:       fmt.Sprintf("C%d", colony.MSN),
		MSN:          colony.MSN,
		BuiltBy:      nations[colony.BuiltByNationId],
		Name:         colony.Name,
		TechLevel:    colony.TechLevel,
		ControlledBy: players[colony.ControlledByPlayerId],
		Planet:       planets[colony.PlanetId],
		Population: wraith.Population{
			ProfessionalQty:     colony.Population.ProfessionalQty,
			SoldierQty:          colony.Population.SoldierQty,
			UnskilledQty:        colony.Population.UnskilledQty,
			UnemployedQty:       colony.Population.UnemployedQty,
			ConstructionCrewQty: colony.Population.ConstructionCrewQty,
			SpyTeamQty:          colony.Population.SpyTeamQty,
			RebelPct:            colony.Population.RebelPct,
		},
		Pay: wraith.Pay{
			ProfessionalPct: colony.Pay.ProfessionalPct,
			SoldierPct:      colony.Pay.SoldierPct,
			UnskilledPct:    colony.Pay.UnskilledPct,
		},
		Rations: wraith.Rations{
			ProfessionalPct: colony.Rations.ProfessionalPct,
			SoldierPct:      colony.Rations.SoldierPct,
			UnskilledPct:    colony.Rations.UnskilledPct,
			UnemployedPct:   colony.Rations.UnemployedPct,
		},
	}
	for _, group := range colony.FactoryGroupIds {
		cors.FactoryGroups = append(cors.FactoryGroups, factoryGroup[group])
	}
	for _, group := range colony.FarmGroupIds {
		cors.FarmGroups = append(cors.FarmGroups, farmGroup[group])
	}
	for _, group := range colony.MineGroupIds {
		cors.MineGroups = append(cors.MineGroups, mineGroup[group])
	}
	for _, unit := range colony.Hull {
		if u, ok := units[unit.UnitId]; ok {
			cors.Hull = append(cors.Hull, &wraith.HullUnit{
				Unit:     u,
				TotalQty: unit.TotalQty,
			})
		}
	}
	for _, unit := range colony.Inventory {
		if u, ok := units[unit.UnitId]; ok {
			cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{
				Unit:      u,
				TotalQty:  unit.TotalQty,
				StowedQty: unit.StowedQty,
			})
		}
	}
	return cors
}

func jdbStarToWraithStar(star *jdb.Star, systems map[int]*wraith.System) *wraith.Star {
	return &wraith.Star{
		Id:       star.Id,
		System:   systems[star.SystemId],
		Sequence: star.Sequence,
		Kind:     star.Kind,
	}
}

func jdbSystemToWraithSystem(system *jdb.System) *wraith.System {
	return &wraith.System{
		Id: system.Id,
		Coords: wraith.Coordinates{
			X: system.Coords.X,
			Y: system.Coords.Y,
			Z: system.Coords.Z,
		},
	}
}
