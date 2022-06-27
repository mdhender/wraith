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
	"github.com/pkg/errors"
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
		return errors.Wrap(ErrMissingField, "short name")
	}
	for _, r := range g.ShortName { // check for invalid runes in the field
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return errors.Wrap(ErrInvalidField, "short name: invalid rune")
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

// FetchGameByName does just that
func (s *Store) FetchGameByName(name string) (*Game, error) {
	panic("!")
}

func (s *Store) GenerateGame(shortName, name, descr string, radius int, startDt time.Time, positions []*PlayerPosition) (*Game, error) {
	return s.genGame(shortName, name, descr, radius, startDt, positions)
}

func (s *Store) SaveGame(game *Game) error {
	panic("!")
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
	var g Game
	row := s.db.QueryRow("select id, short_name, name, descr, current_turn from games where id = ?", id)
	err := row.Scan(&g.Id, &g.ShortName, &g.Name, &g.Description, &g.CurrentTurn)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *Store) fetchGameByName(name string) (*Game, error) {
	var id int
	row := s.db.QueryRow("select id from games where short_name = ?", strings.ToUpper(name))
	err := row.Scan(&id)
	if err != nil {
		return nil, err
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
	for year := 0; year < 10; year++ {
		for quarter := 0; quarter < 4; quarter++ {
			if quarter == 0 && year != 0 {
				continue
			}
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

func (s *Store) saveGame(g *Game) error {
	// get a transaction with a deferred rollback in case things fail
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("saveGame: beginTx: %w", err)
	}
	defer tx.Rollback()

	r, err := tx.ExecContext(s.ctx, "insert into games (short_name, name, current_turn, descr) values (?, ?, ?, ?)",
		g.ShortName, g.Name, g.CurrentTurn, g.Description)
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
		_, err = tx.ExecContext(s.ctx, "insert into turns (game_id, turn, start_dt, end_dt) values (?, ?, ?, ?)",
			g.Id, turn.String(), turn.StartDt, turn.EndDt)
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
		row := tx.QueryRow("select ifnull(id, 0) from users where handle = ?", player.Details[0].Handle)
		err = row.Scan(&uid)
		if err != nil || uid == 0 {
			uid = nobody
			log.Printf("saveGame: player %q: no matching user\n", player.Details[0].Handle)
		}

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
			nation.Id, nation.Details[0].EffTurn, nation.Details[0].EndTurn, nation.Details[0].Name, nation.Details[0].GovtName, nation.Details[0].GovtKind, nation.Details[0].ControlledBy.Id)
		if err != nil {
			return fmt.Errorf("saveGame: nation_dtl: insert: %w", err)
		}
		_, err = tx.ExecContext(s.ctx, "insert into nation_skills (nation_id, efftn, endtn, tech_level, research_points_pool, biology, bureaucracy, gravitics, life_support, manufacturing, military, mining, shields) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			nation.Id, nation.Research[0].EffTurn, nation.Research[0].EndTurn, nation.Research[0].TechLevel, nation.Research[0].ResearchPointsPool, nation.Skills[0].Biology, nation.Skills[0].Bureaucracy, nation.Skills[0].Gravitics, nation.Skills[0].LifeSupport, nation.Skills[0].Manufacturing, nation.Skills[0].Military, nation.Skills[0].Mining, nation.Skills[0].Shields)
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
			r, err := tx.ExecContext(s.ctx, "insert into colonies (game_id, colony_no, planet_id, kind) values (?, ?, ?, ?)",
				g.Id, colony.MSN, colony.Details[0].Location.Id, colony.Kind)
			if err != nil {
				return fmt.Errorf("saveGame: colonies: insert: %w", err)
			}
			id, err := r.LastInsertId()
			if err != nil {
				return fmt.Errorf("saveGame: colonies: lastInsertId: %w", err)
			}
			colony.Id = int(id)

			_, err = tx.ExecContext(s.ctx, "insert into colony_dtl (colony_id, efftn, endtn, name) values (?, ?, ?, ?)",
				colony.Id, colony.Details[0].EffTurn, colony.Details[0].EndTurn, colony.Details[0].Name)
			if err != nil {
				return fmt.Errorf("saveGame: colony_dtl: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into colony_population (colony_id, efftn, endtn, qty_professional, qty_soldier, qty_unskilled, qty_unemployed, qty_construction_crews, qty_spy_teams, rebel_pct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				colony.Id, colony.Population[0].EffTurn, colony.Population[0].EndTurn,
				colony.Population[0].QtyProfessional,
				colony.Population[0].QtySoldier,
				colony.Population[0].QtyUnskilled,
				colony.Population[0].QtyUnemployed,
				colony.Population[0].QtyConstructionCrew,
				colony.Population[0].QtySpyTeam,
				colony.Population[0].RebelPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: colony_population: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into colony_pay (colony_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				colony.Id, colony.Pay[0].EffTurn, colony.Pay[0].EndTurn,
				colony.Pay[0].ProfessionalPct,
				colony.Pay[0].SoldierPct,
				colony.Pay[0].UnskilledPct,
				colony.Pay[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: colony_pay: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into colony_rations (colony_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				colony.Id, colony.Rations[0].EffTurn, colony.Rations[0].EndTurn,
				colony.Rations[0].ProfessionalPct,
				colony.Rations[0].SoldierPct,
				colony.Rations[0].UnskilledPct,
				colony.Rations[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: colony_rations: insert: %w", err)
			}

			for _, hull := range colony.Hull {
				_, err = tx.ExecContext(s.ctx, "insert into colony_hull (colony_id, unit_id, tech_level, efftn, endtn, qty_operational) values (?, ?, ?, ?, ?, ?)",
					colony.Id, hull.Unit.Code, hull.Unit.TechLevel, hull.EffTurn, hull.EndTurn, hull.QtyOperational)
				if err != nil {
					return fmt.Errorf("saveGame: colony_hull: insert: %w", err)
				}
			}

			for _, inventory := range colony.Inventory {
				_, err = tx.ExecContext(s.ctx, "insert into colony_inventory (colony_id, unit_id, tech_level, efftn, endtn, qty_operational, qty_stowed) values (?, ?, ?, ?, ?, ?, ?)",
					colony.Id, inventory.Unit.Code, inventory.Unit.TechLevel, inventory.EffTurn, inventory.EndTurn, inventory.QtyOperational, inventory.QtyStowed)
				if err != nil {
					return fmt.Errorf("saveGame: colony_inventory: insert: %w", err)
				}
			}

			for _, group := range colony.Factories {
				r, err := tx.ExecContext(s.ctx, "insert into colony_factory_group (colony_id, group_no, efftn, endtn, unit_id, tech_level) values (?, ?, ?, ?, ?, ?)",
					colony.Id, group.No, group.EffTurn, group.EndTurn, group.Unit.Code, group.Unit.TechLevel)
				if err != nil {
					return fmt.Errorf("saveGame: colony_factory_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: colony_factory_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					_, err = tx.ExecContext(s.ctx, "insert into colony_factory_group_units (factory_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn, unit.EndTurn, unit.Unit.Code, unit.Unit.TechLevel, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: colony_factory_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into colony_factory_group_stages (factory_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
						group.Id, "0000/0", stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: colony_factory_group_stages: insert: %w", err)
					}
				}
			}

			for _, group := range colony.Mines {
				r, err := tx.ExecContext(s.ctx, "insert into colony_mining_group (colony_id, group_no, efftn, endtn, resource_id) values (?, ?, ?, ?, ?)",
					colony.Id, group.No, group.EffTurn, group.EndTurn, group.Deposit.Id)
				if err != nil {
					return fmt.Errorf("saveGame: colony_mining_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: colony_mining_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					_, err = tx.ExecContext(s.ctx, "insert into colony_mining_group_units (mining_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn, unit.EndTurn, unit.Unit.Code, unit.Unit.TechLevel, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: colony_mining_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into colony_mining_group_stages (mining_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
						group.Id, "0000/0", stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: colony_mining_group_stages: insert: %w", err)
					}
				}
			}
			log.Printf("created nation %3d: colony %3d %8d\n", nation.No, colony.MSN, colony.Id)
		}

		for _, ship := range nation.Ships {
			if ship.Kind != "ship" {
				continue
			}
			r, err := tx.ExecContext(s.ctx, "insert into ships (game_id, ship_no, planet_id, kind) values (?, ?, ?, ?)",
				g.Id, ship.MSN, ship.Details[0].Location.Id, ship.Kind)
			if err != nil {
				return fmt.Errorf("saveGame: ships: insert: %w", err)
			}
			id, err := r.LastInsertId()
			if err != nil {
				return fmt.Errorf("saveGame: ships: lastInsertId: %w", err)
			}
			ship.Id = int(id)

			_, err = tx.ExecContext(s.ctx, "insert into ship_dtl (ship_id, efftn, endtn, name) values (?, ?, ?, ?)",
				ship.Id, ship.Details[0].EffTurn, ship.Details[0].EndTurn, ship.Details[0].Name)
			if err != nil {
				return fmt.Errorf("saveGame: ship_dtl: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into ship_population (ship_id, efftn, endtn, qty_professional, qty_soldier, qty_unskilled, qty_unemployed, qty_construction_crews, qty_spy_teams, rebel_pct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				ship.Id, ship.Population[0].EffTurn, ship.Population[0].EndTurn,
				ship.Population[0].QtyProfessional,
				ship.Population[0].QtySoldier,
				ship.Population[0].QtyUnskilled,
				ship.Population[0].QtyUnemployed,
				ship.Population[0].QtyConstructionCrew,
				ship.Population[0].QtySpyTeam,
				ship.Population[0].RebelPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: ship_dtl: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into ship_pay (ship_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				ship.Id, ship.Pay[0].EffTurn, ship.Pay[0].EndTurn,
				ship.Pay[0].ProfessionalPct,
				ship.Pay[0].SoldierPct,
				ship.Pay[0].UnskilledPct,
				ship.Pay[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: ship_pay: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into ship_rations (ship_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				ship.Id, ship.Rations[0].EffTurn, ship.Rations[0].EndTurn,
				ship.Rations[0].ProfessionalPct,
				ship.Rations[0].SoldierPct,
				ship.Rations[0].UnskilledPct,
				ship.Rations[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: ship_rations: insert: %w", err)
			}

			for _, hull := range ship.Hull {
				_, err = tx.ExecContext(s.ctx, "insert into ship_hull (ship_id, unit_id, tech_level, efftn, endtn, qty_operational) values (?, ?, ?, ?, ?, ?)",
					ship.Id, hull.Unit.Code, hull.Unit.TechLevel, hull.EffTurn, hull.EndTurn, hull.QtyOperational)
				if err != nil {
					return fmt.Errorf("saveGame: ship_hull: insert: %w", err)
				}
			}

			for _, inventory := range ship.Inventory {
				_, err = tx.ExecContext(s.ctx, "insert into ship_inventory (ship_id, unit_id, tech_level, efftn, endtn, qty_operational, qty_stowed) values (?, ?, ?, ?, ?, ?, ?)",
					ship.Id, inventory.Unit.Code, inventory.Unit.TechLevel, inventory.EffTurn, inventory.EndTurn, inventory.QtyOperational, inventory.QtyStowed)
				if err != nil {
					return fmt.Errorf("saveGame: ship_inventory: insert: %w", err)
				}
			}

			for _, group := range ship.Factories {
				r, err := tx.ExecContext(s.ctx, "insert into ship_factory_group (ship_id, group_no, efftn, endtn, unit_id, tech_level) values (?, ?, ?, ?, ?, ?)",
					ship.Id, group.No, group.EffTurn, group.EndTurn, group.Unit.Code, group.Unit.TechLevel)
				if err != nil {
					return fmt.Errorf("saveGame: ship_factory_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: ship_factory_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					_, err = tx.ExecContext(s.ctx, "insert into ship_factory_group_units (factory_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn, unit.EndTurn, unit.Unit.Code, unit.Unit.TechLevel, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: ship_factory_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into ship_factory_group_stages (factory_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
						group.Id, "0000/0", stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: ship_factory_group_stages: insert: %w", err)
					}
				}
			}
			log.Printf("created nation %3d: ship  %3d %8d\n", nation.No, ship.MSN, ship.Id)
		}
	}

	return tx.Commit()
}
