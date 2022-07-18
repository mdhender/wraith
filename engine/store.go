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

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/mdhender/wraith/models"
)

// Open returns an initialized engine
func Open(r *models.Store, options ...Option) (e *Engine, err error) {
	e = &Engine{r: r}
	for _, opt := range options {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

//func (e *Engine) fetchNations() ([]*Nation, error) {
//	var nations []*Nation
//	rows, err := e.db.Query(`select nations.id, nations.nation_no, nations.speciality, nations.descr,
//       nation_dtl.name, nation_dtl.govt_name, nation_dtl.govt_kind,
//       nation_skills.tech_level, nation_skills.research_points_pool
//		from nations, nation_dtl, nation_skills
//		where nations.game_id = ?
//		and (nation_dtl.nation_id = nations.id and nation_dtl.efftn <= ? and ? < nation_dtl.endtn)
//		and (nation_skills.nation_id = nations.id and nation_skills.efftn <= ? and ? < nation_skills.endtn)
//		order by nations.nation_no`,
//		e.game.Id,
//		e.game.CurrentTurn, e.game.CurrentTurn,
//		e.game.CurrentTurn, e.game.CurrentTurn)
//	if err != nil {
//		return nil, err
//	}
//	for rows.Next() {
//		nation := Nation{}
//		err := rows.Scan(&nation.Id, &nation.No, &nation.Speciality, &nation.Description,
//			&nation.Name, &nation.Government.Name, &nation.Government.Kind,
//			&nation.TechLevel, &nation.ResearchPool)
//		if err != nil {
//			return nil, err
//		}
//		nations = append(nations, &nation)
//	}
//	return nations, nil
//}
//
//func (e *Engine) fetchTurn(turn string) (*Turn, error) {
//	t := Turn{GameID: e.game.Id, Turn: turn}
//	row := e.db.QueryRow("select start_dt, end_dt from turns where game_id = ? and turn = ?", t.GameID, t.Turn)
//	err := row.Scan(&t.StartDt, &t.EndDt)
//	if err != nil {
//		return nil, nil
//	}
//	return &t, nil
//}
//
//func (e *Engine) LookupGame(id int) (*Game, error) {
//	var g Game
//	row := e.db.QueryRow("select id, short_name, name, descr, current_turn from games where id = ?", id)
//	err := row.Scan(&g.Id, &g.ShortName, &g.Name, &g.Descr, &g.CurrentTurn)
//	if err != nil {
//		return nil, err
//	}
//	return &g, nil
//}
//
//func (e *Engine) LookupGameByName(shortName string) *Game {
//	var g Game
//	row := e.db.QueryRow("select id, short_name, name, descr, current_turn from games where short_name= ?", shortName)
//	err := row.Scan(&g.Id, &g.ShortName, &g.Name, &g.Descr, &g.CurrentTurn)
//	if err != nil {
//		return nil
//	}
//	return &g
//}
//
//func (e *Engine) saveGame() error {
//	now := time.Now()
//
//	// get a transaction with a deferred rollback in case things fail
//	tx, err := e.db.BeginTx(e.ctx, nil)
//	if err != nil {
//		return fmt.Errorf("saveGame: beginTx: %w", err)
//	}
//	defer tx.Rollback()
//
//	r, err := tx.ExecContext(e.ctx, "insert into games (short_name, name, current_turn, descr) values (?, ?, ?, ?)",
//		e.game.ShortName, e.game.Name, e.game.CurrentTurn, e.game.Descr)
//	if err != nil {
//		return fmt.Errorf("saveGame: insert: 107: %w", err)
//	}
//	id, err := r.LastInsertId()
//	if err != nil {
//		return fmt.Errorf("saveGame: fetchId: %w", err)
//	}
//	e.game.Id = int(id)
//	log.Printf("created game %3d %s\n", int(id), e.game.ShortName)
//
//	for _, turn := range e.game.Turns {
//		_, err = tx.ExecContext(e.ctx, "insert into turns (game_id, turn, start_dt, end_dt) values (?, ?, ?, ?)",
//			e.game.Id, turn.Turn, turn.StartDt, turn.EndDt)
//		if err != nil {
//			return fmt.Errorf("saveGame: insert: 120: %w", err)
//		}
//	}
//
//	var nobody int
//	row := tx.QueryRow("select ifnull(user_id, 0) from user_profile where (effdt <= ? and ? < enddt) and handle = ?", now, now, "nobody")
//	err = row.Scan(&nobody)
//	if err != nil {
//		return fmt.Errorf("saveGame: players: nobody: %w", err)
//	}
//
//	for _, player := range e.game.Players {
//		var uid int
//		row := tx.QueryRow("select ifnull(user_id, 0) from user_profile where (effdt <= ? and ? < enddt) and handle = ?", now, now, player.Handle)
//		err = row.Scan(&uid)
//		if err != nil || uid == 0 {
//			uid = nobody
//			log.Printf("hey: player %q has no matching user\n", player.Handle)
//		}
//
//		r, err := tx.ExecContext(e.ctx, "insert into players (game_id, controlled_by, subject_of) values (?, ?, null)",
//			e.game.Id, uid)
//		if err != nil {
//			return fmt.Errorf("saveGame: players: insert: %w", err)
//		}
//		id, err := r.LastInsertId()
//		if err != nil {
//			return fmt.Errorf("saveGame: players: lastInsertId: %w", err)
//		}
//		player.Id = int(id)
//	}
//
//	for _, system := range e.game.Systems {
//		r, err := tx.ExecContext(e.ctx, "insert into systems (game_id, x, y, z, qty_stars) values (?, ?, ?, ?, ?)",
//			e.game.Id, system.X, system.Y, system.Z, len(system.Stars))
//		if err != nil {
//			return fmt.Errorf("saveGame: systems: insert: %w", err)
//		}
//		id, err := r.LastInsertId()
//		if err != nil {
//			return fmt.Errorf("saveGame: systems: lastInsertId: %w", err)
//		}
//		system.Id = int(id)
//
//		for _, star := range system.Stars {
//			r, err := tx.ExecContext(e.ctx, "insert into stars (system_id, sequence, kind) values (?, ?, ?)",
//				system.Id, star.Sequence, star.Kind)
//			if err != nil {
//				return fmt.Errorf("saveGame: stars: insert: %w", err)
//			}
//			id, err := r.LastInsertId()
//			if err != nil {
//				return fmt.Errorf("saveGame: stars: lastInsertId: %w", err)
//			}
//			star.Id = int(id)
//
//			for orbit, planet := range star.Orbits {
//				if orbit == 0 {
//					continue
//				}
//				homePlanet := "N"
//				if planet.HomePlanet {
//					homePlanet = "Y"
//				}
//				r, err := tx.ExecContext(e.ctx, "insert into planets (star_id, orbit_no, kind, habitability_no, home_planet) values (?, ?, ?, ?, ?)",
//					star.Id, planet.Orbit, planet.Kind, planet.HabitabilityNumber, homePlanet)
//				if err != nil {
//					return fmt.Errorf("saveGame: planet: insert: %w", err)
//				}
//				id, err := r.LastInsertId()
//				if err != nil {
//					return fmt.Errorf("saveGame: planet: lastInsertId: %w", err)
//				}
//				planet.Id = int(id)
//
//				for n, resource := range planet.Resources {
//					resource.No = n + 1
//					r, err := tx.ExecContext(e.ctx, "insert into resources (planet_id, deposit_no, kind, qty_initial, yield_pct) values (?, ?, ?, ?, ?)",
//						planet.Id, resource.No, resource.Kind, resource.InitialQuantity, resource.YieldPct)
//					if err != nil {
//						log.Printf("failed  system %8d: star %8d: orbit %2d: planet %8d: resource %8d %s\n", system.Id, star.Id, planet.Orbit, planet.Id, resource.Id, resource.Kind)
//						return fmt.Errorf("saveGame: resource: insert: %w", err)
//					}
//					id, err := r.LastInsertId()
//					if err != nil {
//						return fmt.Errorf("saveGame: resource: lastInsertId: %w", err)
//					}
//					resource.Id = int(id)
//					//log.Printf("created system %8d: star %8d: orbit %2d: planet %8d: resource %8d %-13s %9d\n", system.Id, star.Id, planet.Orbit, planet.Id, resource.Id, resource.Kind, resource.InitialQuantity)
//				}
//				//log.Printf("created system %8d: star %8d: orbit %2d: planet %8d\n", system.Id, star.Id, orbit, planet.Id)
//			}
//			//log.Printf("created system %8d: star %8d: suffix %q\n", system.Id, star.Id, star.Suffix)
//		}
//		//log.Printf("created system %8d\n", system.Id)
//	}
//
//	for _, nation := range e.game.Nations {
//		r, err := tx.ExecContext(e.ctx, "insert into nations (game_id, nation_no, speciality, descr) values (?, ?, ?, ?)",
//			e.game.Id, nation.No, nation.Speciality, nation.Description)
//		if err != nil {
//			return fmt.Errorf("saveGame: nations: insert: %w", err)
//		}
//		id, err := r.LastInsertId()
//		if err != nil {
//			return fmt.Errorf("saveGame: nations: lastInsertId: %w", err)
//		}
//		nation.Id = int(id)
//		_, err = tx.ExecContext(e.ctx, "insert into nation_dtl (nation_id, efftn, endtn, name, govt_name, govt_kind, controlled_by) values (?, ?, ?, ?, ?, ?, ?)",
//			nation.Id, "0000/0", "9999/9", nation.Name, nation.Government.Name, nation.Government.Kind, nation.ControlledBy.Id)
//		if err != nil {
//			return fmt.Errorf("saveGame: nation_dtl: insert: %w", err)
//		}
//		_, err = tx.ExecContext(e.ctx, "insert into nation_skills (nation_id, efftn, endtn, tech_level, research_points_pool, biology, bureaucracy, gravitics, life_support, manufacturing, military, mining, shields) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
//			nation.Id, "0000/0", "9999/9", nation.TechLevel, nation.ResearchPool, nation.Skills.Biology, nation.Skills.Bureaucracy, nation.Skills.Gravitics, nation.Skills.LifeSupport, nation.Skills.Manufacturing, nation.Skills.Military, nation.Skills.Mining, nation.Skills.Shields)
//		if err != nil {
//			return fmt.Errorf("saveGame: nation_skills: insert: %w", err)
//		}
//		log.Printf("created nation %3d %8d\n", nation.No, nation.Id)
//	}
//
//	for _, nation := range e.game.Nations {
//		for _, colony := range nation.Colonies {
//			r, err := tx.ExecContext(e.ctx, "insert into colonies (game_id, colony_no, planet_id, kind) values (?, ?, ?, ?)",
//				e.game.Id, colony.No, colony.Location.Id, colony.Kind)
//			if err != nil {
//				return fmt.Errorf("saveGame: colonies: insert: %w", err)
//			}
//			id, err := r.LastInsertId()
//			if err != nil {
//				return fmt.Errorf("saveGame: colonies: lastInsertId: %w", err)
//			}
//			colony.Id = int(id)
//
//			_, err = tx.ExecContext(e.ctx, "insert into colony_dtl (colony_id, efftn, endtn, name) values (?, ?, ?, ?)",
//				colony.Id, "0000/0", "9999/9", colony.Name)
//			if err != nil {
//				return fmt.Errorf("saveGame: colony_dtl: insert: %w", err)
//			}
//
//			_, err = tx.ExecContext(e.ctx, "insert into colony_population (colony_id, efftn, endtn, qty_professional, qty_soldier, qty_unskilled, qty_unemployed, qty_construction_crews, qty_spy_teams, rebel_pct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
//				colony.Id, "0000/0", "9999/9",
//				colony.Population.Professional.Qty,
//				colony.Population.Soldier.Qty,
//				colony.Population.Unskilled.Qty,
//				colony.Population.Unemployed.Qty,
//				colony.Population.ConstructionCrews,
//				colony.Population.SpyTeams,
//				colony.Population.RebelPct,
//			)
//			if err != nil {
//				return fmt.Errorf("saveGame: colony_dtl: insert: %w", err)
//			}
//
//			_, err = tx.ExecContext(e.ctx, "insert into colony_rations (colony_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
//				colony.Id, "0000/0", "9999/9",
//				colony.Population.Professional.Ration,
//				colony.Population.Soldier.Ration,
//				colony.Population.Unskilled.Ration,
//				colony.Population.Unemployed.Ration,
//			)
//			if err != nil {
//				return fmt.Errorf("saveGame: colony_rations: insert: %w", err)
//			}
//
//			_, err = tx.ExecContext(e.ctx, "insert into colony_pay (colony_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
//				colony.Id, "0000/0", "9999/9",
//				colony.Population.Professional.Pay,
//				colony.Population.Soldier.Pay,
//				colony.Population.Unskilled.Pay,
//				colony.Population.Unemployed.Pay,
//			)
//			if err != nil {
//				return fmt.Errorf("saveGame: colony_pay: insert: %w", err)
//			}
//
//			for _, inventory := range colony.Hull {
//				_, err = tx.ExecContext(e.ctx, "insert into colony_hull (colony_id, unit_id, tech_level, efftn, endtn, qty_operational) values (?, ?, ?, ?, ?, ?)",
//					colony.Id, inventory.Code, inventory.TechLevel, "0000/0", "9999/9", inventory.OperationalQty)
//				if err != nil {
//					return fmt.Errorf("saveGame: colony_hull: insert: %w", err)
//				}
//			}
//
//			for _, inventory := range colony.Inventory {
//				_, err = tx.ExecContext(e.ctx, "insert into colony_inventory (colony_id, unit_id, tech_level, efftn, endtn, qty_operational, qty_stowed) values (?, ?, ?, ?, ?, ?, ?)",
//					colony.Id, inventory.Code, inventory.TechLevel, "0000/0", "9999/9", inventory.OperationalQty, inventory.StowedQty)
//				if err != nil {
//					return fmt.Errorf("saveGame: colony_inventory: insert: %w", err)
//				}
//			}
//
//			for _, group := range colony.FactoryGroups {
//				r, err := tx.ExecContext(e.ctx, "insert into colony_factory_group (colony_id, group_no, efftn, endtn, unit_id, tech_level) values (?, ?, ?, ?, ?, ?)",
//					colony.Id, group.No, "0000/0", "9999/9", group.BuildCode, group.BuildTechLevel)
//				if err != nil {
//					return fmt.Errorf("saveGame: colony_factory_group: insert: %w", err)
//				}
//				id, err := r.LastInsertId()
//				if err != nil {
//					return fmt.Errorf("saveGame: colony_factory_group: lastInsertId: %w", err)
//				}
//				group.Id = int(id)
//				for _, unit := range group.Units {
//					_, err = tx.ExecContext(e.ctx, "insert into colony_factory_group_units (factory_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
//						group.Id, "0000/0", "9999/9", "FCT", unit.TechLevel, unit.Qty)
//					if err != nil {
//						return fmt.Errorf("saveGame: colony_factory_group_Units: insert: %w", err)
//					}
//					_, err = tx.ExecContext(e.ctx, "insert into colony_factory_group_stages (factory_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
//						group.Id, "0000/0", unit.Stages[0], unit.Stages[1], unit.Stages[2], 0)
//					if err != nil {
//						return fmt.Errorf("saveGame: colony_factory_group_Units: insert: %w", err)
//					}
//				}
//			}
//
//			for _, group := range colony.MiningGroups {
//				r, err := tx.ExecContext(e.ctx, "insert into colony_mining_group (colony_id, group_no, efftn, endtn, resource_id) values (?, ?, ?, ?, ?)",
//					colony.Id, group.No, "0000/0", "9999/9", group.Deposit.Id)
//				if err != nil {
//					return fmt.Errorf("saveGame: colony_mining_group: insert: %w", err)
//				}
//				id, err := r.LastInsertId()
//				if err != nil {
//					return fmt.Errorf("saveGame: colony_mining_group: lastInsertId: %w", err)
//				}
//				group.Id = int(id)
//				for _, unit := range group.Units {
//					_, err = tx.ExecContext(e.ctx, "insert into colony_mining_group_units (mining_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
//						group.Id, "0000/0", "9999/9", "MIN", unit.TechLevel, unit.Qty)
//					if err != nil {
//						return fmt.Errorf("saveGame: colony_mining_group_units: insert: %w", err)
//					}
//					_, err = tx.ExecContext(e.ctx, "insert into colony_mining_group_stages (mining_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
//						group.Id, "0000/0", unit.Stages[0], unit.Stages[1], unit.Stages[2], 0)
//					if err != nil {
//						return fmt.Errorf("saveGame: colony_mining_group_units: insert: %w", err)
//					}
//				}
//			}
//			log.Printf("created nation %3d: colony %3d %8d\n", nation.No, colony.No, colony.Id)
//		}
//
//		for _, ship := range nation.Ships {
//			r, err := tx.ExecContext(e.ctx, "insert into ships (game_id, ship_no, planet_id, kind) values (?, ?, ?, ?)",
//				e.game.Id, ship.No, ship.Location.Id, ship.Kind)
//			if err != nil {
//				return fmt.Errorf("saveGame: ships: insert: %w", err)
//			}
//			id, err := r.LastInsertId()
//			if err != nil {
//				return fmt.Errorf("saveGame: ships: lastInsertId: %w", err)
//			}
//			ship.Id = int(id)
//
//			_, err = tx.ExecContext(e.ctx, "insert into ship_dtl (ship_id, efftn, endtn, name) values (?, ?, ?, ?)",
//				ship.Id, "0000/0", "9999/9", ship.Name)
//			if err != nil {
//				return fmt.Errorf("saveGame: ship_dtl: insert: %w", err)
//			}
//
//			_, err = tx.ExecContext(e.ctx, "insert into ship_population (ship_id, efftn, endtn, qty_professional, qty_soldier, qty_unskilled, qty_unemployed, qty_construction_crews, qty_spy_teams, rebel_pct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
//				ship.Id, "0000/0", "9999/9",
//				ship.Population.Professional.Qty,
//				ship.Population.Soldier.Qty,
//				ship.Population.Unskilled.Qty,
//				ship.Population.Unemployed.Qty,
//				ship.Population.ConstructionCrews,
//				ship.Population.SpyTeams,
//				ship.Population.RebelPct,
//			)
//			if err != nil {
//				return fmt.Errorf("saveGame: ship_dtl: insert: %w", err)
//			}
//
//			_, err = tx.ExecContext(e.ctx, "insert into ship_rations (ship_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
//				ship.Id, "0000/0", "9999/9",
//				ship.Population.Professional.Ration,
//				ship.Population.Soldier.Ration,
//				ship.Population.Unskilled.Ration,
//				ship.Population.Unemployed.Ration,
//			)
//			if err != nil {
//				return fmt.Errorf("saveGame: ship_rations: insert: %w", err)
//			}
//
//			_, err = tx.ExecContext(e.ctx, "insert into ship_pay (ship_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
//				ship.Id, "0000/0", "9999/9",
//				ship.Population.Professional.Pay,
//				ship.Population.Soldier.Pay,
//				ship.Population.Unskilled.Pay,
//				ship.Population.Unemployed.Pay,
//			)
//			if err != nil {
//				return fmt.Errorf("saveGame: ship_pay: insert: %w", err)
//			}
//
//			for _, inventory := range ship.Hull {
//				_, err = tx.ExecContext(e.ctx, "insert into ship_hull (ship_id, unit_id, tech_level, efftn, endtn, qty_operational) values (?, ?, ?, ?, ?, ?)",
//					ship.Id, inventory.Code, inventory.TechLevel, "0000/0", "9999/9", inventory.OperationalQty)
//				if err != nil {
//					return fmt.Errorf("saveGame: ship_hull: insert: %w", err)
//				}
//			}
//
//			for _, inventory := range ship.Inventory {
//				_, err = tx.ExecContext(e.ctx, "insert into ship_inventory (ship_id, unit_id, tech_level, efftn, endtn, qty_operational, qty_stowed) values (?, ?, ?, ?, ?, ?, ?)",
//					ship.Id, inventory.Code, inventory.TechLevel, "0000/0", "9999/9", inventory.OperationalQty, inventory.StowedQty)
//				if err != nil {
//					return fmt.Errorf("saveGame: ship_inventory: insert: %w", err)
//				}
//			}
//
//			for _, group := range ship.FactoryGroups {
//				r, err := tx.ExecContext(e.ctx, "insert into ship_factory_group (ship_id, group_no, efftn, endtn, unit_id, tech_level) values (?, ?, ?, ?, ?, ?)",
//					ship.Id, group.No, "0000/0", "9999/9", group.BuildCode, group.BuildTechLevel)
//				if err != nil {
//					return fmt.Errorf("saveGame: ship_factory_group: insert: %w", err)
//				}
//				id, err := r.LastInsertId()
//				if err != nil {
//					return fmt.Errorf("saveGame: ship_factory_group: lastInsertId: %w", err)
//				}
//				group.Id = int(id)
//				for _, unit := range group.Units {
//					_, err = tx.ExecContext(e.ctx, "insert into ship_factory_group_units (factory_group_id, efftn, endtn, unit_id, tech_level, qty_operational) values (?, ?, ?, ?, ?, ?)",
//						group.Id, "0000/0", "9999/9", "FCT", unit.TechLevel, unit.Qty)
//					if err != nil {
//						return fmt.Errorf("saveGame: ship_factory_group_units: insert: %w", err)
//					}
//					_, err = tx.ExecContext(e.ctx, "insert into ship_factory_group_stages (factory_group_id, turn, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?)",
//						group.Id, "0000/0", unit.Stages[0], unit.Stages[1], unit.Stages[2], 0)
//					if err != nil {
//						return fmt.Errorf("saveGame: ship_factory_group_units: insert: %w", err)
//					}
//				}
//			}
//			log.Printf("created nation %3d: ship  %3d %8d\n", nation.No, ship.No, ship.Id)
//		}
//	}
//
//	return tx.Commit()
//}
//
//// Load retrieves a game from the store
//func (e *Engine) Load(id string) error {
//	if e == nil {
//		return ErrNoEngine
//	} else if e.db == nil {
//		return ErrNoStore
//	}
//	e.reset()
//
//	log.Printf("loading %q\n", id)
//	e.game = e.LookupGameByName(id)
//	if e.game == nil {
//		return ErrNoGame
//	}
//
//	log.Printf("loading %q: turn %q\n", e.game.ShortName, e.game.CurrentTurn)
//	turn, err := e.fetchTurn(e.game.CurrentTurn)
//	if err != nil {
//		return err
//	} else if turn == nil {
//		return ErrNoTurn
//	}
//	log.Printf("turn %v\n", *turn)
//
//	log.Printf("loading %q: nations\n", e.game.ShortName)
//	e.game.Nations, err = e.fetchNations()
//	if err != nil {
//		return err
//	} else if e.game.Nations == nil {
//		return ErrNoNation
//	}
//	for _, nation := range e.game.Nations {
//		log.Printf("game %q: nation %d %q\n", e.game.ShortName, nation.Id, nation.Name)
//	}
//
//	return nil
//}
//
//func (e *Engine) Save() error {
//	if e == nil {
//		return ErrNoEngine
//	} else if e.db == nil {
//		return ErrNoStore
//	} else if e.game == nil {
//		return ErrNoGame
//	}
//
//	return fmt.Errorf("engine.Save: %w", ErrNotImplemented)
//}
