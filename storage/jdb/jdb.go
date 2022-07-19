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

// Package jdb implements a simple data store for game data using JSON files.
package jdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func Extract(db *sql.DB, ctx context.Context, gameId int) (*Game, error) {
	g := &Game{Id: gameId}
	if err := g.extractGame(db); err != nil {
		return nil, fmt.Errorf("jdb: extract: %w", err)
	}
	return g, nil
}

func (g *Game) extractGame(db *sql.DB) error {
	var startDt, endDt time.Time
	row := db.QueryRow(`
		select g.short_name, g.name,
		       t.year, t.quarter, t.start_dt, t.end_dt
			from games g
			inner join turns t on g.id = t.game_id and g.current_turn = t.turn
			where g.id = ?`, g.Id)
	err := row.Scan(&g.ShortName, &g.Name,
		&g.Turn.Year, &g.Turn.Quarter, &startDt, &endDt)
	if err != nil {
		return fmt.Errorf("extractGame: %w", err)
	}
	g.Turn.StartDt = startDt.Format(time.RFC3339)
	g.Turn.EndDt = endDt.Format(time.RFC3339)

	turn := fmt.Sprintf("%04d/%d", g.Turn.Year, g.Turn.Quarter)

	if err = g.extractUnits(db); err != nil {
		return fmt.Errorf("extractGame: %w", err)
	} else if err = g.extractPlayers(db, turn); err != nil {
		return fmt.Errorf("extractGame: %w", err)
	} else if err = g.extractSystems(db, turn); err != nil {
		return fmt.Errorf("extractGame: %w", err)
	} else if err = g.extractNations(db); err != nil {
		return fmt.Errorf("extractGame: %w", err)
	}

	// error correction because i didn't load home world data
	for _, colony := range g.SurfaceColonies {
		for _, nation := range g.Nations {
			if nation.HomePlanetId == 0 && nation.ControlledByPlayerId == colony.ControlledByPlayerId {
				nation.HomePlanetId = colony.PlanetId
				break
			}
		}
	}

	// error correction for built by fields, consumer goods, and life support units
	if g.Turn.Year == 0 && g.Turn.Quarter == 0 {
		var cngd, lsu *Unit
		for _, u := range g.Units {
			switch u.Code {
			case "CNGD":
				cngd = u
			case "LSP-1":
				lsu = u
			}
		}
		for _, colony := range g.EnclosedColonies {
			if colony.BuiltByNationId == 0 {
				for _, player := range g.Players {
					if player.Id == colony.ControlledByPlayerId {
						colony.BuiltByNationId = player.MemberOf
						break
					}
				}
			}
		}
		for _, colony := range g.OrbitalColonies {
			totalPop := colony.Population.ProfessionalQty + colony.Population.SoldierQty + colony.Population.UnskilledQty + colony.Population.UnemployedQty
			totalCngd := (375*colony.Population.ProfessionalQty + 250*colony.Population.SoldierQty + 125*colony.Population.UnskilledQty) / 1000
			if colony.BuiltByNationId == 0 {
				for _, player := range g.Players {
					if player.Id == colony.ControlledByPlayerId {
						colony.BuiltByNationId = player.MemberOf
						break
					}
				}
			}
			for _, u := range colony.Hull {
				if u.UnitId == lsu.Id && u.TotalQty < (totalPop+totalPop/16) {
					u.TotalQty = totalPop + totalPop/16
				}
			}
			for _, u := range colony.Inventory {
				if u.UnitId == cngd.Id && (u.TotalQty+u.StowedQty)/4 < (totalCngd+totalCngd/16) {
					u.TotalQty = 4 * (totalCngd + totalCngd/16)
					u.StowedQty = u.TotalQty
				}
			}
		}
		for _, colony := range g.SurfaceColonies {
			if colony.BuiltByNationId == 0 {
				for _, player := range g.Players {
					if player.Id == colony.ControlledByPlayerId {
						colony.BuiltByNationId = player.MemberOf
						break
					}
				}
			}
		}
	}

	return nil
}

func (g *Game) extractUnits(db *sql.DB) error {
	rows, err := db.Query(`
		select id, code, tech_level, name, descr, mass_per_unit, volume_per_unit, hudnut, stowed_volume_per_unit
		from units
		order by id`)
	if err != nil {
		return fmt.Errorf("extractUnits: %w", err)
	}
	for rows.Next() {
		var hudnut string
		unit := &Unit{}
		err := rows.Scan(&unit.Id, &unit.Code, &unit.TechLevel, &unit.Name, &unit.Description, &unit.MassPerUnit, &unit.VolumePerUnit, &hudnut, &unit.StowedVolumePerUnit)
		if err != nil {
			return fmt.Errorf("extractUnits: %w", err)
		}
		unit.Kind = unit.Description
		unit.Hudnut = hudnut == "Y"

		unit.MetsPerUnit, unit.NonMetsPerUnit, _, unit.FuelPerUnitPerTurn, _ = unitAttributes(unit.Kind, unit.TechLevel)

		g.Units = append(g.Units, unit)
	}
	return nil
}

func (g *Game) extractPlayers(db *sql.DB, turn string) error {
	rows, err := db.Query(`
		select u.id, p.id, pd.handle, np.nation_id, ifnull(pd.subject_of, 0)
		from users u
			left join player_dtl pd on u.id = pd.controlled_by and (pd.efftn <= ? and ? < pd.endtn)
			left join players p on pd.player_id = p.id
		    left join nation_player np on p.id = np.player_id
		where p.game_id = ?`, turn, turn, g.Id)
	if err != nil {
		return fmt.Errorf("extractPlayers: %w", err)
	}
	for rows.Next() {
		player := &Player{}
		err := rows.Scan(&player.UserId, &player.Id, &player.Name, &player.MemberOf, &player.ReportsToPlayerId)
		if err != nil {
			return fmt.Errorf("extractPlayers: %w", err)
		}
		g.Players = append(g.Players, player)
	}
	return nil
}

func (g *Game) extractNations(db *sql.DB) error {
	rows, err := db.Query(`
		select n.id, n.nation_no, n.speciality,
			   nd.name, nd.govt_name, nd.govt_kind, nd.controlled_by,
			   nr.tech_level, nr.research_points_pool
		from games g
			inner join nations n on n.game_id = g.id
			inner join nation_dtl nd on n.id = nd.nation_id and (nd.efftn <= g.current_turn and g.current_turn < nd.endtn)
			inner join nation_research nr on n.id = nr.nation_id and (nr.efftn <= g.current_turn and g.current_turn < nr.endtn)
		where g.id = ?
		order by n.nation_no`, g.Id)
	if err != nil {
		return fmt.Errorf("extractNations: %w", err)
	}
	for rows.Next() {
		nation := &Nation{}
		err := rows.Scan(&nation.Id, &nation.No, &nation.Speciality,
			&nation.Name, &nation.GovtName, &nation.GovtKind, &nation.ControlledByPlayerId,
			&nation.TechLevel, &nation.ResearchPointsPool)
		if err != nil {
			return fmt.Errorf("extractNations: %w", err)
		}
		g.Nations = append(g.Nations, nation)
	}
	return nil
}

func (g *Game) extractSystems(db *sql.DB, turn string) error {
	rows, err := db.Query(`
		select s.id, s.x, s.y, s.z
		from systems s
		where s.game_id = ?
		order by s.id`, g.Id)
	if err != nil {
		return fmt.Errorf("extractSystems: %w", err)
	}
	for rows.Next() {
		system := &System{}
		err := rows.Scan(&system.Id, &system.Coords.X, &system.Coords.Y, &system.Coords.Z)
		if err != nil {
			return fmt.Errorf("extractSystems: %w", err)
		} else if err = g.extractStars(db, system, turn); err != nil {
			return fmt.Errorf("extractSystems: %w", err)
		}
		g.Systems = append(g.Systems, system)
	}
	return nil
}

func (g *Game) extractStars(db *sql.DB, system *System, turn string) error {
	rows, err := db.Query(`
		select s.id, s.sequence, s.kind
		from stars s
		where s.system_id = ?
		order by s.id`, system.Id)
	if err != nil {
		return fmt.Errorf("extractStars: %w", err)
	}
	for rows.Next() {
		star := &Star{SystemId: system.Id}
		err := rows.Scan(&star.Id, &star.Sequence, &star.Kind)
		if err != nil {
			return fmt.Errorf("extractStars: %w", err)
		} else if err = g.extractPlanets(db, star, turn); err != nil {
			return fmt.Errorf("extractStars: %w", err)
		}
		system.StarIds = append(system.StarIds, star.Id)
		g.Stars = append(g.Stars, star)
	}
	return nil
}

func (g *Game) extractPlanets(db *sql.DB, star *Star, turn string) error {
	rows, err := db.Query(`
		select p.id, p.orbit_no, p.kind, ifnull(pd.habitability_no, 0)
		from planets p
			left join planet_dtl pd on p.id = pd.planet_id and (pd.efftn <= ? and ? < pd.endtn)
		where p.star_id = ?
		order by p.orbit_no`, turn, turn, star.Id)
	if err != nil {
		return fmt.Errorf("extractPlanets: %w", err)
	}
	for rows.Next() {
		planet := &Planet{SystemId: star.SystemId, StarId: star.Id}
		err := rows.Scan(&planet.Id, &planet.OrbitNo, &planet.Kind, &planet.HabitabilityNo)
		if err != nil {
			return fmt.Errorf("extractPlanets: %w", err)
		} else if err = g.extractDeposits(db, planet, turn); err != nil {
			return fmt.Errorf("extractPlanets: %w", err)
		} else if err = g.extractSurfaceColonies(db, planet, turn); err != nil {
			return fmt.Errorf("extractPlanets: %w", err)
		} else if err = g.extractEnclosedColonies(db, planet, turn); err != nil {
			return fmt.Errorf("extractPlanets: %w", err)
		} else if err = g.extractOrbitalColonies(db, planet, turn); err != nil {
			return fmt.Errorf("extractPlanets: %w", err)
		} else if err = g.extractShips(db, planet, turn); err != nil {
			return fmt.Errorf("extractPlanets: %w", err)
		}
		star.PlanetIds = append(star.PlanetIds, planet.Id)
		g.Planets = append(g.Planets, planet)
	}
	return nil
}

func (g *Game) extractDeposits(db *sql.DB, planet *Planet, turn string) error {
	rows, err := db.Query(`
		select r.id, r.deposit_no, r.unit_id, r.qty_initial, r.yield_pct,
			   ifnull(rd.remaining_qty, 0), ifnull(rd.controlled_by, 0)
		from resources r
			left join resource_dtl rd on r.id = rd.resource_id and (rd.efftn <= ? and ? < rd.endtn)
		where r.planet_id = ?
		order by r.deposit_no`, turn, turn, planet.Id)
	if err != nil {
		return fmt.Errorf("extractDeposits: %w", err)
	}
	for rows.Next() {
		deposit := &Deposit{PlanetId: planet.Id}
		err := rows.Scan(&deposit.Id, &deposit.No, &deposit.UnitId, &deposit.InitialQty, &deposit.YieldPct,
			&deposit.RemainingQty, &deposit.ControlledByColonyId)
		if err != nil {
			return fmt.Errorf("extractDeposits: %w", err)
		}
		planet.DepositIds = append(planet.DepositIds, deposit.Id)
		g.Deposits = append(g.Deposits, deposit)
	}
	return nil
}

func (g *Game) extractSurfaceColonies(db *sql.DB, planet *Planet, turn string) error {
	rows, err := db.Query(`
		select c.id, c.msn,
			   cd.name, cd.tech_level, ifnull(cd.controlled_by, 0),
			   cl.planet_id,
			   cp.qty_professional, cp.qty_soldier, cp.qty_unskilled, cp.qty_unemployed, cp.qty_construction_crews, cp.qty_spy_teams, cp.rebel_pct,
			   cpp.professional_pct, cpp.soldier_pct, cpp.unskilled_pct,
			   cr.professional_pct, cr.soldier_pct, cr.unskilled_pct, cr.unemployed_pct
		from games g
			inner join cors c on g.id = c.game_id and c.kind = 'open'
			inner join cors_dtl cd on c.id = cd.cors_id and (cd.efftn <= g.current_turn and g.current_turn < cd.endtn)
			inner join cors_loc cl on c.id = cl.cors_id and (cl.efftn <= g.current_turn and g.current_turn < cl.endtn)
			inner join cors_population cp on c.id = cp.cors_id and (cp.efftn <= g.current_turn and g.current_turn < cp.endtn)
			inner join cors_pay cpp on c.id = cpp.cors_id and (cpp.efftn <= g.current_turn and g.current_turn < cpp.endtn)
			inner join cors_rations cr on c.id = cr.cors_id and (cr.efftn <= g.current_turn and g.current_turn < cr.endtn)
		where g.id = ?
		and cl.planet_id = ?
		order by c.id`, g.Id, planet.Id)
	if err != nil {
		return fmt.Errorf("extractSurfaceColonies: %w", err)
	}
	for rows.Next() {
		colony := &SurfaceColony{}
		err := rows.Scan(&colony.Id, &colony.MSN,
			&colony.Name, &colony.TechLevel, &colony.ControlledByPlayerId,
			&colony.PlanetId,
			&colony.Population.ProfessionalQty, &colony.Population.SoldierQty, &colony.Population.UnskilledQty, &colony.Population.UnemployedQty, &colony.Population.ConstructionCrewQty, &colony.Population.SpyTeamQty, &colony.Population.RebelPct,
			&colony.Pay.ProfessionalPct, &colony.Pay.SoldierPct, &colony.Pay.UnskilledPct,
			&colony.Rations.ProfessionalPct, &colony.Rations.SoldierPct, &colony.Rations.UnskilledPct, &colony.Rations.UnemployedPct)
		if err != nil {
			return fmt.Errorf("extractSurfaceColonies: %w", err)
		} else if colony.Hull, err = g.extractHull(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractSurfaceColonies: %w", err)
		} else if colony.Inventory, err = g.extractInventory(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractSurfaceColonies: %w", err)
		}
		if groups, err := g.extractFactoryGroups(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractSurfaceColonies: %w", err)
		} else {
			for _, group := range groups {
				colony.FactoryGroupIds = append(colony.FactoryGroupIds, group.Id)
			}
		}
		if groups, err := g.extractFarmGroups(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractSurfaceColonies: %w", err)
		} else {
			for _, group := range groups {
				colony.FarmGroupIds = append(colony.FarmGroupIds, group.Id)
			}
		}
		if groups, err := g.extractMineGroups(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractSurfaceColonies: %w", err)
		} else {
			for _, group := range groups {
				colony.MineGroupIds = append(colony.MineGroupIds, group.Id)
			}
		}
		planet.SurfaceColonyIds = append(planet.SurfaceColonyIds, colony.Id)
		g.SurfaceColonies = append(g.SurfaceColonies, colony)
	}
	return nil
}

func (g *Game) extractHull(db *sql.DB, corsId int, turn string) ([]*HullUnit, error) {
	var units []*HullUnit
	rows, err := db.Query(`
		select unit_id, qty_operational
		from cors_hull
		where cors_id = ?
		  and (efftn <= ? and ? < endtn)
		order by unit_id`, corsId, turn, turn)
	if err != nil {
		return nil, fmt.Errorf("extractHull: %w", err)
	}
	for rows.Next() {
		unit := &HullUnit{}
		if err := rows.Scan(&unit.UnitId, &unit.TotalQty); err != nil {
			return nil, fmt.Errorf("extractHull: %w", err)
		}
		units = append(units, unit)
	}
	return units, nil
}

func (g *Game) extractInventory(db *sql.DB, corsId int, turn string) ([]*InventoryUnit, error) {
	var units []*InventoryUnit
	rows, err := db.Query(`
		select unit_id, qty_operational, qty_stowed
		from cors_inventory
		where cors_id = ?
		  and (efftn <= ? and ? < endtn)
		order by unit_id`, corsId, turn, turn)
	if err != nil {
		return nil, fmt.Errorf("extractInventory: %w", err)
	}
	for rows.Next() {
		var operQty int
		unit := &InventoryUnit{}
		if err := rows.Scan(&unit.UnitId, &operQty, &unit.StowedQty); err != nil {
			return nil, fmt.Errorf("extractInventory: %w", err)
		}
		unit.TotalQty = operQty + unit.StowedQty
		units = append(units, unit)
	}
	return units, nil
}

func (g *Game) extractFactoryGroups(db *sql.DB, corsId int, turn string) ([]*FactoryGroup, error) {
	var groups []*FactoryGroup
	rows, err := db.Query(`
		select c.id, c.group_no, c.unit_id,
			   ifnull(cfgs.qty_stage_1, 0), ifnull(cfgs.qty_stage_2, 0), ifnull(cfgs.qty_stage_3, 0), ifnull(cfgs.qty_stage_4, 0)
		from cors_factory_group c
			left join cors_factory_group_stages cfgs on c.id = cfgs.factory_group_id and cfgs.turn = ?
		where c.cors_id = ?
		  and (c.efftn <= ? and ? < c.endtn)
		order by c.group_no, c.unit_id`, turn, corsId, turn, turn)
	if err != nil {
		return nil, fmt.Errorf("extractFactoryGroups: %w", err)
	}
	for rows.Next() {
		group := &FactoryGroup{CorSId: corsId}
		if err := rows.Scan(&group.Id, &group.No, &group.Product,
			&group.Stage1Qty, &group.Stage2Qty, &group.Stage3Qty, &group.Stage4Qty); err != nil {
			return nil, fmt.Errorf("extractFactoryGroups: %w", err)
		} else if group.Units, err = g.extractFactoryGroupUnits(db, group.Id, turn); err != nil {
			return nil, fmt.Errorf("extractFactoryGroups: %w", err)
		}
		groups = append(groups, group)
		g.FactoryGroups = append(g.FactoryGroups, group)
	}
	return groups, nil
}

func (g *Game) extractFactoryGroupUnits(db *sql.DB, groupId int, turn string) ([]*FactoryGroupUnits, error) {
	var units []*FactoryGroupUnits
	rows, err := db.Query(`
		select unit_id, qty_operational
		from cors_factory_group_units
		where factory_group_id = ?
		  and (efftn <= ? and ? < endtn)
		order by unit_id`, groupId, turn, turn)
	if err != nil {
		return nil, fmt.Errorf("extractFactoryGroupUnits: %w", err)
	}
	for rows.Next() {
		unit := &FactoryGroupUnits{}
		if err := rows.Scan(&unit.UnitId, &unit.TotalQty); err != nil {
			return nil, fmt.Errorf("extractFactoryGroupUnits: %w", err)
		}
		units = append(units, unit)
	}
	return units, nil
}

func (g *Game) extractFarmGroups(db *sql.DB, corsId int, turn string) ([]*FarmGroup, error) {
	var groups []*FarmGroup
	rows, err := db.Query(`
		select c.id, c.group_no,
			   ifnull(cfgs.qty_stage_1, 0), ifnull(cfgs.qty_stage_2, 0), ifnull(cfgs.qty_stage_3, 0), ifnull(cfgs.qty_stage_4, 0)
		from cors_farm_group c
			left join cors_farm_group_stages cfgs on c.id = cfgs.farm_group_id and cfgs.turn = ?
		where c.cors_id = ?
		  and (c.efftn <= ? and ? < c.endtn)
		order by c.group_no, c.unit_id`, turn, corsId, turn, turn)
	if err != nil {
		return nil, fmt.Errorf("extractFarmGroups: %w", err)
	}
	for rows.Next() {
		group := &FarmGroup{CorSId: corsId}
		if err := rows.Scan(&group.Id, &group.No,
			&group.Stage1Qty, &group.Stage2Qty, &group.Stage3Qty, &group.Stage4Qty); err != nil {
			return nil, fmt.Errorf("extractFarmGroups: %w", err)
		} else if group.Units, err = g.extractFarmGroupUnits(db, group.Id, turn); err != nil {
			return nil, fmt.Errorf("extractFarmGroups: %w", err)
		}
		groups = append(groups, group)
		g.FarmGroups = append(g.FarmGroups, group)
	}
	return groups, nil
}

func (g *Game) extractFarmGroupUnits(db *sql.DB, groupId int, turn string) ([]*FarmGroupUnits, error) {
	var units []*FarmGroupUnits
	rows, err := db.Query(`
		select unit_id, qty_operational
		from cors_farm_group_units
		where farm_group_id = ?
		  and (efftn <= ? and ? < endtn)
		order by unit_id`, groupId, turn, turn)
	if err != nil {
		return nil, fmt.Errorf("extractFarmGroupUnits: %w", err)
	}
	for rows.Next() {
		unit := &FarmGroupUnits{}
		if err := rows.Scan(&unit.UnitId, &unit.TotalQty); err != nil {
			return nil, fmt.Errorf("extractFarmGroupUnits: %w", err)
		}
		units = append(units, unit)
	}
	return units, nil
}

func (g *Game) extractMineGroups(db *sql.DB, corsId int, turn string) ([]*MineGroup, error) {
	var groups []*MineGroup
	rows, err := db.Query(`
		select c.id, c.group_no, c.resource_id,
			   cmgu.unit_id, cmgu.qty_operational,
			   ifnull(cmgs.qty_stage_1, 0), ifnull(cmgs.qty_stage_2, 0), ifnull(cmgs.qty_stage_3, 0), ifnull(cmgs.qty_stage_4, 0)
		from cors_mining_group c
			left join cors_mining_group_units cmgu on c.id = cmgu.mining_group_id and (cmgu.efftn <= ? and ? < cmgu.endtn)
			left join cors_mining_group_stages cmgs on c.id = cmgs.mining_group_id and cmgs.turn = ?
		where c.cors_id = ?
		  and (c.efftn <= ? and ? < c.endtn)
		order by c.group_no`, turn, turn, turn, corsId, turn, turn)
	if err != nil {
		return nil, fmt.Errorf("extractMineGroups: %w", err)
	}
	for rows.Next() {
		group := &MineGroup{ColonyId: corsId}
		if err := rows.Scan(&group.Id, &group.No, &group.DepositId,
			&group.UnitId, &group.TotalQty,
			&group.Stage1Qty, &group.Stage2Qty, &group.Stage3Qty, &group.Stage4Qty); err != nil {
			return nil, fmt.Errorf("extractMineGroups: %w", err)
		}
		groups = append(groups, group)
		g.MineGroups = append(g.MineGroups, group)
	}
	return groups, nil
}

func (g *Game) extractEnclosedColonies(db *sql.DB, planet *Planet, turn string) error {
	rows, err := db.Query(`
		select c.id, c.msn,
			   cd.name, cd.tech_level, ifnull(cd.controlled_by, 0),
			   cl.planet_id,
			   cp.qty_professional, cp.qty_soldier, cp.qty_unskilled, cp.qty_unemployed, cp.qty_construction_crews, cp.qty_spy_teams, cp.rebel_pct,
			   cpp.professional_pct, cpp.soldier_pct, cpp.unskilled_pct,
			   cr.professional_pct, cr.soldier_pct, cr.unskilled_pct, cr.unemployed_pct
		from games g
			inner join cors c on g.id = c.game_id and c.kind = 'enclosed'
			inner join cors_dtl cd on c.id = cd.cors_id and (cd.efftn <= g.current_turn and g.current_turn < cd.endtn)
			inner join cors_loc cl on c.id = cl.cors_id and (cl.efftn <= g.current_turn and g.current_turn < cl.endtn)
			inner join cors_population cp on c.id = cp.cors_id and (cp.efftn <= g.current_turn and g.current_turn < cp.endtn)
			inner join cors_pay cpp on c.id = cpp.cors_id and (cpp.efftn <= g.current_turn and g.current_turn < cpp.endtn)
			inner join cors_rations cr on c.id = cr.cors_id and (cr.efftn <= g.current_turn and g.current_turn < cr.endtn)
		where g.id = ?
		and cl.planet_id = ?
		order by c.id`, g.Id, planet.Id)
	if err != nil {
		return fmt.Errorf("extractEnclosedColonies: %w", err)
	}
	for rows.Next() {
		colony := &EnclosedColony{}
		err := rows.Scan(&colony.Id, &colony.MSN,
			&colony.Name, &colony.TechLevel, &colony.ControlledByPlayerId,
			&colony.PlanetId,
			&colony.Population.ProfessionalQty, &colony.Population.SoldierQty, &colony.Population.UnskilledQty, &colony.Population.UnemployedQty, &colony.Population.ConstructionCrewQty, &colony.Population.SpyTeamQty, &colony.Population.RebelPct,
			&colony.Pay.ProfessionalPct, &colony.Pay.SoldierPct, &colony.Pay.UnskilledPct,
			&colony.Rations.ProfessionalPct, &colony.Rations.SoldierPct, &colony.Rations.UnskilledPct, &colony.Rations.UnemployedPct)
		if err != nil {
			return fmt.Errorf("extractEnclosedColonies: %w", err)
		} else if colony.Hull, err = g.extractHull(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractEnclosedColonies: %w", err)
		} else if colony.Inventory, err = g.extractInventory(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractEnclosedColonies: %w", err)
		}
		if groups, err := g.extractFactoryGroups(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractEnclosedColonies: %w", err)
		} else {
			for _, group := range groups {
				colony.FactoryGroupIds = append(colony.FactoryGroupIds, group.Id)
			}
		}
		if groups, err := g.extractFarmGroups(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractEnclosedColonies: %w", err)
		} else {
			for _, group := range groups {
				colony.FarmGroupIds = append(colony.FarmGroupIds, group.Id)
			}
		}
		if groups, err := g.extractMineGroups(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractEnclosedColonies: %w", err)
		} else {
			for _, group := range groups {
				colony.MineGroupIds = append(colony.MineGroupIds, group.Id)
			}
		}
		planet.EnclosedColonyIds = append(planet.EnclosedColonyIds, colony.Id)
		g.EnclosedColonies = append(g.EnclosedColonies, colony)
	}
	return nil
}

func (g *Game) extractOrbitalColonies(db *sql.DB, planet *Planet, turn string) error {
	rows, err := db.Query(`
		select c.id, c.msn,
			   cd.name, cd.tech_level, ifnull(cd.controlled_by, 0),
			   cl.planet_id,
			   cp.qty_professional, cp.qty_soldier, cp.qty_unskilled, cp.qty_unemployed, cp.qty_construction_crews, cp.qty_spy_teams, cp.rebel_pct,
			   cpp.professional_pct, cpp.soldier_pct, cpp.unskilled_pct,
			   cr.professional_pct, cr.soldier_pct, cr.unskilled_pct, cr.unemployed_pct
		from games g
			inner join cors c on g.id = c.game_id and c.kind = 'orbital'
			inner join cors_dtl cd on c.id = cd.cors_id and (cd.efftn <= g.current_turn and g.current_turn < cd.endtn)
			inner join cors_loc cl on c.id = cl.cors_id and (cl.efftn <= g.current_turn and g.current_turn < cl.endtn)
			inner join cors_population cp on c.id = cp.cors_id and (cp.efftn <= g.current_turn and g.current_turn < cp.endtn)
			inner join cors_pay cpp on c.id = cpp.cors_id and (cpp.efftn <= g.current_turn and g.current_turn < cpp.endtn)
			inner join cors_rations cr on c.id = cr.cors_id and (cr.efftn <= g.current_turn and g.current_turn < cr.endtn)
		where g.id = ?
		and cl.planet_id = ?
		order by c.id`, g.Id, planet.Id)
	if err != nil {
		return fmt.Errorf("extractOrbitalColonies: %w", err)
	}
	for rows.Next() {
		colony := &OrbitalColony{}
		err := rows.Scan(&colony.Id, &colony.MSN,
			&colony.Name, &colony.TechLevel, &colony.ControlledByPlayerId,
			&colony.PlanetId,
			&colony.Population.ProfessionalQty, &colony.Population.SoldierQty, &colony.Population.UnskilledQty, &colony.Population.UnemployedQty, &colony.Population.ConstructionCrewQty, &colony.Population.SpyTeamQty, &colony.Population.RebelPct,
			&colony.Pay.ProfessionalPct, &colony.Pay.SoldierPct, &colony.Pay.UnskilledPct,
			&colony.Rations.ProfessionalPct, &colony.Rations.SoldierPct, &colony.Rations.UnskilledPct, &colony.Rations.UnemployedPct)
		if err != nil {
			return fmt.Errorf("extractOrbitalColonies: %w", err)
		} else if colony.Hull, err = g.extractHull(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractOrbitalColonies: %w", err)
		} else if colony.Inventory, err = g.extractInventory(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractOrbitalColonies: %w", err)
		}
		if groups, err := g.extractFactoryGroups(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractOrbitalColonies: %w", err)
		} else {
			for _, group := range groups {
				colony.FactoryGroupIds = append(colony.FactoryGroupIds, group.Id)
			}
		}
		if groups, err := g.extractFarmGroups(db, colony.Id, turn); err != nil {
			return fmt.Errorf("extractOrbitalColonies: %w", err)
		} else {
			for _, group := range groups {
				colony.FarmGroupIds = append(colony.FarmGroupIds, group.Id)
			}
		}
		planet.OrbitalColonyIds = append(planet.OrbitalColonyIds, colony.Id)
		g.OrbitalColonies = append(g.OrbitalColonies, colony)
	}
	return nil
}

func (g *Game) extractShips(db *sql.DB, planet *Planet, turn string) error {
	rows, err := db.Query(`
		select c.id, c.msn,
			   cd.name, cd.tech_level, ifnull(cd.controlled_by, 0),
			   cl.planet_id,
			   cp.qty_professional, cp.qty_soldier, cp.qty_unskilled, cp.qty_unemployed, cp.qty_construction_crews, cp.qty_spy_teams, cp.rebel_pct,
			   cpp.professional_pct, cpp.soldier_pct, cpp.unskilled_pct,
			   cr.professional_pct, cr.soldier_pct, cr.unskilled_pct, cr.unemployed_pct
		from games g
			inner join cors c on g.id = c.game_id and c.kind = 'ship'
			inner join cors_dtl cd on c.id = cd.cors_id and (cd.efftn <= g.current_turn and g.current_turn < cd.endtn)
			inner join cors_loc cl on c.id = cl.cors_id and (cl.efftn <= g.current_turn and g.current_turn < cl.endtn)
			inner join cors_population cp on c.id = cp.cors_id and (cp.efftn <= g.current_turn and g.current_turn < cp.endtn)
			inner join cors_pay cpp on c.id = cpp.cors_id and (cpp.efftn <= g.current_turn and g.current_turn < cpp.endtn)
			inner join cors_rations cr on c.id = cr.cors_id and (cr.efftn <= g.current_turn and g.current_turn < cr.endtn)
		where g.id = ?
		and cl.planet_id = ?
		order by c.id`, g.Id, planet.Id)
	if err != nil {
		return fmt.Errorf("extractShips: %w", err)
	}
	for rows.Next() {
		ship := &Ship{}
		err := rows.Scan(&ship.Id, &ship.MSN,
			&ship.Name, &ship.TechLevel, &ship.ControlledByPlayerId,
			&ship.PlanetId,
			&ship.Population.ProfessionalQty, &ship.Population.SoldierQty, &ship.Population.UnskilledQty, &ship.Population.UnemployedQty, &ship.Population.ConstructionCrewQty, &ship.Population.SpyTeamQty, &ship.Population.RebelPct,
			&ship.Pay.ProfessionalPct, &ship.Pay.SoldierPct, &ship.Pay.UnskilledPct,
			&ship.Rations.ProfessionalPct, &ship.Rations.SoldierPct, &ship.Rations.UnskilledPct, &ship.Rations.UnemployedPct)
		if err != nil {
			return fmt.Errorf("extractShips: %w", err)
		} else if ship.Hull, err = g.extractHull(db, ship.Id, turn); err != nil {
			return fmt.Errorf("extractShips: %w", err)
		} else if ship.Inventory, err = g.extractInventory(db, ship.Id, turn); err != nil {
			return fmt.Errorf("extractShips: %w", err)
		}
		if groups, err := g.extractFactoryGroups(db, ship.Id, turn); err != nil {
			return fmt.Errorf("extractShips: %w", err)
		} else {
			for _, group := range groups {
				ship.FactoryGroupIds = append(ship.FactoryGroupIds, group.Id)
			}
		}
		if groups, err := g.extractFarmGroups(db, ship.Id, turn); err != nil {
			return fmt.Errorf("extractShips: %w", err)
		} else {
			for _, group := range groups {
				ship.FarmGroupIds = append(ship.FarmGroupIds, group.Id)
			}
		}
		planet.ShipIds = append(planet.ShipIds, ship.Id)
		g.Ships = append(g.Ships, ship)
	}
	return nil
}
