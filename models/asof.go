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

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// fetch game by id and turn
func (s *Store) fetchGameByIdAsOf(gameId int, asOfTurn string) (*Game, error) {
	now := time.Now()
	started := now

	// get a transaction with a deferred rollback in case things fail
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: beginTx: %w", err)
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	// fetch game
	game, turn := &Game{}, &Turn{}
	row := tx.QueryRow("select id, short_name, name, descr, year, quarter from games, turns where id = ? and turn = ?", gameId, asOfTurn)
	err = row.Scan(&game.Id, &game.ShortName, &game.Name, &game.Description, &turn.Year, &turn.Quarter)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: game: %w", gameId, err)
	} else if game.Id == 0 {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: game: %w", gameId, ErrNoDataFound)
	}
	game.CurrentTurn = turn
	log.Printf("fetchGameByIdAsOf: %d: game: elapsed %v\n", gameId, time.Now().Sub(started))

	// fetch units
	game.Units = make(map[int]*Unit)
	rows, err := tx.Query("select id, code, tech_level, name, descr, mass_per_unit, volume_per_unit, hudnut, stowed_volume_per_unit from units")
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: units: %w", gameId, err)
	}
	for rows.Next() {
		var hudnut string
		unit := &Unit{}
		err := rows.Scan(&unit.Id, &unit.Code, &unit.TechLevel, &unit.Name, &unit.Description, &unit.MassPerUnit, &unit.VolumePerUnit, &hudnut, &unit.StowedVolumePerUnit)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: units: %w", gameId, err)
		}
		unit.Hudnut = hudnut == "Y"
		game.Units[unit.Id] = unit
	}
	log.Printf("fetchGameByIdAsOf: %d: units: fetched %d units\n", gameId, len(game.Units))
	log.Printf("fetchGameByIdAsOf: %d: units: elapsed %v\n", gameId, time.Now().Sub(started))

	// fetch users
	game.Users = make(map[int]*User)
	rows, err = tx.Query("select u.id, u.handle, up.effdt, up.enddt, up.email, up.handle from users u, user_profile up where up.user_id = u.id and (up.effdt <= ? and ? < up.enddt) order by id", asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: users: %w", gameId, err)
	}
	for rows.Next() {
		user := &User{Profiles: []*UserProfile{{}}}
		err := rows.Scan(&user.Id, &user.Handle, &user.Profiles[0].EffDt, &user.Profiles[0].EndDt, &user.Profiles[0].Email, &user.Profiles[0].Handle)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: users: %w", gameId, err)
		}
		game.Users[user.Id] = user
	}
	log.Printf("fetchGameByIdAsOf: %d: users: fetched %d users\n", gameId, len(game.Users))
	log.Printf("fetchGameByIdAsOf: %d: users: elapsed %v\n", gameId, time.Now().Sub(started))

	// fetch players as of the turn
	game.Players = make(map[int]*Player)
	rows, err = tx.Query("select id from players where game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: players: %s: %w", gameId, asOfTurn, err)
	}
	for rows.Next() {
		player := &Player{Game: game}
		err := rows.Scan(&player.Id)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: players: %s: %w", gameId, asOfTurn, err)
		}
		game.Players[player.Id] = player
	}

	// fetch player details as of the turn
	rows, err = tx.Query(`
		select player_id, handle, controlled_by, ifnull(subject_of, 0)
		from players p, player_dtl pd
		where game_id = ?
		  and pd.player_id = p.id
		  and (efftn <= ? and ? < endtn)
		order by player_id`, gameId, asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: playerDetails: %s: %w", gameId, asOfTurn, err)
	}
	var player *Player
	for rows.Next() {
		var playerId, controlledById, subjectOfId int

		detail := &PlayerDetail{}
		err := rows.Scan(&playerId, &detail.Handle, &controlledById, &subjectOfId)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: playerDetails: %s: %w", gameId, asOfTurn, err)
		}
		if controlledById != 0 {
			detail.ControlledBy = game.Users[controlledById]
		}
		if subjectOfId != 0 {
			detail.SubjectOf = game.Players[subjectOfId]
		}
		if player == nil || player.Id != playerId {
			player = game.Players[playerId]
		}
		detail.Player = player
		player.Details = append(player.Details, detail)
	}
	log.Printf("fetchGameByIdAsOf: %d: players: fetched %d players\n", gameId, len(game.Players))
	log.Printf("fetchGameByIdAsOf: %d: players: elapsed %v\n", gameId, time.Now().Sub(started))

	// fetch nations as of the turn
	game.Nations = make(map[int]*Nation)
	rows, err = tx.Query(`
		select n.id, n.nation_no, n.speciality, n.descr,
			   nd.name, nd.govt_name, nd.govt_kind, nd.controlled_by,
			   nr.tech_level, nr.research_points_pool,
			   ns.biology, ns.bureaucracy, ns.gravitics, ns.life_support, ns.manufacturing, ns.military, ns.mining, ns.shields
		from nations n, nation_dtl nd, nation_research nr, nation_skills ns
		where n.game_id = ?
		  and (nd.nation_id = n.id and nd.efftn <= ? and ? <nd.endtn)
		  and (nr.nation_id = n.id and nr.efftn <= ? and ? <nr.endtn)
		  and (ns.nation_id = n.id and ns.efftn <= ? and ? <ns.endtn)
		order by n.nation_no`, gameId, asOfTurn, asOfTurn, asOfTurn, asOfTurn, asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: nations: %s: %w", gameId, asOfTurn, err)
	}
	for rows.Next() {
		nation := &Nation{Game: game}
		detail := &NationDetail{Nation: nation}
		nation.Details = append(nation.Details, detail)
		research := &NationResearch{Nation: nation}
		nation.Research = append(nation.Research, research)
		skills := &NationSkills{Nation: nation}
		nation.Skills = append(nation.Skills, skills)
		var controlledById int
		err := rows.Scan(&nation.Id, &nation.No, &nation.Speciality, &nation.Description,
			&detail.Name, &detail.GovtName, &detail.GovtKind, &controlledById,
			&research.TechLevel, &research.ResearchPointsPool,
			&skills.Biology, &skills.Bureaucracy, &skills.Gravitics, &skills.LifeSupport, &skills.Manufacturing, &skills.Military, &skills.Mining, &skills.Shields)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: nations: %s: %w", gameId, asOfTurn, err)
		}
		detail.ControlledBy = game.Players[controlledById]
		game.Nations[nation.Id] = nation
	}
	log.Printf("fetchGameByIdAsOf: %d: nations: fetched %d nations\n", gameId, len(game.Nations))
	log.Printf("fetchGameByIdAsOf: %d: nations: elapsed %v\n", gameId, time.Now().Sub(started))

	// cross link nations and players
	rows, err = tx.Query("select nation_id, player_id from nation_player where nation_id in (select id from nations where game_id = ?) order by player_id", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: memberOf: %w", gameId, err)
	}
	numPlayers := 0
	for rows.Next() {
		numPlayers++
		var nationId, playerId int
		err := rows.Scan(&nationId, &playerId)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: memberOf: %w", gameId, err)
		}
		nation, player := game.Nations[nationId], game.Players[playerId]
		player.MemberOf = nation
		nation.Players = append(nation.Players, player)
	}
	log.Printf("fetchGameByIdAsOf: %d: memberOf: fetched %d players\n", gameId, numPlayers)
	log.Printf("fetchGameByIdAsOf: %d: memberOf: elapsed %v\n", gameId, time.Now().Sub(started))

	// fetch systems, stars, and planets
	game.Systems, game.Stars, game.Planets = make(map[int]*System), make(map[int]*Star), make(map[int]*Planet)
	rows, err = tx.Query(`
		select sy.id, st.id, pl.id, x, y, z, st.sequence, st.kind, pl.orbit_no, pl.kind, ifnull(pd.controlled_by, 0), pd.habitability_no
		from systems sy, stars st, planets pl, planet_dtl pd
		where sy.game_id = ?
		  and st.system_id = sy.id
		  and pl.star_id = st.id
		  and pd.planet_id = pl.id
		  and (efftn <= ? and ? < endtn)
		order by x, y, z, sequence, orbit_no`, gameId, asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: systems: %w", gameId, err)
	}
	var system *System
	var star *Star
	for rows.Next() {
		var coords Coordinates
		var systemId, starId int
		var seq, kind string
		var controlledById int

		planet := &Planet{}
		planetDetail := &PlanetDetail{Planet: planet}
		planet.Details = append(planet.Details, planetDetail)

		err := rows.Scan(&systemId, &starId, &planet.Id, &coords.X, &coords.Y, &coords.Z, &seq, &kind, &planet.OrbitNo, &planet.Kind, &controlledById, &planetDetail.HabitabilityNo)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: systems: %w", gameId, err)
		}
		if system == nil || system.Id != systemId {
			system = &System{Game: game, Id: systemId, Coords: coords}
			game.Systems[system.Id] = system
		}
		if star == nil || star.Id != starId {
			star = &Star{System: system, Id: starId, Sequence: seq, Kind: kind}
			game.Stars[star.Id] = star
			system.Stars = append(system.Stars, star)
		}
		planet.Star = star
		star.Orbits = append(star.Orbits, planet)
		if controlledById != 0 {
			planetDetail.ControlledBy = game.Nations[controlledById]
		}
		game.Planets[planet.Id] = planet
	}
	log.Printf("fetchGameByIdAsOf: %d: systems: fetched %8d systems\n", gameId, len(game.Systems))
	log.Printf("fetchGameByIdAsOf: %d: systems: fetched %8d stars\n", gameId, len(game.Stars))
	log.Printf("fetchGameByIdAsOf: %d: systems: fetched %8d planets\n", gameId, len(game.Planets))
	log.Printf("fetchGameByIdAsOf: %d: systems: elapsed %v\n", gameId, time.Now().Sub(started))

	game.Resources = make(map[int]*NaturalResource)
	rows, err = tx.Query(`
		select pl.id,
			   r.id, r.deposit_no, r.kind, r.qty_initial, r.yield_pct,
			   rd.remaining_qty, ifnull(rd.controlled_by, 0)
		from systems sy, stars st, planets pl, resources r, resource_dtl rd
		where sy.game_id = ?
		  and st.system_id = sy.id
		  and pl.star_id = st.id
		  and r.planet_id = pl.id
		  and (rd.resource_id = r.id and rd.efftn <= ? and ? < rd.endtn)
		order by pl.id, r.deposit_no`, gameId, asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: systems: resources: %w", gameId, err)
	}
	var resource *NaturalResource
	numDeposits := 0
	for rows.Next() {
		var planetId, resourceId, depositNo, qtyInitial, qtyRemaining, controlledById int
		var yieldPct float64
		var kind string

		err := rows.Scan(&planetId, &resourceId, &depositNo, &kind, &qtyInitial, &yieldPct, &qtyRemaining, &controlledById)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: systems: resources: %w", gameId, err)
		}
		if resource == nil || resource.Id != resourceId {
			resource = &NaturalResource{Planet: game.Planets[planetId], Id: resourceId, No: depositNo, Kind: kind, QtyInitial: qtyInitial, YieldPct: yieldPct}
			game.Resources[resource.Id] = resource
		}
		detail := &NaturalResourceDetail{NaturalResource: resource, QtyRemaining: qtyRemaining}
		if controlledById != 0 {
			detail.ControlledBy = game.Nations[controlledById]
		}
		numDeposits++
		resource.Details = append(resource.Details, detail)
	}
	log.Printf("fetchGameByIdAsOf: %d: systems: fetched %8d resources\n", gameId, len(game.Resources))
	log.Printf("fetchGameByIdAsOf: %d: systems: fetched %8d deposits\n", gameId, numDeposits)
	log.Printf("fetchGameByIdAsOf: %d: systems: elapsed %v\n", gameId, time.Now().Sub(started))

	// fetch colonies and ships
	game.CorS, game.Colonies, game.Ships = make(map[int]*ColonyOrShip), make(map[int]*ColonyOrShip), make(map[int]*ColonyOrShip)
	rows, err = tx.Query(`
		select cors.id, cors.msn, cors.kind,
			   cdtl.name, cdtl.tech_level, ifnull(cdtl.controlled_by, 0),
			   cloc.planet_id,
			   cpay.professional_pct, cpay.soldier_pct, cpay.unskilled_pct, cpay.unemployed_pct,
			   crat.professional_pct, crat.soldier_pct, crat.unskilled_pct, crat.unemployed_pct,
			   cpop.qty_professional, cpop.qty_soldier, cpop.qty_unskilled, cpop.qty_unemployed, cpop.qty_construction_crews, cpop.qty_spy_teams, cpop.rebel_pct
		from cors, cors_dtl cdtl, cors_loc cloc, cors_pay cpay, cors_rations crat, cors_population cpop
		where cors.game_id = ?
		  and (cdtl.cors_id = cors.id and cdtl.efftn <= ? and ? < cdtl.endtn)
		  and (cloc.cors_id = cors.id and cloc.efftn <= ? and ? < cloc.endtn)
		  and (cpay.cors_id = cors.id and cpay.efftn <= ? and ? < cpay.endtn)
		  and (crat.cors_id = cors.id and crat.efftn <= ? and ? < crat.endtn)
		  and (cpop.cors_id = cors.id and cpop.efftn <= ? and ? < cpop.endtn)
		order by cors.msn`, gameId, asOfTurn, asOfTurn, asOfTurn, asOfTurn, asOfTurn, asOfTurn, asOfTurn, asOfTurn, asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: cors: %w", gameId, err)
	}
	for rows.Next() {
		var corsId, corsMsn, clocPlanetId, cdtlTechLevel, controlledById int
		var corsKind, cdtlName string
		var proPay, solPay, unsPay, unePay, proRat, solRat, unsRat, uneRat, rebelPct float64
		var proQty, solQty, unsQty, uneQty, conQty, spyQty int
		err := rows.Scan(&corsId, &corsMsn, &corsKind,
			&cdtlName, &cdtlTechLevel, &controlledById,
			&clocPlanetId,
			&proPay, &solPay, &unsPay, &unePay,
			&proRat, &solRat, &unsRat, &uneRat,
			&proQty, &solQty, &unsQty, &uneQty, &conQty, &spyQty, &rebelPct)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: cors: %w", gameId, err)
		}
		cors := &ColonyOrShip{
			Game:   game,
			Planet: game.Planets[clocPlanetId],
			Id:     corsId,
			MSN:    corsMsn,
			Kind:   corsKind,
		}
		cors.Details = []*CSDetail{{
			CS:           cors,
			Name:         cdtlName,
			TechLevel:    cdtlTechLevel,
			ControlledBy: game.Players[controlledById],
		}}
		cors.Locations = []*CSLocation{{
			CS:       cors,
			Location: game.Planets[clocPlanetId],
		}}
		cors.Pay = []*CSPay{{
			CS:              cors,
			ProfessionalPct: proPay,
			SoldierPct:      solPay,
			UnskilledPct:    unsPay,
			UnemployedPct:   unePay,
		}}
		cors.Population = []*CSPopulation{{
			CS:                  cors,
			QtyProfessional:     proQty,
			QtySoldier:          solQty,
			QtyUnskilled:        unsQty,
			QtyUnemployed:       uneQty,
			QtyConstructionCrew: conQty,
			QtySpyTeam:          spyQty,
			RebelPct:            rebelPct,
		}}
		cors.Rations = []*CSRations{{
			CS:              cors,
			ProfessionalPct: proRat,
			SoldierPct:      solRat,
			UnskilledPct:    unsRat,
			UnemployedPct:   uneRat,
		}}
		game.CorS[cors.Id] = cors
		if corsKind == "ship" {
			game.Ships[cors.Id] = cors
		} else {
			game.Colonies[cors.Id] = cors
		}
	}
	log.Printf("fetchGameByIdAsOf: %d: cors: fetched %8d cors\n", gameId, len(game.CorS))
	log.Printf("fetchGameByIdAsOf: %d: cors: fetched %8d colonies\n", gameId, len(game.Colonies))
	log.Printf("fetchGameByIdAsOf: %d: cors: fetched %8d ships\n", gameId, len(game.Ships))
	log.Printf("fetchGameByIdAsOf: %d: cors: elapsed %v\n", gameId, time.Now().Sub(started))

	// cross link nations and colonies or ships
	numLinks, numNulls, numShips, numColonies := 0, 0, 0, 0
	for _, cors := range game.CorS {
		if cors.Details == nil || cors.Details[0].ControlledBy == nil {
			numNulls++
			continue
		}
		nation := cors.Details[0].ControlledBy.MemberOf
		if cors.Kind == "ship" {
			numShips++
			nation.Ships = append(nation.Ships, cors)
		} else {
			numColonies++
			nation.Colonies = append(nation.Colonies, cors)
		}
		numLinks++
		nation.CorS = append(nation.CorS, cors)
	}
	log.Printf("fetchGameByIdAsOf: %d: cors: linked  %8d cors\n", gameId, numLinks)
	log.Printf("fetchGameByIdAsOf: %d: cors: linked  %8d colonies\n", gameId, numColonies)
	log.Printf("fetchGameByIdAsOf: %d: cors: linked  %8d ships\n", gameId, numShips)
	log.Printf("fetchGameByIdAsOf: %d: cors: linked  %8d nulls\n", gameId, numNulls)
	log.Printf("fetchGameByIdAsOf: %d: cors: elapsed %v\n", gameId, time.Now().Sub(started))

	// sort cors
	for _, nation := range game.Nations {
		for _, cs := range [][]*ColonyOrShip{nation.CorS, nation.Colonies, nation.Ships} {
			for i := 0; i < len(cs); i++ {
				for j := i + 1; j < len(cs); j++ {
					if cs[j].MSN < cs[i].MSN {
						cs[i], cs[j] = cs[j], cs[i]
					}
				}
			}
		}
	}
	log.Printf("fetchGameByIdAsOf: %d: cors: elapsed %v\n", gameId, time.Now().Sub(started))

	// fetch cors hulls
	rows, err = tx.Query(`
		select cors.id,
			   ch.unit_id, ch.tech_level, ch.qty_operational
		from cors, cors_hull ch
		where cors.game_id = ?
		and (ch.cors_id = cors.id and ch.efftn <= ? and ? < ch.endtn)
		order by cors.id, ch.unit_id, ch.tech_level`, gameId, asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: corsHull: %w", gameId, err)
	}
	var numHulls int
	for rows.Next() {
		numHulls++
		var corsId, unitId, techLevel, qtyOperational int
		err := rows.Scan(&corsId, &unitId, &techLevel, &qtyOperational)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: corsHull: %w", gameId, err)
		}
		hull := &CSHull{
			CS:              game.CorS[corsId],
			Unit:            game.Units[unitId],
			TechLevel:       techLevel,
			QtyOperational:  qtyOperational,
			MassOperational: 0,
			TotalMass:       0,
		}
		hull.CS.Hull = append(hull.CS.Hull, hull)
	}
	log.Printf("fetchGameByIdAsOf: %d: cors: fetched %8d hulls\n", gameId, numHulls)
	log.Printf("fetchGameByIdAsOf: %d: cors: elapsed %v\n", gameId, time.Now().Sub(started))

	// fetch cors inventory
	rows, err = tx.Query(`
		select cors.id,
			   ci.unit_id, ci.tech_level, ci.qty_operational, ci.qty_stowed
		from cors, cors_inventory ci
		where cors.game_id = ?
		and (ci.cors_id = cors.id and ci.efftn <= ? and ? < ci.endtn)
		order by cors.id, ci.unit_id, ci.tech_level`, gameId, asOfTurn, asOfTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByIdAsOf: %d: corsInventory: %w", gameId, err)
	}
	var numInventory int
	for rows.Next() {
		numInventory++
		var corsId, unitId, techLevel, qtyOperational, qtyStowed int
		err := rows.Scan(&corsId, &unitId, &techLevel, &qtyOperational, &qtyStowed)
		if err != nil {
			return nil, fmt.Errorf("fetchGameByIdAsOf: %d: corsInventory: %w", gameId, err)
		}
		inventory := &CSInventory{
			CS:             game.CorS[corsId],
			Unit:           game.Units[unitId],
			TechLevel:      techLevel,
			QtyOperational: qtyOperational,
			QtyStowed:      qtyStowed,
		}
		inventory.CS.Inventory = append(inventory.CS.Inventory, inventory)
	}
	log.Printf("fetchGameByIdAsOf: %d: cors: fetched %8d inventories\n", gameId, numInventory)
	log.Printf("fetchGameByIdAsOf: %d: cors: elapsed %v\n", gameId, time.Now().Sub(started))

	return game, nil
}
