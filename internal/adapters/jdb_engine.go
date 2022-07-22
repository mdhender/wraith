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
	"strings"
	"time"
)

// JdbGameToWraithEngine converts a Game to an Engine.
func JdbGameToWraithEngine(jg *jdb.Game) (*wraith.Engine, error) {
	var err error

	e := &wraith.Engine{
		Version:         "0.1.0",
		Colonies:        make(map[string]*wraith.CorS),
		CorSById:        make(map[int]*wraith.CorS),
		Deposits:        make(map[int]*wraith.Deposit),
		FactoryGroups:   make(map[int]*wraith.FactoryGroup),
		FarmGroups:      make(map[int]*wraith.FarmGroup),
		MineGroups:      make(map[int]*wraith.MineGroup),
		Nations:         make(map[int]*wraith.Nation),
		Planets:         make(map[int]*wraith.Planet),
		Players:         make(map[int]*wraith.Player),
		Ships:           make(map[string]*wraith.CorS),
		Stars:           make(map[int]*wraith.Star),
		Systems:         make(map[int]*wraith.System),
		Units:           make(map[int]*wraith.Unit),
		UnitsFromString: make(map[string]*wraith.Unit),
	}

	e.Game.Id = jg.Id
	e.Game.Code = jg.ShortName
	e.Game.Name = jg.Name
	e.Game.Turn.Year = jg.Turn.Year
	e.Game.Turn.Quarter = jg.Turn.Quarter
	if e.Game.Turn.StartDt, err = time.Parse(time.RFC3339, jg.Turn.StartDt); err != nil {
		return nil, err
	}
	if e.Game.Turn.EndDt, err = time.Parse(time.RFC3339, jg.Turn.EndDt); err != nil {
		return nil, err
	}

	for _, unit := range jg.Units {
		u := jdbUnitToWraithUnit(unit)
		e.Units[unit.Id] = u
		e.UnitsFromString[unit.Code] = u
		e.UnitsFromString[strings.ToLower(u.Name)] = u
	}

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
		s.System.Stars = append(s.System.Stars, s)
	}

	for _, planet := range jg.Planets {
		p := jdbPlanetToWraithPlanet(planet, e.Systems, e.Stars)
		e.Planets[p.Id] = p
		p.Star.Planets = append(p.Star.Planets, p)
	}

	for _, deposit := range jg.Deposits {
		d := jdbDepositToWraithDeposit(deposit, e.Planets, e.CorSById, e.Units)
		e.Deposits[d.Id] = d
		d.Planet.Deposits = append(d.Planet.Deposits, d)
	}

	for _, nation := range jg.Nations {
		n := jdbNationToWraithNation(nation, e.Players, e.Planets)
		e.Nations[n.Id] = n
		e.Players[nation.ControlledByPlayerId].MemberOf = n
	}

	for _, group := range jg.FactoryGroups {
		g := jdbFactoryGroupToWraithFactoryGroup(group, e.CorSById, e.Units)
		e.FactoryGroups[g.Id] = g
	}

	var food *wraith.Unit
	for _, u := range e.Units {
		if u.Code == "FOOD" {
			food = u
			break
		}
	}
	for _, group := range jg.FarmGroups {
		g := jdbFarmGroupToWraithFarmGroup(group, food, e.CorSById, e.Units)
		e.FarmGroups[g.Id] = g
	}

	for _, group := range jg.MineGroups {
		g := jdbMineGroupToWraithMineGroup(group, e.CorSById, e.Deposits, e.Units)
		e.MineGroups[g.Id] = g
	}

	for _, colony := range jg.SurfaceColonies {
		c := jdbSurfaceColonyToWraithColony(colony, e.FactoryGroups, e.FarmGroups, e.MineGroups, e.Nations, e.Planets, e.Players, e.Units)
		e.Colonies[c.HullId] = c
		e.CorSById[c.Id] = c
	}
	for _, colony := range jg.EnclosedColonies {
		c := jdbEnclosedColonyToWraithColony(colony, e.FactoryGroups, e.FarmGroups, e.MineGroups, e.Nations, e.Planets, e.Players, e.Units)
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

	for _, group := range jg.FactoryGroups {
		e.FactoryGroups[group.Id].CorS = e.CorSById[group.CorSId]
	}

	for _, group := range jg.FarmGroups {
		e.FarmGroups[group.Id].CorS = e.CorSById[group.CorSId]
	}

	for _, group := range jg.MineGroups {
		e.MineGroups[group.Id].CorS = e.CorSById[group.ColonyId]
	}

	for _, colony := range e.Colonies {
		if colony.ControlledBy != nil {
			colony.ControlledBy.Colonies = append(colony.ControlledBy.Colonies, colony)
		}
		colony.Planet.Colonies = append(colony.Planet.Colonies, colony)
	}

	for _, ship := range e.Ships {
		if ship.ControlledBy != nil {
			ship.ControlledBy.Colonies = append(ship.ControlledBy.Colonies, ship)
		}
		ship.Planet.Ships = append(ship.Planet.Ships, ship)
	}

	for _, planet := range e.Planets {
		sort.Sort(planet.Colonies)
		sort.Sort(planet.Deposits)
		sort.Sort(planet.Ships)
	}

	for _, player := range e.Players {
		sort.Sort(player.Colonies)
		sort.Sort(player.Ships)
	}

	for _, cs := range e.CorSById {
		if cs.Id > e.Seq {
			e.Seq = cs.Id
		}
	}
	for _, dp := range e.Deposits {
		if dp.Id > e.Seq {
			e.Seq = dp.Id
		}
	}
	for _, fg := range e.FactoryGroups {
		if fg.Id > e.Seq {
			e.Seq = fg.Id
		}
	}
	for _, fg := range e.FarmGroups {
		if fg.Id > e.Seq {
			e.Seq = fg.Id
		}
	}
	for _, mg := range e.MineGroups {
		if mg.Id > e.Seq {
			e.Seq = mg.Id
		}
	}
	for _, n := range e.Nations {
		if n.Id > e.Seq {
			e.Seq = n.Id
		}
	}
	for _, p := range e.Planets {
		if p.Id > e.Seq {
			e.Seq = p.Id
		}
	}
	for _, p := range e.Players {
		if p.Id > e.Seq {
			e.Seq = p.Id
		}
	}
	for _, s := range e.Stars {
		if s.Id > e.Seq {
			e.Seq = s.Id
		}
	}
	for _, s := range e.Systems {
		if s.Id > e.Seq {
			e.Seq = s.Id
		}
	}
	for _, u := range e.Units {
		if u.Id > e.Seq {
			e.Seq = u.Id
		}
	}
	return e, nil
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

func jdbEnclosedColonyToWraithColony(colony *jdb.EnclosedColony, factoryGroup map[int]*wraith.FactoryGroup, farmGroup map[int]*wraith.FarmGroup, mineGroup map[int]*wraith.MineGroup, nations map[int]*wraith.Nation, planets map[int]*wraith.Planet, players map[int]*wraith.Player, units map[int]*wraith.Unit) *wraith.CorS {
	cors := &wraith.CorS{
		Id:           colony.Id,
		Kind:         "enclosed",
		HullId:       fmt.Sprintf("C%d", colony.MSN),
		MSN:          colony.MSN,
		BuiltBy:      nations[colony.BuiltByNationId],
		Name:         colony.Name,
		TechLevel:    colony.TechLevel,
		ControlledBy: players[colony.ControlledByPlayerId],
		Planet:       planets[colony.PlanetId],
		Population: wraith.Population{
			ProfessionalQty:        colony.Population.ProfessionalQty,
			SoldierQty:             colony.Population.SoldierQty,
			UnskilledQty:           colony.Population.UnskilledQty,
			UnemployedQty:          colony.Population.UnemployedQty,
			ConstructionCrewQty:    colony.Population.ConstructionCrewQty,
			SpyTeamQty:             colony.Population.SpyTeamQty,
			RebelPct:               colony.Population.RebelPct,
			BirthsPriorTurn:        colony.Population.BirthsPriorTurn,
			NaturalDeathsPriorTurn: colony.Population.NaturalDeathsPriorTurn,
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
			cors.Hull = append(cors.Hull, &wraith.InventoryUnit{
				Unit:      u,
				ActiveQty: unit.TotalQty,
			})
		}
	}
	for _, unit := range colony.Inventory {
		if u, ok := units[unit.UnitId]; ok {
			switch u.Kind {
			case "consumer-goods", "fuel", "metallics", "military-supplies", "non-metallics":
				cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{Unit: u, StowedQty: unit.TotalQty})
			default:
				cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{
					Unit:      u,
					ActiveQty: unit.TotalQty - unit.StowedQty,
					StowedQty: unit.StowedQty,
				})
			}
		}
	}
	return cors
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
		g.Units = append(g.Units, &wraith.InventoryUnit{
			Unit:      units[u.UnitId],
			ActiveQty: u.TotalQty,
		})
	}
	return g
}

func jdbFarmGroupToWraithFarmGroup(group *jdb.FarmGroup, food *wraith.Unit, cors map[int]*wraith.CorS, units map[int]*wraith.Unit) *wraith.FarmGroup {
	g := &wraith.FarmGroup{
		CorS:     cors[group.CorSId],
		Id:       group.Id,
		No:       group.No,
		Product:  food,
		StageQty: [4]int{group.Stage1Qty, group.Stage2Qty, group.Stage3Qty, group.Stage4Qty},
	}
	for _, u := range group.Units {
		g.Units = append(g.Units, &wraith.InventoryUnit{
			Unit:      units[u.UnitId],
			ActiveQty: u.TotalQty,
		})
	}
	return g
}

func jdbMineGroupToWraithMineGroup(group *jdb.MineGroup, cors map[int]*wraith.CorS, deposits map[int]*wraith.Deposit, units map[int]*wraith.Unit) *wraith.MineGroup {
	g := &wraith.MineGroup{
		CorS:    cors[group.ColonyId],
		Id:      group.Id,
		No:      group.No,
		Deposit: deposits[group.DepositId],
		Unit: &wraith.InventoryUnit{
			Unit:      units[group.UnitId],
			ActiveQty: group.TotalQty,
			StowedQty: 0,
		},
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
			ProfessionalQty:        colony.Population.ProfessionalQty,
			SoldierQty:             colony.Population.SoldierQty,
			UnskilledQty:           colony.Population.UnskilledQty,
			UnemployedQty:          colony.Population.UnemployedQty,
			ConstructionCrewQty:    colony.Population.ConstructionCrewQty,
			SpyTeamQty:             colony.Population.SpyTeamQty,
			RebelPct:               colony.Population.RebelPct,
			BirthsPriorTurn:        colony.Population.BirthsPriorTurn,
			NaturalDeathsPriorTurn: colony.Population.NaturalDeathsPriorTurn,
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
			cors.Hull = append(cors.Hull, &wraith.InventoryUnit{
				Unit:      u,
				ActiveQty: unit.TotalQty,
			})
		}
	}
	for _, unit := range colony.Inventory {
		if u, ok := units[unit.UnitId]; ok {
			switch u.Kind {
			case "consumer-goods", "fuel", "metallics", "military-supplies", "non-metallics":
				cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{Unit: u, StowedQty: unit.TotalQty})
			default:
				cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{
					Unit:      u,
					ActiveQty: unit.TotalQty - unit.StowedQty,
					StowedQty: unit.StowedQty,
				})
			}
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
			ProfessionalQty:        ship.Population.ProfessionalQty,
			SoldierQty:             ship.Population.SoldierQty,
			UnskilledQty:           ship.Population.UnskilledQty,
			UnemployedQty:          ship.Population.UnemployedQty,
			ConstructionCrewQty:    ship.Population.ConstructionCrewQty,
			SpyTeamQty:             ship.Population.SpyTeamQty,
			RebelPct:               ship.Population.RebelPct,
			NaturalDeathsPriorTurn: ship.Population.NaturalDeathsPriorTurn,
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
			cors.Hull = append(cors.Hull, &wraith.InventoryUnit{
				Unit:      u,
				ActiveQty: unit.TotalQty,
			})
		}
	}
	for _, unit := range ship.Inventory {
		if u, ok := units[unit.UnitId]; ok {
			switch u.Kind {
			case "consumer-goods", "fuel", "metallics", "military-supplies", "non-metallics":
				cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{Unit: u, StowedQty: unit.TotalQty})
			default:
				cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{
					Unit:      u,
					ActiveQty: unit.TotalQty - unit.StowedQty,
					StowedQty: unit.StowedQty,
				})
			}
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
			ProfessionalQty:        colony.Population.ProfessionalQty,
			SoldierQty:             colony.Population.SoldierQty,
			UnskilledQty:           colony.Population.UnskilledQty,
			UnemployedQty:          colony.Population.UnemployedQty,
			ConstructionCrewQty:    colony.Population.ConstructionCrewQty,
			SpyTeamQty:             colony.Population.SpyTeamQty,
			RebelPct:               colony.Population.RebelPct,
			BirthsPriorTurn:        colony.Population.BirthsPriorTurn,
			NaturalDeathsPriorTurn: colony.Population.NaturalDeathsPriorTurn,
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
			cors.Hull = append(cors.Hull, &wraith.InventoryUnit{
				Unit:      u,
				ActiveQty: unit.TotalQty,
			})
		}
	}
	for _, unit := range colony.Inventory {
		if u, ok := units[unit.UnitId]; ok {
			switch u.Kind {
			case "consumer-goods", "fuel", "metallics", "military-supplies", "non-metallics":
				cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{Unit: u, StowedQty: unit.TotalQty})
			default:
				cors.Inventory = append(cors.Inventory, &wraith.InventoryUnit{
					Unit:      u,
					ActiveQty: unit.TotalQty - unit.StowedQty,
					StowedQty: unit.StowedQty,
				})
			}
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

func jdbUnitToWraithUnit(unit *jdb.Unit) *wraith.Unit {
	return &wraith.Unit{
		Id:                    unit.Id,
		Kind:                  unit.Kind,
		Code:                  unit.Code,
		TechLevel:             unit.TechLevel,
		Name:                  unit.Name,
		Description:           unit.Description,
		MassPerUnit:           unit.MassPerUnit,
		VolumePerUnit:         unit.VolumePerUnit,
		Hudnut:                unit.Hudnut,
		StowedVolumePerUnit:   unit.StowedVolumePerUnit,
		FuelPerUnitPerTurn:    unit.FuelPerUnitPerTurn,
		MetsPerUnitPerTurn:    unit.MetsPerUnit,
		NonMetsPerUnitPerTurn: unit.NonMetsPerUnit,
	}
}
