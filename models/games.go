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
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"
)

// CreateGame adds a new game to the store if it passes validation
func (s *Store) CreateGame(g *Game) error {
	if s.db == nil {
		return ErrNoConnection
	}

	g.ShortName = strings.ToUpper(strings.TrimSpace(g.ShortName))
	if g.ShortName == "" {
		return fmt.Errorf("short name: %w", ErrMissingField)
	}
	for _, r := range g.ShortName { // check for invalid runes in the field
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return fmt.Errorf("short name: invalid rune: %w", ErrInvalidField)
		}
	}

	// check for duplicate keys
	var count int
	row := s.db.QueryRow("select ifnull(count(id), 0) from games where short_name = ?", g.ShortName)
	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("createGame: count: %w", err)
	} else if count != 0 {
		return ErrDuplicateKey
	}

	g.Name = strings.TrimSpace(g.Name)
	if g.Name == "" {
		g.Name = g.ShortName
	}
	g.Description = strings.TrimSpace(g.Description)
	if g.Description == "" {
		g.Description = g.Name
	}

	return nil
}

// DeleteGame removes a game from the store.
func (s *Store) DeleteGame(id int) error {
	return s.deleteGame(id)
}

// DeleteGameByName removes a game from the store.
func (s *Store) DeleteGameByName(shortName string) error {
	return s.deleteGameByName(shortName)
}

// FetchGame fetches a game by id
func (s *Store) FetchGame(id int) (*Game, error) {
	return s.fetchGame(id)
}

// FetchGameByName does just that
func (s *Store) FetchGameByName(name string) (*Game, error) {
	return s.fetchGameByName(name)
}

func (s *Store) GenerateGame(shortName, name, descr string, radius int, startDt time.Time, positions []*PlayerPosition) (*Game, error) {
	return s.genGame(shortName, name, descr, radius, startDt, positions)
}

// LookupGame looks up a game by id
func (s *Store) LookupGame(id int) (*Game, error) {
	return s.lookupGame(id)
}

// LookupGameByName does just that
func (s *Store) LookupGameByName(name string) (*Game, error) {
	return s.lookupGameByName(name)
}

func (s *Store) SaveGame(game *Game) error {
	return s.saveGame(game)
}

func (s *Store) deleteGame(id int) error {
	if s.db == nil {
		return ErrNoConnection
	}
	_, err := s.db.Exec("delete from games where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) deleteGameByName(shortName string) error {
	if s.db == nil {
		return ErrNoConnection
	}
	_, err := s.db.Exec("delete from games where short_name = ?", shortName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) fetchGame(id int) (*Game, error) {
	// get a transaction with a deferred rollback in case things fail
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: beginTx: %w", err)
	}
	defer tx.Rollback()

	g := &Game{}

	row := s.db.QueryRow("select id, short_name, name, descr, current_turn from games where id = ?", id)
	var currentTurn string
	err = row.Scan(&g.Id, &g.ShortName, &g.Name, &g.Description, &currentTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: %w", id, err)
	}

	rows, err := s.db.Query("select no, year, quarter, start_dt, end_dt from turns where game_id = ? order by turn", g.Id)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: turns: %w", id, err)
	}
	for rows.Next() {
		turn := &Turn{}
		err := rows.Scan(&turn.No, &turn.Year, &turn.Quarter, &turn.StartDt, &turn.EndDt)
		if err != nil {
			return nil, fmt.Errorf("fetchGame: turns: %w", err)
		}
		g.Turns = append(g.Turns, turn)
		if currentTurn == turn.String() {
			g.CurrentTurn = turn
		}
	}

	rows, err = s.db.Query("select id, ifnull(controlled_by, 0), ifnull(subject_of, 0) from players where game_id = ? order by id", g.Id)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: players: %w", id, err)
	}
	for rows.Next() {
		player := &Player{Game: g}
		var controlledBy, subjectOf int
		err := rows.Scan(&player.Id, &controlledBy, &subjectOf)
		if err != nil {
			return nil, fmt.Errorf("fetchGame: players: %+v", err)
		}
		log.Printf("fetchGame: player %d controlled_by %d subject_of %d\n", player.Id, controlledBy, subjectOf)
		g.Players = append(g.Players, player)
	}

	rows, err = s.db.Query("select id, nation_no, speciality, descr from nations where game_id = ? order by nation_no", g.Id)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: nations: %w", id, err)
	}
	for rows.Next() {
		nation := &Nation{Game: g}
		err := rows.Scan(&nation.Id, &nation.No, &nation.Speciality, &nation.Description)
		if err != nil {
			return nil, fmt.Errorf("fetchGame: nations: %+v", err)
		}
		g.Nations = append(g.Nations, nation)
	}

	log.Printf("fetchGame: need to link nations to players\n")

	return g, tx.Commit()
}

func (s *Store) fetchGameByName(name string) (*Game, error) {
	name = strings.ToUpper(strings.TrimSpace(name))
	var id int
	row := s.db.QueryRow("select id from games where short_name = ?", name)
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByName: %q: %w", name, err)
	}
	return s.fetchGame(id)
}

func (s *Store) genGame(shortName, name, descr string, radius int, startDt time.Time, positions []*PlayerPosition) (*Game, error) {
	shortName = strings.ToUpper(strings.TrimSpace(shortName))
	if shortName == "" {
		return nil, fmt.Errorf("short name: %w", ErrMissingField)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = shortName
	}
	descr = strings.TrimSpace(descr)
	if descr == "" {
		descr = shortName
	}

	// delete values
	err := s.deleteGameByName(shortName)
	if err != nil {
		return nil, fmt.Errorf("createGame: %w", err)
	}

	effTurn, endTurn := &Turn{}, &Turn{Year: 9999, Quarter: 4}

	game := &Game{
		ShortName:   shortName,
		Name:        name,
		Description: descr,
	}

	// convert positions to players
	for _, position := range positions {
		user, err := s.fetchUserByHandle(position.UserHandle)
		if err != nil {
			return nil, fmt.Errorf("user: %q: %w", position.UserHandle, ErrNoSuchUser)
		}
		player := &Player{
			Game:         game,
			ControlledBy: user,
			SubjectOf:    nil,
			Details:      nil,
		}
		player.Details = []*PlayerDetail{{
			Player:  player,
			EffTurn: effTurn,
			EndTurn: endTurn,
			Handle:  position.PlayerHandle,
		}}
		game.Players = append(game.Players, player)
		log.Printf("createGame: created position %q\n", position.PlayerHandle)
	}

	systemsPerRing := len(positions)
	totalSystems := radius * systemsPerRing
	log.Printf("createGame: systems per ring %3d estimated systems %6d\n", systemsPerRing, totalSystems)
	rings := mkrings(radius, systemsPerRing)
	numPoints := 0
	for d := 0; d <= radius; d++ {
		numPoints += len(rings[d])
		log.Printf("createGame: ring %2d: %5d\n", d, len(rings[d]))
	}
	log.Printf("createGame:   total: %5d\n", numPoints)

	turnNo, turnDuration := 0, 2*7*24*time.Hour // assume two-week turns
	effDt := startDt
	endDt := effDt.Add(turnDuration)
	game.Turns = append(game.Turns, &Turn{No: turnNo, Year: 0, Quarter: 0, StartDt: effDt, EndDt: endDt})
	effDt = endDt
	endDt = effDt.Add(turnDuration)
	turnNo++
	for year := 1; year <= 10; year++ {
		for quarter := 1; quarter <= 4; quarter++ {
			game.Turns = append(game.Turns, &Turn{No: turnNo, Year: year, Quarter: quarter, StartDt: effDt, EndDt: endDt})
			effDt = endDt
			endDt = effDt.Add(turnDuration)
			turnNo++
		}
	}

	systemId, ring, colonyNo := 0, 5, 0

	// generate nations and their home systems
	for no, position := range positions {
		// warning: assumes that player was created for this game
		player := game.Players[no]

		systemId++
		coords := rings[ring][0]
		rings[ring] = rings[ring][1:]

		system := s.genHomeSystem(systemId)
		system.Ring, system.Coords = ring, coords
		game.Systems = append(game.Systems, system)

		planet := system.Stars[0].Orbits[3]
		nation := s.genNation(no+1, planet, player, position)
		colonyNo++
		nation.Colonies[0].MSN = colonyNo
		colonyNo++
		nation.Colonies[1].MSN = colonyNo

		game.Nations = append(game.Nations, nation)
	}

	// generate the remainder of the systems
	for ring := 0; ring < len(rings); ring++ {
		for _, coords := range rings[ring] {
			systemId++
			system := s.genSystem(systemId)
			system.Ring, system.Coords = ring, coords
			game.Systems = append(game.Systems, system)
		}
	}

	return game, nil
}

func (s *Store) lookupGame(id int) (*Game, error) {
	row := s.db.QueryRow("select id, short_name, name, descr, current_turn from games where id = ?", id)
	var g Game
	var currentTurn string
	err := row.Scan(&g.Id, &g.ShortName, &g.Name, &g.Description, &currentTurn)
	if err != nil {
		return nil, fmt.Errorf("lookupGame: %d: %w", id, err)
	}

	return &g, nil
}

func (s *Store) lookupGameByName(name string) (*Game, error) {
	name = strings.ToUpper(strings.TrimSpace(name))
	var id int
	row := s.db.QueryRow("select id from games where short_name = ?", name)
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("lookupGameByName: %q: %w", name, err)
	}
	return s.lookupGame(id)
}

func (s *Store) saveGame(g *Game) error {
	// get a transaction with a deferred rollback in case things fail
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("saveGame: beginTx: %w", err)
	}
	defer tx.Rollback()

	if g.CurrentTurn == nil {
		g.CurrentTurn = g.Turns[0]
	}

	r, err := tx.ExecContext(s.ctx, "insert into games (short_name, name, current_turn, descr) values (?, ?, ?, ?)",
		g.ShortName, g.Name, g.CurrentTurn.String(), g.Description)
	if err != nil {
		return fmt.Errorf("saveGame: games: insert: %w", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("saveGame: games: fetchId: %w", err)
	}
	g.Id = int(id)
	log.Printf("created game %3d %s\n", g.Id, g.ShortName)

	for _, turn := range g.Turns {
		_, err = tx.ExecContext(s.ctx, "insert into turns (game_id, no, year, quarter, turn, start_dt, end_dt) values (?, ?, ?, ?, ?, ?, ?)",
			g.Id, turn.No, turn.Year, turn.Quarter, turn.String(), turn.StartDt, turn.EndDt)
		if err != nil {
			return fmt.Errorf("saveGame: turns: insert: %w", err)
		}
	}

	var nobody int
	row := tx.QueryRow("select ifnull(id, 0) from users where handle = ?", "nobody")
	err = row.Scan(&nobody)
	if err != nil {
		return fmt.Errorf("saveGame: users: nobody: %w", err)
	}

	for _, player := range g.Players {
		var uid int
		row := tx.QueryRow("select ifnull(id, 0) from users where handle = ?", player.ControlledBy.Handle)
		err = row.Scan(&uid)
		if err != nil {
			return fmt.Errorf("saveGame: users: %w", err)
		} else if uid == 0 {
			return fmt.Errorf("saveGame: users: %q: %w", player.ControlledBy.Handle, ErrNoSuchUser)
		}
		log.Printf("saveGame: mapped %8d to player %q\n", uid, player.Details[0].Handle)

		r, err := tx.ExecContext(s.ctx, "insert into players (game_id, controlled_by, subject_of) values (?, ?, null)",
			g.Id, uid)
		if err != nil {
			return fmt.Errorf("saveGame: players: insert: %w", err)
		}
		id, err := r.LastInsertId()
		if err != nil {
			return fmt.Errorf("saveGame: players: lastInsertId: %w", err)
		}
		player.Id = int(id)
	}

	for _, system := range g.Systems {
		r, err := tx.ExecContext(s.ctx, "insert into systems (game_id, x, y, z, qty_stars) values (?, ?, ?, ?, ?)",
			g.Id, system.Coords.X, system.Coords.Y, system.Coords.Z, len(system.Stars))
		if err != nil {
			return fmt.Errorf("saveGame: systems: insert: %w", err)
		}
		id, err := r.LastInsertId()
		if err != nil {
			return fmt.Errorf("saveGame: systems: lastInsertId: %w", err)
		}
		system.Id = int(id)

		for _, star := range system.Stars {
			r, err := tx.ExecContext(s.ctx, "insert into stars (system_id, sequence, kind) values (?, ?, ?)",
				system.Id, star.Sequence, star.Kind)
			if err != nil {
				return fmt.Errorf("saveGame: stars: insert: %w", err)
			}
			id, err := r.LastInsertId()
			if err != nil {
				return fmt.Errorf("saveGame: stars: lastInsertId: %w", err)
			}
			star.Id = int(id)

			for orbit, planet := range star.Orbits {
				if orbit == 0 {
					continue
				}
				homePlanet := "N"
				if planet.HomePlanet {
					homePlanet = "Y"
				}
				r, err := tx.ExecContext(s.ctx, "insert into planets (star_id, orbit_no, kind, habitability_no, home_planet) values (?, ?, ?, ?, ?)",
					star.Id, planet.OrbitNo, planet.Kind, planet.HabitabilityNo, homePlanet)
				if err != nil {
					return fmt.Errorf("saveGame: planet: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: planet: lastInsertId: %w", err)
				}
				planet.Id = int(id)

				for n, deposit := range planet.Deposits {
					deposit.No = n + 1
					r, err := tx.ExecContext(s.ctx, "insert into resources (planet_id, deposit_no, kind, qty_initial, yield_pct) values (?, ?, ?, ?, ?)",
						planet.Id, deposit.No, deposit.Kind, deposit.QtyInitial, deposit.YieldPct)
					if err != nil {
						log.Printf("failed  system %8d: star %8d: orbit %2d: planet %8d: resource %8d %s\n", system.Id, star.Id, planet.OrbitNo, planet.Id, deposit.Id, deposit.Kind)
						return fmt.Errorf("saveGame: deposit: insert: %w", err)
					}
					id, err := r.LastInsertId()
					if err != nil {
						return fmt.Errorf("saveGame: deposit: lastInsertId: %w", err)
					}
					deposit.Id = int(id)
					//log.Printf("created system %8d: star %8d: orbit %2d: planet %8d: resource %8d %-13s %9d\n", system.Id, star.Id, planet.Orbit, planet.Id, resource.Id, resource.Kind, resource.InitialQuantity)
				}
				//log.Printf("created system %8d: star %8d: orbit %2d: planet %8d\n", system.Id, star.Id, orbit, planet.Id)
			}
			//log.Printf("created system %8d: star %8d: suffix %q\n", system.Id, star.Id, star.Suffix)
		}
		//log.Printf("created system %8d\n", system.Id)
	}

	for _, nation := range g.Nations {
		r, err := tx.ExecContext(s.ctx, "insert into nations (game_id, nation_no, speciality, descr) values (?, ?, ?, ?)",
			g.Id, nation.No, nation.Speciality, nation.Description)
		if err != nil {
			return fmt.Errorf("saveGame: nations: insert: %w", err)
		}
		id, err := r.LastInsertId()
		if err != nil {
			return fmt.Errorf("saveGame: nations: lastInsertId: %w", err)
		}
		nation.Id = int(id)
		_, err = tx.ExecContext(s.ctx, "insert into nation_dtl (nation_id, efftn, endtn, name, govt_name, govt_kind, controlled_by) values (?, ?, ?, ?, ?, ?, ?)",
			nation.Id, nation.Details[0].EffTurn.String(), nation.Details[0].EndTurn.String(), nation.Details[0].Name, nation.Details[0].GovtName, nation.Details[0].GovtKind, nation.Details[0].ControlledBy.Id)
		if err != nil {
			return fmt.Errorf("saveGame: nation_dtl: insert: %w", err)
		}
		_, err = tx.ExecContext(s.ctx, "insert into nation_skills (nation_id, efftn, endtn, tech_level, research_points_pool, biology, bureaucracy, gravitics, life_support, manufacturing, military, mining, shields) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			nation.Id, nation.Research[0].EffTurn.String(), nation.Research[0].EndTurn.String(), nation.Research[0].TechLevel, nation.Research[0].ResearchPointsPool, nation.Skills[0].Biology, nation.Skills[0].Bureaucracy, nation.Skills[0].Gravitics, nation.Skills[0].LifeSupport, nation.Skills[0].Manufacturing, nation.Skills[0].Military, nation.Skills[0].Mining, nation.Skills[0].Shields)
		if err != nil {
			return fmt.Errorf("saveGame: nation_skills: insert: %w", err)
		}
		log.Printf("created nation %3d %8d\n", nation.No, nation.Id)
	}

	for _, nation := range g.Nations {
		for _, colony := range nation.Colonies {
			if colony.Kind == "ship" {
				continue
			}
			r, err := tx.ExecContext(s.ctx, "insert into cors (game_id, msn, kind, planet_id) values (?, ?, ?, ?)",
				g.Id, colony.MSN, colony.Kind, colony.Details[0].Location.Id)
			if err != nil {
				return fmt.Errorf("saveGame: colonies: insert: %w", err)
			}
			id, err := r.LastInsertId()
			if err != nil {
				return fmt.Errorf("saveGame: colonies: lastInsertId: %w", err)
			}
			colony.Id = int(id)

			_, err = tx.ExecContext(s.ctx, "insert into cors_dtl (cors_id, efftn, endtn, name, controlled_by) values (?, ?, ?, ?, ?)",
				colony.Id, colony.Details[0].EffTurn.String(), colony.Details[0].EndTurn.String(), colony.Details[0].Name, colony.Details[0].ControlledBy.Id)
			if err != nil {
				return fmt.Errorf("saveGame: colony: cors_dtl: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_population (cors_id, efftn, endtn, qty_professional, qty_soldier, qty_unskilled, qty_unemployed, qty_construction_crews, qty_spy_teams, rebel_pct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				colony.Id, colony.Population[0].EffTurn.String(), colony.Population[0].EndTurn.String(),
				colony.Population[0].QtyProfessional,
				colony.Population[0].QtySoldier,
				colony.Population[0].QtyUnskilled,
				colony.Population[0].QtyUnemployed,
				colony.Population[0].QtyConstructionCrew,
				colony.Population[0].QtySpyTeam,
				colony.Population[0].RebelPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: colony: cors_population: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_pay (cors_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				colony.Id, colony.Pay[0].EffTurn.String(), colony.Pay[0].EndTurn.String(),
				colony.Pay[0].ProfessionalPct,
				colony.Pay[0].SoldierPct,
				colony.Pay[0].UnskilledPct,
				colony.Pay[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: colony: cors_pay: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_rations (cors_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				colony.Id, colony.Rations[0].EffTurn.String(), colony.Rations[0].EndTurn.String(),
				colony.Rations[0].ProfessionalPct,
				colony.Rations[0].SoldierPct,
				colony.Rations[0].UnskilledPct,
				colony.Rations[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: colony: cors_rations: insert: %w", err)
			}

			for _, hull := range colony.Hull {
				if hull.Unit.Id == 0 {
					hull.Unit.Id = s.lookupUnitIdByCode(hull.Unit.Code)
				}
				_, err = tx.ExecContext(s.ctx, "insert into cors_hull (cors_id, unit_id, tech_level, efftn, endtn, qty_operational) values (?, ?, ?, ?, ?, ?)",
					colony.Id, hull.Unit.Id, hull.Unit.TechLevel, hull.EffTurn.String(), hull.EndTurn.String(), hull.QtyOperational)
				if err != nil {
					return fmt.Errorf("saveGame: colony: cors_hull: insert: %w", err)
				}
			}

			for _, inventory := range colony.Inventory {
				if inventory.Unit.Id == 0 {
					inventory.Unit.Id = s.lookupUnitIdByCode(inventory.Unit.Code)
				}
				_, err = tx.ExecContext(s.ctx, "insert into cors_inventory (cors_id, unit_id, tech_level, efftn, endtn, qty_operational, qty_stowed) values (?, ?, ?, ?, ?, ?, ?)",
					colony.Id, inventory.Unit.Id, inventory.Unit.TechLevel, inventory.EffTurn.String(), inventory.EndTurn.String(), inventory.QtyOperational, inventory.QtyStowed)
				if err != nil {
					return fmt.Errorf("saveGame: colony: cors_inventory: insert: %w", err)
				}
			}

			for _, group := range colony.Factories {
				if group.Unit.Id == 0 {
					group.Unit.Id = s.lookupUnitIdByCode(group.Unit.Code)
				}
				r, err := tx.ExecContext(s.ctx, "insert into cors_factory_group (cors_id, group_no, efftn, endtn, unit_id, tech_level) values (?, ?, ?, ?, ?, ?)",
					colony.Id, group.No, group.EffTurn.String(), group.EndTurn.String(), group.Unit.Id, group.Unit.TechLevel)
				if err != nil {
					return fmt.Errorf("saveGame: colony: cors_factory_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: colony: cors_factory_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					if unit.Unit.Id == 0 {
						unit.Unit.Id = s.lookupUnitIdByCode(unit.Unit.Code)
					}
					_, err = tx.ExecContext(s.ctx, "insert into cors_factory_group_units (factory_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn.String(), unit.EndTurn.String(), unit.Unit.Id, unit.Unit.TechLevel, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: colony: cors_factory_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into cors_factory_group_stages (factory_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
						group.Id, "0000/0", stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: colony: cors_factory_group_stages: insert: %w", err)
					}
				}
			}

			for _, group := range colony.Mines {
				r, err := tx.ExecContext(s.ctx, "insert into cors_mining_group (cors_id, group_no, efftn, endtn, resource_id) values (?, ?, ?, ?, ?)",
					colony.Id, group.No, group.EffTurn.String(), group.EndTurn.String(), group.Deposit.Id)
				if err != nil {
					return fmt.Errorf("saveGame: colony: cors_mining_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: colony: cors_mining_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					if unit.Unit.Id == 0 {
						unit.Unit.Id = s.lookupUnitIdByCode(unit.Unit.Code)
					}
					_, err = tx.ExecContext(s.ctx, "insert into cors_mining_group_units (mining_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn.String(), unit.EndTurn.String(), unit.Unit.Id, unit.Unit.TechLevel, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: colony: cors_mining_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into cors_mining_group_stages (mining_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
						group.Id, "0000/0", stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: colony: cors_mining_group_stages: insert: %w", err)
					}
				}
			}
			log.Printf("created nation %3d: colony %3d %8d\n", nation.No, colony.MSN, colony.Id)
		}

		for _, ship := range nation.Ships {
			if ship.Kind != "ship" {
				continue
			}
			r, err := tx.ExecContext(s.ctx, "insert into cors (game_id, msn, kind, planet_id) values (?, ?, ?, ?)",
				g.Id, ship.MSN, ship.Kind, ship.Details[0].Location.Id)
			if err != nil {
				return fmt.Errorf("saveGame: ships: insert: %w", err)
			}
			id, err := r.LastInsertId()
			if err != nil {
				return fmt.Errorf("saveGame: ships: lastInsertId: %w", err)
			}
			ship.Id = int(id)

			_, err = tx.ExecContext(s.ctx, "insert into cors_dtl (cors_id, efftn, endtn, name, controlled_by) values (?, ?, ?, ?, ?)",
				ship.Id, ship.Details[0].EffTurn.String(), ship.Details[0].EndTurn.String(), ship.Details[0].Name, ship.Details[0].ControlledBy.Id)
			if err != nil {
				return fmt.Errorf("saveGame: ship: cors_dtl: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_population (cors_id, efftn, endtn, qty_professional, qty_soldier, qty_unskilled, qty_unemployed, qty_construction_crews, qty_spy_teams, rebel_pct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				ship.Id, ship.Population[0].EffTurn.String(), ship.Population[0].EndTurn.String(),
				ship.Population[0].QtyProfessional,
				ship.Population[0].QtySoldier,
				ship.Population[0].QtyUnskilled,
				ship.Population[0].QtyUnemployed,
				ship.Population[0].QtyConstructionCrew,
				ship.Population[0].QtySpyTeam,
				ship.Population[0].RebelPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: ship: cors_dtl: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_pay (cors_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				ship.Id, ship.Pay[0].EffTurn.String(), ship.Pay[0].EndTurn.String(),
				ship.Pay[0].ProfessionalPct,
				ship.Pay[0].SoldierPct,
				ship.Pay[0].UnskilledPct,
				ship.Pay[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: ship: cors_pay: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_rations (cors_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				ship.Id, ship.Rations[0].EffTurn.String(), ship.Rations[0].EndTurn.String(),
				ship.Rations[0].ProfessionalPct,
				ship.Rations[0].SoldierPct,
				ship.Rations[0].UnskilledPct,
				ship.Rations[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: ship: cors_rations: insert: %w", err)
			}

			for _, hull := range ship.Hull {
				if hull.Unit.Id == 0 {
					hull.Unit.Id = s.lookupUnitIdByCode(hull.Unit.Code)
				}
				_, err = tx.ExecContext(s.ctx, "insert into cors_hull (cors_id, unit_id, tech_level, efftn, endtn, qty_operational) values (?, ?, ?, ?, ?, ?)",
					ship.Id, hull.Unit.Id, hull.Unit.TechLevel, hull.EffTurn.String(), hull.EndTurn.String(), hull.QtyOperational)
				if err != nil {
					return fmt.Errorf("saveGame: ship: cors_hull: insert: %w", err)
				}
			}

			for _, inventory := range ship.Inventory {
				if inventory.Unit.Id == 0 {
					inventory.Unit.Id = s.lookupUnitIdByCode(inventory.Unit.Code)
				}
				_, err = tx.ExecContext(s.ctx, "insert into cors_inventory (cors_id, unit_id, tech_level, efftn, endtn, qty_operational, qty_stowed) values (?, ?, ?, ?, ?, ?, ?)",
					ship.Id, inventory.Unit.Id, inventory.Unit.TechLevel, inventory.EffTurn.String(), inventory.EndTurn.String(), inventory.QtyOperational, inventory.QtyStowed)
				if err != nil {
					return fmt.Errorf("saveGame: ship: cors_inventory: insert: %w", err)
				}
			}

			for _, group := range ship.Factories {
				if group.Unit.Id == 0 {
					group.Unit.Id = s.lookupUnitIdByCode(group.Unit.Code)
				}
				r, err := tx.ExecContext(s.ctx, "insert into cors_factory_group (cors_id, group_no, efftn, endtn, unit_id, tech_level) values (?, ?, ?, ?, ?, ?)",
					ship.Id, group.No, group.EffTurn.String(), group.EndTurn.String(), group.Unit.Id, group.Unit.TechLevel)
				if err != nil {
					return fmt.Errorf("saveGame: ship: cors_factory_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: ship: cors_factory_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					if unit.Unit.Id == 0 {
						unit.Unit.Id = s.lookupUnitIdByCode(unit.Unit.Code)
					}
					_, err = tx.ExecContext(s.ctx, "insert into cors_factory_group_units (factory_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn.String(), unit.EndTurn.String(), unit.Unit.Id, unit.Unit.TechLevel, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: ship: cors_factory_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into cors_factory_group_stages (factory_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
						group.Id, "0000/0", stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: ship: cors_factory_group_stages: insert: %w", err)
					}
				}
			}
			log.Printf("created nation %3d: ship  %3d %8d\n", nation.No, ship.MSN, ship.Id)
		}
	}

	return tx.Commit()
}
