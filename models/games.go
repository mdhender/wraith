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

// FetchGameByNameAsOf fetches a game by name and turn
func (s *Store) FetchGameByNameAsOf(name string, asOfTurn string) (*Game, error) {
	game, err := s.lookupGameByName(name)
	if err != nil {
		return nil, err
	}
	return s.fetchGameByIdAsOf(game.Id, asOfTurn)
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

// fetch all cors
func (s *Store) fetchAllCorS(gameId int, turns map[string]*Turn, units map[int]*Unit) (map[int]*ColonyOrShip, error) {
	var corsId, controlledBy, unitId int
	var effTurn, endTurn string

	cors := make(map[int]*ColonyOrShip)

	rows, err := s.db.Query("select id, msn, kind from cors where game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: %w", gameId, err)
	}
	for rows.Next() {
		cs := &ColonyOrShip{}
		err := rows.Scan(&cs.Id, &cs.MSN, &cs.Kind)
		if err != nil {
			return nil, fmt.Errorf("fetchAllCorS: %d: %w", gameId, err)
		}
		cors[cs.Id] = cs
	}

	details := make(map[int][]*CSDetail)
	rows, err = s.db.Query("select cors_id, efftn, endtn, name, controlled_by from cors c, cors_dtl cd where c.game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: detail: %w", gameId, err)
	}
	for rows.Next() {
		row := &CSDetail{}
		err := rows.Scan(&corsId, &effTurn, &endTurn, &row.Name, &controlledBy)
		if err != nil {
			return nil, fmt.Errorf("fetchAllCorS: %d: detail: %w", gameId, err)
		}
		row.CS, row.EffTurn, row.EndTurn = cors[corsId], turns[effTurn], turns[endTurn]
		details[corsId] = append(details[corsId], row)
	}

	hulls := make(map[int][]*CSHull)
	rows, err = s.db.Query("select cors_id, efftn, endtn, unit_id, tech_level, qty_operational from cors c, cors_hull cd where c.game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: hull: %w", gameId, err)
	}
	for rows.Next() {
		row := &CSHull{}
		err := rows.Scan(&corsId, &effTurn, &endTurn, &unitId, &row.TechLevel, &row.QtyOperational)
		if err != nil {
			return nil, fmt.Errorf("fetchAllCorS: %d: hull: %w", gameId, err)
		}
		row.CS, row.EffTurn, row.EndTurn, row.Unit = cors[corsId], turns[effTurn], turns[endTurn], units[unitId]
		hulls[corsId] = append(hulls[corsId], row)
	}

	inventories := make(map[int][]*CSInventory)
	rows, err = s.db.Query("select cors_id, efftn, endtn, unit_id, tech_level, qty_operational, qty_stowed from cors c, cors_inventory cd where c.game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: inventory: %w", gameId, err)
	}
	for rows.Next() {
		row := &CSInventory{}
		err := rows.Scan(&corsId, &effTurn, &endTurn, &unitId, &row.TechLevel, &row.QtyOperational, &row.QtyStowed)
		if err != nil {
			return nil, fmt.Errorf("fetchAllCorS: %d: inventory: %w", gameId, err)
		}
		row.CS, row.EffTurn, row.EndTurn, row.Unit = cors[corsId], turns[effTurn], turns[endTurn], units[unitId]
		inventories[corsId] = append(inventories[corsId], row)
	}

	pays := make(map[int][]*CSPay)
	rows, err = s.db.Query("select cors_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct from cors c, cors_pay cd where c.game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: pay: %w", gameId, err)
	}
	for rows.Next() {
		row := &CSPay{}
		err := rows.Scan(&corsId, &effTurn, &endTurn, &row.ProfessionalPct, &row.SoldierPct, &row.UnskilledPct, &row.UnemployedPct)
		if err != nil {
			return nil, fmt.Errorf("fetchAllCorS: %d: pay: %w", gameId, err)
		}
		row.CS, row.EffTurn, row.EndTurn = cors[corsId], turns[effTurn], turns[endTurn]
		pays[corsId] = append(pays[corsId], row)
	}

	populations := make(map[int][]*CSPopulation)
	rows, err = s.db.Query("select cors_id, efftn, endtn, qty_professional, qty_soldier, qty_unskilled, qty_unemployed, qty_construction_crews, qty_spy_teams, rebel_pct from cors c, cors_population cd where c.game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: population: %w", gameId, err)
	}
	for rows.Next() {
		row := &CSPopulation{}
		err := rows.Scan(&corsId, &effTurn, &endTurn, &row.QtyProfessional, &row.QtySoldier, &row.QtyUnskilled, &row.QtyUnemployed, &row.QtyConstructionCrew, &row.QtySpyTeam, &row.RebelPct)
		if err != nil {
			return nil, fmt.Errorf("fetchAllCorS: %d: population: %w", gameId, err)
		}
		row.CS, row.EffTurn, row.EndTurn = cors[corsId], turns[effTurn], turns[endTurn]
		populations[corsId] = append(populations[corsId], row)
	}

	rations := make(map[int][]*CSRations)
	rows, err = s.db.Query("select cors_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct from cors c, cors_rations cd where c.game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: rations: %w", gameId, err)
	}
	for rows.Next() {
		row := &CSRations{}
		err := rows.Scan(&corsId, &effTurn, &endTurn, &row.ProfessionalPct, &row.SoldierPct, &row.UnskilledPct, &row.UnemployedPct)
		if err != nil {
			return nil, fmt.Errorf("fetchAllCorS: %d: rations: %w", gameId, err)
		}
		row.CS, row.EffTurn, row.EndTurn = cors[corsId], turns[effTurn], turns[endTurn]
		rations[corsId] = append(rations[corsId], row)
	}

	rows, err = s.db.Query("select cors_id, group_no, efftn, endtn from cors c, cors_factory_group cd where c.game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: factories: %w", gameId, err)
	}
	for rows.Next() {
	}

	log.Printf("fetchCorS: farm group not implemented!\n")
	//rows, err = s.db.Query("select cors_id, group_no, efftn, endtn from cors c, cors_farm_group cd where c.game_id = ?", gameId)
	//if err != nil {
	//	return nil, fmt.Errorf("fetchAllCorS: %d: farms: %w", gameId, err)
	//}
	//for rows.Next() {
	//}

	rows, err = s.db.Query("select cors_id, group_no, efftn, endtn from cors c, cors_mining_group cd where c.game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchAllCorS: %d: mines: %w", gameId, err)
	}
	for rows.Next() {
	}

	return cors, nil
}

// fetch game
func (s *Store) fetchGame(id int) (*Game, error) {
	now := time.Now()
	started := now

	//endOfTurns := &Turn{No: 9999*4 + 4, Year: 9999, Quarter: 4}

	// get a transaction with a deferred rollback in case things fail
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: beginTx: %w", err)
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	// fetch game
	g := &Game{}
	row := s.db.QueryRow("select id, short_name, name, descr, current_turn from games where id = ?", id)
	var currentTurn string
	err = row.Scan(&g.Id, &g.ShortName, &g.Name, &g.Description, &currentTurn)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: %w", id, err)
	} else if g.Id == 0 {
		return nil, fmt.Errorf("fetchGame: %d: %w", id, ErrNoDataFound)
	}
	log.Printf("fetchGame: %d: elapsed %v\n", id, time.Now().Sub(started))

	// fetch units
	units := make(map[int]*Unit)
	log.Printf("fetchGame: %d: todo: fix units\n", id)
	log.Printf("fetchGame: %d: elapsed %v\n", id, time.Now().Sub(started))

	// fetch users
	users, err := s.fetchUsers(tx, now)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: %w", id, err)
	}
	log.Printf("fetchGame: fetched %d users\n", len(users))
	log.Printf("fetchGame: %d: elapsed %v\n", id, time.Now().Sub(started))

	// fetch turns
	g.Turns, err = s.fetchTurns(g.Id)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: %w", id, err)
	}
	g.CurrentTurn = g.Turns[currentTurn]
	log.Printf("fetchGame: currentTurn %q\n", g.CurrentTurn.String())
	log.Printf("fetchGame: %d: elapsed %v\n", id, time.Now().Sub(started))

	// fetch systems
	g.Systems, err = s.fetchSystems(g.Id, g, g.Turns)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: %w", id, err)
	}
	log.Printf("fetchGame: fetched %d systems\n", len(g.Systems))
	log.Printf("fetchGame: %d: elapsed %v\n", id, time.Now().Sub(started))

	// fetch players
	g.Players, err = s.fetchPlayers(g.Id, g, users, g.Turns)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: %w", id, err)
	}
	for _, player := range g.Players {
		for _, detail := range player.Details {
			if detail.SubjectOf == nil {
				log.Printf("fetchGame: player %d %s controlled_by %q subject_of %q\n", player.Id, detail.Handle, detail.ControlledBy.Handle, "nobody")
			} else {
				log.Printf("fetchGame: player %d %s controlled_by %q subject_of %d\n", player.Id, detail.Handle, detail.ControlledBy.Handle, detail.SubjectOf.Id)
			}
		}
	}
	log.Printf("fetchGame: fetched %d players\n", len(g.Players))
	log.Printf("fetchGame: %d: elapsed %v\n", id, time.Now().Sub(started))

	// fetch nations
	g.Nations, err = s.fetchNations(g.Id, g, g.Players, g.Turns, units)
	if err != nil {
		return nil, fmt.Errorf("fetchGame: %d: nations: %w", id, err)
	}
	log.Printf("fetchGame: fetched %d nations\n", len(g.Nations))
	log.Printf("fetchGame: %d: elapsed %v\n", id, time.Now().Sub(started))

	return g, tx.Commit()
}

func (s *Store) fetchGameByName(name string) (*Game, error) {
	name = strings.ToUpper(strings.TrimSpace(name))
	var id int
	row := s.db.QueryRow("select id from games where short_name = ?", name)
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("fetchGameByName: %q: %w", name, err)
	} else if id == 0 {
		return nil, fmt.Errorf("fetchGame: %q: %w", name, ErrNoDataFound)
	}
	return s.fetchGame(id)
}

// fetch nation details
func (s *Store) fetchNationDetails(nationId int, nation *Nation, players map[int]*Player, turns map[string]*Turn) ([]*NationDetail, error) {
	details := []*NationDetail{}

	rows, err := s.db.Query("select efftn, endtn, name, govt_name, govt_kind, ifnull(controlled_by, 0) from nation_dtl where nation_id = ? order by efftn", nationId)
	if err != nil {
		return nil, fmt.Errorf("fetchNationDetails: %d: %w", nationId, err)
	}
	for rows.Next() {
		detail := &NationDetail{Nation: nation}
		var controlledBy int
		var effTurn, endTurn string
		err := rows.Scan(&effTurn, &endTurn, &detail.Name, &detail.GovtName, &detail.GovtKind, &controlledBy)
		if err != nil {
			return nil, fmt.Errorf("fetchNationDetails: %d: %w", nationId, err)
		}
		detail.EffTurn = turns[effTurn]
		detail.EndTurn = turns[endTurn]
		if controlledBy == 0 {
			return nil, fmt.Errorf("fetchNationDetails: %d: %w", controlledBy, ErrNoDataFound)
		}
		if controlledBy != 0 {
			player, ok := players[controlledBy]
			if !ok {
				return nil, fmt.Errorf("fetchNationDetails: %d: %w", controlledBy, ErrNoDataFound)
			}
			detail.ControlledBy = player
		}
		details = append(details, detail)
	}

	return details, nil
}

// fetch nation research
func (s *Store) fetchNationResearch(nationId int, nation *Nation, turns map[string]*Turn) ([]*NationResearch, error) {
	researchs := []*NationResearch{}

	rows, err := s.db.Query("select efftn, endtn, tech_level, research_points_pool  from nation_research where nation_id = ? order by efftn", nationId)
	if err != nil {
		return nil, fmt.Errorf("fetchNationResearch: %d: %w", nationId, err)
	}
	for rows.Next() {
		research := &NationResearch{Nation: nation}
		var effTurn, endTurn string
		err := rows.Scan(&effTurn, &endTurn, &research.TechLevel, &research.ResearchPointsPool)
		if err != nil {
			return nil, fmt.Errorf("fetchNationResearch: %d: %w", nationId, err)
		}
		research.EffTurn = turns[effTurn]
		research.EndTurn = turns[endTurn]
		researchs = append(researchs, research)
	}

	return researchs, nil
}

// fetch nation skills
func (s *Store) fetchNationSkills(nationId int, nation *Nation, turns map[string]*Turn) ([]*NationSkills, error) {
	skills := []*NationSkills{}

	rows, err := s.db.Query("select efftn, endtn, biology, bureaucracy, gravitics, life_support, manufacturing, military, mining, shields from nation_skills where nation_id = ? order by efftn", nationId)
	if err != nil {
		return nil, fmt.Errorf("fetchNationSkills: %d: %w", nationId, err)
	}
	for rows.Next() {
		skill := &NationSkills{Nation: nation}
		var effTurn, endTurn string
		err := rows.Scan(&effTurn, &endTurn, &skill.Biology, &skill.Bureaucracy, &skill.Gravitics, &skill.LifeSupport, &skill.Manufacturing, &skill.Military, &skill.Mining, &skill.Shields)
		if err != nil {
			return nil, fmt.Errorf("fetchNationSkills: %d: %w", nationId, err)
		}
		skill.EffTurn = turns[effTurn]
		skill.EndTurn = turns[endTurn]
		skills = append(skills, skill)
	}

	return skills, nil
}

// fetch nations
func (s *Store) fetchNations(gameId int, game *Game, players map[int]*Player, turns map[string]*Turn, units map[int]*Unit) (map[int]*Nation, error) {
	nations := make(map[int]*Nation)

	rows, err := s.db.Query("select id, nation_no, speciality, descr from nations where game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchNations: %d: %w", gameId, err)
	}
	for rows.Next() {
		nation := &Nation{Game: game}
		err := rows.Scan(&nation.Id, &nation.No, &nation.Speciality, &nation.Description)
		if err != nil {
			return nil, fmt.Errorf("fetchNations: %d: %w", gameId, err)
		}
		nation.Details, err = s.fetchNationDetails(nation.Id, nation, players, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchNations: %d: %w", gameId, err)
		}
		nation.Research, err = s.fetchNationResearch(nation.Id, nation, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchNations: %d: %w", gameId, err)
		}
		nation.Skills, err = s.fetchNationSkills(nation.Id, nation, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchNations: %d: %w", gameId, err)
		}
		nations[nation.Id] = nation
	}

	cs, err := s.fetchAllCorS(gameId, turns, units)
	if err != nil {
		return nil, fmt.Errorf("fetchNations: %d: %w", gameId, err)
	}
	log.Printf("fetchNations: fetched %d colonies/ships\n", len(cs))

	return nations, nil
}

// fetch nations as of turn
func (s *Store) fetchNationsAsOf(tx *sql.Tx, gameId int, asOfTurn string, game *Game) (map[int]*Nation, error) {
	nations := make(map[int]*Nation)

	rows, err := tx.Query("select id, nation_no, speciality, descr from nations where game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchNations: %d: %w", gameId, err)
	}
	for rows.Next() {
		nation := &Nation{Game: game}
		err := rows.Scan(&nation.Id, &nation.No, &nation.Speciality, &nation.Description)
		if err != nil {
			return nil, fmt.Errorf("fetchNationsAsOf: %d: %w", gameId, err)
		}
		nations[nation.Id] = nation
	}

	//for _, nation := range nations {
	//	nation.Details, err = s.fetchNationDetailsAsOf(tx, nation.Id, nation, players, turns)
	//	if err != nil {
	//		return nil, fmt.Errorf("fetchNationsAsOf: %d: %w", gameId, err)
	//	}
	//}
	//for _, nation := range nations {
	//	nation.Research, err = s.fetchNationResearchAsOf(tx, nation.Id, nation, turns)
	//	if err != nil {
	//		return nil, fmt.Errorf("fetchNationsAsOf: %d: %w", gameId, err)
	//	}
	//}
	//for _, nation := range nations {
	//	nation.Skills, err = s.fetchNationSkillsAsOf(tx, nation.Id, nation, turns)
	//	if err != nil {
	//		return nil, fmt.Errorf("fetchNationsAsOf: %d: %w", gameId, err)
	//	}
	//}
	//
	//cs, err := s.fetchCorSAsOf(tx, gameId, turns, units)
	//if err != nil {
	//	return nil, fmt.Errorf("fetchNationsAsOf: %d: %w", gameId, err)
	//}
	//log.Printf("fetchNationsAsOf: %d: fetched %d colonies/ships\n", gameId, len(cs))

	return nations, nil
}

// fetch natural resource details
func (s *Store) fetchNaturalResourceDetails(resourceId int, resource *NaturalResource, turns map[string]*Turn) ([]*NaturalResourceDetail, error) {
	details := []*NaturalResourceDetail{}

	rows, err := s.db.Query("select efftn, endtn, remaining_qty, ifnull(controlled_by, 0) from resource_dtl where resource_id = ? order by efftn", resourceId)
	if err != nil {
		return nil, fmt.Errorf("fetchNaturalResourceDetails: %d: %w", resourceId, err)
	}
	for rows.Next() {
		detail := &NaturalResourceDetail{NaturalResource: resource}
		var effTurn, endTurn string
		var controlledBy int
		err := rows.Scan(&effTurn, &endTurn, &detail.QtyRemaining, &controlledBy)
		if err != nil {
			return nil, fmt.Errorf("fetchNaturalResourceDetails: %d: %w", resourceId, err)
		}
		detail.EffTurn = turns[effTurn]
		detail.EndTurn = turns[endTurn]
		details = append(details, detail)
	}

	return details, nil
}

// fetch natural resources
func (s *Store) fetchNaturalResources(planetId int, planet *Planet, turns map[string]*Turn) ([]*NaturalResource, error) {
	resources := []*NaturalResource{}

	rows, err := s.db.Query("select id, deposit_no, unit_id, qty_initial, yield_pct from resources where planet_id = ? order by deposit_no", planetId)
	if err != nil {
		return nil, fmt.Errorf("fetchNaturalResources: %d: %w", planetId, err)
	}
	for rows.Next() {
		var unitId int
		resource := &NaturalResource{Planet: planet}
		err := rows.Scan(&resource.Id, &resource.No, &unitId, &resource.QtyInitial, &resource.YieldPct)
		if err != nil {
			return nil, fmt.Errorf("fetchNaturalResources: %d: %w", planetId, err)
		}
		resource.Unit = s.lookupUnit(unitId)
		resource.Details, err = s.fetchNaturalResourceDetails(resource.Id, resource, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchNaturalResources: %d: %w", planetId, err)
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// fetch planet details
func (s *Store) fetchPlanetDetails(planetId int, planet *Planet, turns map[string]*Turn) ([]*PlanetDetail, error) {
	details := []*PlanetDetail{}

	rows, err := s.db.Query("select efftn, endtn, ifnull(controlled_by, 0), habitability_no from planet_dtl where planet_id = ? order by efftn", planetId)
	if err != nil {
		return nil, fmt.Errorf("fetchPlanetDetails: %d: %w", planetId, err)
	}
	for rows.Next() {
		detail := &PlanetDetail{Planet: planet}
		var effTurn, endTurn string
		var controlledBy int
		err := rows.Scan(&effTurn, &endTurn, &controlledBy, &detail.HabitabilityNo)
		if err != nil {
			return nil, fmt.Errorf("fetchPlanetDetails: %d: %w", planetId, err)
		}
		detail.EffTurn = turns[effTurn]
		detail.EndTurn = turns[endTurn]
		details = append(details, detail)
	}

	return details, nil
}

// fetch planets
func (s *Store) fetchPlanets(starId int, star *Star, turns map[string]*Turn) ([]*Planet, error) {
	planets := []*Planet{}

	rows, err := s.db.Query("select id, orbit_no, kind, home_planet from planets where star_id = ?", starId)
	if err != nil {
		return nil, fmt.Errorf("fetchPlanets: %d: %w", starId, err)
	}
	for rows.Next() {
		planet := &Planet{Star: star}
		var homePlanet string
		err := rows.Scan(&planet.Id, &planet.OrbitNo, &planet.Kind, &homePlanet)
		if err != nil {
			return nil, fmt.Errorf("fetchPlanets: %d: %w", starId, err)
		}
		planet.HomePlanet = homePlanet == "Y"
		planet.Details, err = s.fetchPlanetDetails(planet.Id, planet, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchPlanets: %d: %w", starId, err)
		}
		planet.Deposits, err = s.fetchNaturalResources(planet.Id, planet, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchPlanets: %d: %w", starId, err)
		}
		//planet.Colonies, err = s.fetchPlanetColonies(planet.Id, planet)
		//if err != nil {
		//	return nil, fmt.Errorf("fetchPlanets: %d: %w", starId, err)
		//}
		//planet.Ships, err = s.fetchPlanetShips(planet.Id, planet)
		//if err != nil {
		//	return nil, fmt.Errorf("fetchPlanets: %d: %w", starId, err)
		//}
		planets = append(planets, planet)
	}

	return planets, nil
}

// fetch player details
func (s *Store) fetchPlayerDetails(playerId int, users map[int]*User, players map[int]*Player, turns map[string]*Turn) ([]*PlayerDetail, error) {
	player, ok := players[playerId]
	if !ok {
		return nil, fmt.Errorf("fetchPlayerDetails: %d: %w", playerId, ErrNoDataFound)
	}

	details := []*PlayerDetail{}

	rows, err := s.db.Query("select efftn, endtn, handle, ifnull(controlled_by, 0), ifnull(subject_of, 0) from player_dtl where player_id = ? order by efftn", playerId)
	if err != nil {
		return nil, fmt.Errorf("fetchPlayerDetails: %d: %w", playerId, err)
	}
	for rows.Next() {
		dtl := &PlayerDetail{Player: player}
		var effTurn, endTurn string
		var controlledBy, subjectOf int
		err := rows.Scan(&effTurn, &endTurn, &dtl.Handle, &controlledBy, &subjectOf)
		if err != nil {
			return nil, fmt.Errorf("fetchPlayerDetails: %d: %w", playerId, err)
		}
		dtl.EffTurn = turns[effTurn]
		dtl.EndTurn = turns[endTurn]
		if controlledBy == 0 {
			return nil, fmt.Errorf("fetchPlayerDetails: %d: %w", playerId, ErrNoDataFound)
		}
		if controlledBy != 0 {
			user, ok := users[controlledBy]
			if !ok {
				return nil, fmt.Errorf("fetchPlayerDetails: %d: %w", playerId, ErrNoDataFound)
			}
			dtl.ControlledBy = user
		}
		if subjectOf != 0 {
			reportsTo, ok := players[subjectOf]
			if !ok {
				return nil, fmt.Errorf("fetchPlayerDetails: %d: %w", playerId, ErrNoDataFound)
			}
			dtl.SubjectOf = reportsTo
		}
		details = append(details, dtl)
	}

	return details, nil
}

// fetch players
func (s *Store) fetchPlayers(gameId int, game *Game, users map[int]*User, turns map[string]*Turn) (map[int]*Player, error) {
	players := make(map[int]*Player)

	rows, err := s.db.Query("select id from players where game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchPlayers: %d: %w", gameId, err)
	}
	for rows.Next() {
		player := &Player{Game: game}
		err := rows.Scan(&player.Id)
		if err != nil {
			return nil, fmt.Errorf("fetchPlayers: %d: %w", gameId, err)
		}
		players[player.Id] = player
	}

	// fetch the details for each player
	for _, player := range players {
		player.Details, err = s.fetchPlayerDetails(player.Id, users, players, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchPlayers: %d: %w", gameId, err)
		}
	}

	return players, nil
}

// fetch stars
func (s *Store) fetchStars(systemId int, system *System, turns map[string]*Turn) ([]*Star, error) {
	stars := []*Star{}

	rows, err := s.db.Query("select id, sequence, kind from stars where system_id = ?", systemId)
	if err != nil {
		return nil, fmt.Errorf("fetchStars: %d: %w", systemId, err)
	}
	for rows.Next() {
		star := &Star{System: system}
		err := rows.Scan(&star.Id, &star.Sequence, &star.Kind)
		if err != nil {
			return nil, fmt.Errorf("fetchStars: %d: %w", systemId, err)
		}
		star.Orbits, err = s.fetchPlanets(star.Id, star, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchStars: %d: %w", systemId, err)
		}
		stars = append(stars, star)
	}

	return stars, nil
}

// fetch systems
func (s *Store) fetchSystems(gameId int, game *Game, turns map[string]*Turn) (map[int]*System, error) {
	systems := make(map[int]*System)

	rows, err := s.db.Query("select id, x, y, z from systems where game_id = ?", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchSystems: %d: %w", gameId, err)
	}
	for rows.Next() {
		system := &System{Game: game}
		err := rows.Scan(&system.Id, &system.Coords.X, &system.Coords.Y, &system.Coords.Z)
		if err != nil {
			return nil, fmt.Errorf("fetchSystems: %d: %w", gameId, err)
		}
		system.Stars, err = s.fetchStars(system.Id, system, turns)
		if err != nil {
			return nil, fmt.Errorf("fetchSystems: %d: %w", gameId, err)
		}
		systems[system.Id] = system
	}

	return systems, nil
}

// fetch turn
func (s *Store) fetchTurn(gameId int, turn string) (*Turn, error) {
	row := s.db.QueryRow("select no, year, quarter, start_dt, end_dt from turns where game_id = ? and turn = ?", gameId, turn)
	t := &Turn{}
	err := row.Scan(&t.No, &t.Year, &t.Quarter, &t.StartDt, &t.EndDt)
	if err != nil {
		return nil, fmt.Errorf("fetchTurn: %d %q: %w", gameId, turn, err)
	}
	return t, nil
}

// fetch turns
func (s *Store) fetchTurns(gameId int) (map[string]*Turn, error) {
	turns := make(map[string]*Turn)

	rows, err := s.db.Query("select no, year, quarter, start_dt, end_dt from turns where game_id = ? order by turn", gameId)
	if err != nil {
		return nil, fmt.Errorf("fetchTurns: %d: %w", gameId, err)
	}
	for rows.Next() {
		turn := &Turn{}
		err := rows.Scan(&turn.No, &turn.Year, &turn.Quarter, &turn.StartDt, &turn.EndDt)
		if err != nil {
			return nil, fmt.Errorf("fetchTurns: %d: %w", gameId, err)
		}
		turns[turn.String()] = turn
	}

	return turns, nil
}

// fetch units
func (s *Store) fetchUnits(tx *sql.Tx) (map[int]*Unit, error) {
	units := make(map[int]*Unit)

	rows, err := tx.Query("select id, code, name, descr from units")
	if err != nil {
		return nil, fmt.Errorf("fetchUnits: %w", err)
	}
	for rows.Next() {
		unit := &Unit{}
		err := rows.Scan(&unit.Id, &unit.Code, &unit.Name, &unit.Description)
		if err != nil {
			return nil, fmt.Errorf("fetchUnits: %w", err)
		}
		units[unit.Id] = unit
	}

	return units, nil
}

// fetch users
func (s *Store) fetchUsers(tx *sql.Tx, asOf time.Time) (map[int]*User, error) {
	users := make(map[int]*User)

	rows, err := tx.Query("select u.id, u.handle, up.effdt, up.enddt, up.email, up.handle from users u, user_profile up where up.user_id = u.id and (up.effdt <= ? and ? < up.enddt) order by id", asOf, asOf)
	if err != nil {
		return nil, fmt.Errorf("fetchUsers: %w", err)
	}
	for rows.Next() {
		user := &User{Profiles: []*UserProfile{{}}}
		err := rows.Scan(&user.Id, &user.Handle, &user.Profiles[0].EffDt, &user.Profiles[0].EndDt, &user.Profiles[0].Email, &user.Profiles[0].Handle)
		if err != nil {
			return nil, fmt.Errorf("fetchUsers: %w", err)
		}
		//log.Printf("fetchUsers: user %d %q profile %q %q %v %v\n",
		//user.Id, user.Handle,
		//user.Profiles[0].Handle, user.Profiles[0].Email, user.Profiles[0].EffDt, user.Profiles[0].EndDt)
		users[user.Id] = user
	}

	return users, nil
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
		Nations:     make(map[int]*Nation),
		Players:     make(map[int]*Player),
		Stars:       make(map[int]*Star),
		Systems:     make(map[int]*System),
		Turns:       make(map[string]*Turn),
	}

	// convert positions to players
	for no, position := range positions {
		user, err := s.fetchUserByHandle(position.UserHandle)
		if err != nil {
			return nil, fmt.Errorf("user: %q: %w", position.UserHandle, ErrNoDataFound)
		} else if user == nil {
			return nil, fmt.Errorf("user: %q: %w", position.UserHandle, ErrNoDataFound)
		}
		player := &Player{
			Id:       no + 1,
			Game:     game,
			MemberOf: nil,
			Details:  nil,
		}
		player.Details = []*PlayerDetail{{
			Player:       player,
			EffTurn:      effTurn,
			EndTurn:      endTurn,
			Handle:       position.PlayerHandle,
			ControlledBy: user,
			SubjectOf:    nil,
		}}
		game.Players[player.Id] = player
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
	turn := &Turn{No: turnNo, Year: 0, Quarter: 0, StartDt: effDt, EndDt: endDt}
	game.Turns[turn.String()] = turn
	effDt = endDt
	endDt = effDt.Add(turnDuration)
	turnNo++
	for year := 1; year <= 10; year++ {
		for quarter := 1; quarter <= 4; quarter++ {
			turn = &Turn{No: turnNo, Year: year, Quarter: quarter, StartDt: effDt, EndDt: endDt}
			game.Turns[turn.String()] = turn
			effDt = endDt
			endDt = effDt.Add(turnDuration)
			turnNo++
		}
	}

	systemId, ring, colonyNo := 0, 5, 0

	// generate nations and their home systems
	for no, position := range positions {
		// warning: assumes that player was created for this game
		player, ok := game.Players[no+1]
		if !ok {
			log.Printf("no %d: position id %d\n", no, position.Id)
			return nil, fmt.Errorf("genGame: player %d: %w", no+1, ErrNoDataFound)
		}

		systemId++
		coords := rings[ring][0]
		rings[ring] = rings[ring][1:]

		system := s.genHomeSystem(systemId)
		system.Ring, system.Coords = ring, coords
		game.Systems[system.Id] = system

		planet := system.Stars[0].Orbits[3]
		nation := s.genNation(no+1, planet, player, position)
		colonyNo++
		nation.Colonies[0].MSN = colonyNo
		colonyNo++
		nation.Colonies[1].MSN = colonyNo

		player.MemberOf = nation // link the player to its nation

		game.Nations[nation.No] = nation
		//log.Printf("genGame: nation no %d: controlledBy %v\n", nation.No, nation.Details[0].ControlledBy)
	}

	// generate the remainder of the systems
	for ring := 0; ring < len(rings); ring++ {
		for _, coords := range rings[ring] {
			systemId++
			system := s.genSystem(systemId)
			system.Ring, system.Coords = ring, coords
			game.Systems[system.Id] = system
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
	g.CurrentTurn, err = s.fetchTurn(g.Id, currentTurn)
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
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	if g.CurrentTurn == nil {
		g.CurrentTurn = g.Turns["0000/0"]
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
		row := tx.QueryRow("select ifnull(id, 0) from users where handle = ?", player.Details[0].ControlledBy.Handle)
		err = row.Scan(&uid)
		if err != nil {
			return fmt.Errorf("saveGame: users: %w", err)
		} else if uid == 0 {
			return fmt.Errorf("saveGame: users: %q: %w", player.Details[0].ControlledBy.Handle, ErrNoDataFound)
		}
		log.Printf("saveGame: mapped %8d to player %q\n", uid, player.Details[0].Handle)

		r, err := tx.ExecContext(s.ctx, "insert into players (game_id) values (?)", g.Id)
		if err != nil {
			return fmt.Errorf("saveGame: players: insert: %w", err)
		}
		id, err := r.LastInsertId()
		if err != nil {
			return fmt.Errorf("saveGame: players: lastInsertId: %w", err)
		}
		player.Id = int(id)
		for _, detail := range player.Details {
			if detail.SubjectOf == nil {
				_, err = tx.ExecContext(s.ctx, "insert into player_dtl (player_id, efftn, endtn, handle, controlled_by) values (?, ?, ?, ?, ?)",
					player.Id, detail.EffTurn.String(), detail.EndTurn.String(), detail.Handle, detail.ControlledBy.Id)
			} else {
				_, err = tx.ExecContext(s.ctx, "insert into player_dtl (player_id, efftn, endtn, handle, controlled_by, subject_of) values (?, ?, ?, ?, ?, ?)",
					player.Id, detail.EffTurn.String(), detail.EndTurn.String(), detail.Handle, detail.ControlledBy.Id, detail.SubjectOf.Id)
			}
			if err != nil {
				return fmt.Errorf("saveGame: player_dtl: insert: %w", err)
			}
		}
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
				r, err := tx.ExecContext(s.ctx, "insert into planets (star_id, orbit_no, kind, home_planet) values (?, ?, ?, ?)",
					star.Id, planet.OrbitNo, planet.Kind, homePlanet)
				if err != nil {
					return fmt.Errorf("saveGame: planet: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: planet: lastInsertId: %w", err)
				}
				planet.Id = int(id)

				for _, detail := range planet.Details {
					if detail.ControlledBy == nil {
						_, err = tx.ExecContext(s.ctx, "insert into planet_dtl (planet_id, efftn, endtn, habitability_no) values (?, ?, ?, ?)",
							planet.Id, detail.EffTurn.String(), detail.EndTurn.String(), detail.HabitabilityNo)
					} else {
						_, err = tx.ExecContext(s.ctx, "insert into planet_dtl (planet_id, efftn, endtn, controlled_by, habitability_no) values (?, ?, ?, ?, ?)",
							planet.Id, detail.EffTurn.String(), detail.EndTurn.String(), detail.ControlledBy.Id, detail.HabitabilityNo)
					}
					if err != nil {
						return fmt.Errorf("saveGame: planetDetail: insert: %w", err)
					}
				}

				for n, deposit := range planet.Deposits {
					deposit.No = n + 1
					r, err := tx.ExecContext(s.ctx, "insert into resources (planet_id, deposit_no, unit_id, qty_initial, yield_pct) values (?, ?, ?, ?, ?)",
						planet.Id, deposit.No, deposit.Unit.Id, deposit.QtyInitial, deposit.YieldPct)
					if err != nil {
						log.Printf("failed  system %8d: star %8d: orbit %2d: planet %8d: resource %8d %s\n", system.Id, star.Id, planet.OrbitNo, planet.Id, deposit.Id, deposit.Unit.Code)
						return fmt.Errorf("saveGame: deposit: insert: %w", err)
					}
					id, err := r.LastInsertId()
					if err != nil {
						return fmt.Errorf("saveGame: deposit: lastInsertId: %w", err)
					}
					deposit.Id = int(id)
					for _, detail := range deposit.Details {
						if detail.ControlledBy == nil {
							_, err = tx.ExecContext(s.ctx, "insert into resource_dtl (resource_id, efftn, endtn, remaining_qty) values (?, ?, ?, ?)",
								deposit.Id, detail.EffTurn.String(), detail.EndTurn.String(), detail.QtyRemaining)
						} else {
							_, err = tx.ExecContext(s.ctx, "insert into resource_dtl (resource_id, efftn, endtn, remaining_qty, controlled_by) values (?, ?, ?, ?, ?)",
								deposit.Id, detail.EffTurn.String(), detail.EndTurn.String(), detail.QtyRemaining, detail.ControlledBy)
						}
						if err != nil {
							log.Printf("failed  system %8d: star %8d: orbit %2d: planet %8d: resource %8d %s\n", system.Id, star.Id, planet.OrbitNo, planet.Id, deposit.Id, deposit.Unit.Code)
							return fmt.Errorf("saveGame: depositDetail: insert: %w", err)
						}
					}
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
		if nation.Details[0].ControlledBy == nil {
			return fmt.Errorf("saveGame: nation %d: nationDetails: controlledBy: nil: %w", nation.Id, ErrNoDataFound)
		}
		_, err = tx.ExecContext(s.ctx, "insert into nation_dtl (nation_id, efftn, endtn, name, govt_name, govt_kind, controlled_by) values (?, ?, ?, ?, ?, ?, ?)",
			nation.Id, nation.Details[0].EffTurn.String(), nation.Details[0].EndTurn.String(), nation.Details[0].Name, nation.Details[0].GovtName, nation.Details[0].GovtKind, nation.Details[0].ControlledBy.Id)
		if err != nil {
			return fmt.Errorf("saveGame: nation_dtl: insert: %w", err)
		}
		_, err = tx.ExecContext(s.ctx, "insert into nation_research (nation_id, efftn, endtn, tech_level, research_points_pool) values (?, ?, ?, ?, ?)",
			nation.Id, nation.Research[0].EffTurn.String(), nation.Research[0].EndTurn.String(), nation.Research[0].TechLevel, nation.Research[0].ResearchPointsPool)
		if err != nil {
			return fmt.Errorf("saveGame: nation_research: insert: %w", err)
		}
		_, err = tx.ExecContext(s.ctx, "insert into nation_skills (nation_id, efftn, endtn, biology, bureaucracy, gravitics, life_support, manufacturing, military, mining, shields) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			nation.Id, nation.Skills[0].EffTurn.String(), nation.Skills[0].EndTurn.String(), nation.Skills[0].Biology, nation.Skills[0].Bureaucracy, nation.Skills[0].Gravitics, nation.Skills[0].LifeSupport, nation.Skills[0].Manufacturing, nation.Skills[0].Military, nation.Skills[0].Mining, nation.Skills[0].Shields)
		if err != nil {
			return fmt.Errorf("saveGame: nation_skills: insert: %w", err)
		}
		log.Printf("created nation %3d %8d\n", nation.No, nation.Id)
	}

	for _, nation := range g.Nations {
		numColonies, numShips := 0, 0
		for _, cs := range nation.CorS {
			if cs.Kind == "ship" {
				numShips++
			} else {
				numColonies++
			}
			r, err := tx.ExecContext(s.ctx, "insert into cors (game_id, msn, kind) values (?, ?, ?)",
				g.Id, cs.MSN, cs.Kind)
			if err != nil {
				return fmt.Errorf("saveGame: cors: insert: %w", err)
			}
			id, err := r.LastInsertId()
			if err != nil {
				return fmt.Errorf("saveGame: cors: lastInsertId: %w", err)
			}
			cs.Id = int(id)

			_, err = tx.ExecContext(s.ctx, "insert into cors_dtl (cors_id, efftn, endtn, name, tech_level, controlled_by) values (?, ?, ?, ?, ?, ?)",
				cs.Id, cs.Details[0].EffTurn.String(), cs.Details[0].EndTurn.String(), cs.Details[0].Name, cs.Details[0].TechLevel, cs.Details[0].ControlledBy.Id)
			if err != nil {
				return fmt.Errorf("saveGame: cors_dtl: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_loc (cors_id, efftn, endtn, planet_id) values (?, ?, ?, ?)",
				cs.Id, cs.Locations[0].EffTurn.String(), cs.Locations[0].EndTurn.String(), cs.Locations[0].Location.Id)
			if err != nil {
				return fmt.Errorf("saveGame: cors_loc: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_population (cors_id, efftn, endtn, qty_professional, qty_soldier, qty_unskilled, qty_unemployed, qty_construction_crews, qty_spy_teams, rebel_pct) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				cs.Id, cs.Population[0].EffTurn.String(), cs.Population[0].EndTurn.String(),
				cs.Population[0].QtyProfessional,
				cs.Population[0].QtySoldier,
				cs.Population[0].QtyUnskilled,
				cs.Population[0].QtyUnemployed,
				cs.Population[0].QtyConstructionCrew,
				cs.Population[0].QtySpyTeam,
				cs.Population[0].RebelPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: cors_population: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_pay (cors_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				cs.Id, cs.Pay[0].EffTurn.String(), cs.Pay[0].EndTurn.String(),
				cs.Pay[0].ProfessionalPct,
				cs.Pay[0].SoldierPct,
				cs.Pay[0].UnskilledPct,
				cs.Pay[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: cors_pay: insert: %w", err)
			}

			_, err = tx.ExecContext(s.ctx, "insert into cors_rations (cors_id, efftn, endtn, professional_pct, soldier_pct, unskilled_pct, unemployed_pct) values (?, ?, ?, ?, ?, ?, ?)",
				cs.Id, cs.Rations[0].EffTurn.String(), cs.Rations[0].EndTurn.String(),
				cs.Rations[0].ProfessionalPct,
				cs.Rations[0].SoldierPct,
				cs.Rations[0].UnskilledPct,
				cs.Rations[0].UnemployedPct,
			)
			if err != nil {
				return fmt.Errorf("saveGame: cors_rations: insert: %w", err)
			}

			for _, hull := range cs.Hull {
				if hull.Unit.Id == 0 {
					hull.Unit.Id = s.lookupUnitIdByCode(hull.Unit.Code)
				}
				_, err = tx.ExecContext(s.ctx, "insert into cors_hull (cors_id, unit_id, tech_level, efftn, endtn, qty_operational) values (?, ?, ?, ?, ?, ?)",
					cs.Id, hull.Unit.Id, hull.Unit.TechLevel, hull.EffTurn.String(), hull.EndTurn.String(), hull.QtyOperational)
				if err != nil {
					return fmt.Errorf("saveGame: cors_hull: insert: %w", err)
				}
			}

			for _, inventory := range cs.Inventory {
				if inventory.Unit.Id == 0 {
					inventory.Unit.Id = s.lookupUnitIdByCode(inventory.Unit.Code)
				}
				_, err = tx.ExecContext(s.ctx, "insert into cors_inventory (cors_id, unit_id, tech_level, efftn, endtn, qty_operational, qty_stowed) values (?, ?, ?, ?, ?, ?, ?)",
					cs.Id, inventory.Unit.Id, inventory.Unit.TechLevel, inventory.EffTurn.String(), inventory.EndTurn.String(), inventory.QtyOperational, inventory.QtyStowed)
				if err != nil {
					return fmt.Errorf("saveGame: cors_inventory: insert: %w", err)
				}
			}

			for _, group := range cs.Factories {
				if group.Unit.Id == 0 {
					group.Unit.Id = s.lookupUnitIdByCode(group.Unit.Code)
				}
				r, err := tx.ExecContext(s.ctx, "insert into cors_factory_group (cors_id, group_no, efftn, endtn, unit_id) values (?, ?, ?, ?, ?)",
					cs.Id, group.No, group.EffTurn.String(), group.EndTurn.String(), group.Unit.Id)
				if err != nil {
					return fmt.Errorf("saveGame: cors_factory_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: cors_factory_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					if unit.Unit.Id == 0 {
						unit.Unit.Id = s.lookupUnitIdByCode(unit.Unit.Code)
					}
					_, err = tx.ExecContext(s.ctx, "insert into cors_factory_group_units (factory_group_id, efftn, endtn, unit_id, qty_operational) values (?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn.String(), unit.EndTurn.String(), unit.Unit.Id, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: cors_factory_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into cors_factory_group_stages (factory_group_id, turn, unit_id, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?, ?)",
						group.Id, stage.Turn.String(), group.Unit.Id, stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: cors_factory_group_stages: insert: %w", err)
					}
				}
			}

			for _, group := range cs.Farms {
				if group.Unit.Id == 0 {
					group.Unit.Id = s.lookupUnitIdByCode(group.Unit.Code)
				}
				r, err := tx.ExecContext(s.ctx, "insert into cors_farm_group (cors_id, group_no, efftn, endtn, unit_id) values (?, ?, ?, ?, ?)",
					cs.Id, group.No, group.EffTurn.String(), group.EndTurn.String(), group.Unit.Id)
				if err != nil {
					return fmt.Errorf("saveGame: cors_farm_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: cors_farm_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					if unit.Unit.Id == 0 {
						unit.Unit.Id = s.lookupUnitIdByCode(unit.Unit.Code)
					}
					_, err = tx.ExecContext(s.ctx, "insert into cors_farm_group_units (farm_group_id, efftn, endtn, unit_id, qty_operational) values (?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn.String(), unit.EndTurn.String(), unit.Unit.Id, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: cors_farm_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into cors_farm_group_stages (farm_group_id, turn, unit_id, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?, ?)",
						group.Id, stage.Turn.String(), group.Unit.Id, stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: cors_farm_group_stages: insert: %w", err)
					}
				}
			}

			for _, group := range cs.Mines {
				if group.Deposit.Unit.Id == 0 {
					group.Deposit.Unit.Id = s.lookupUnitIdByCode(group.Deposit.Unit.Code)
				}
				r, err := tx.ExecContext(s.ctx, "insert into cors_mining_group (cors_id, group_no, efftn, endtn, resource_id) values (?, ?, ?, ?, ?)",
					cs.Id, group.No, group.EffTurn.String(), group.EndTurn.String(), group.Deposit.Id)
				if err != nil {
					return fmt.Errorf("saveGame: cors_mining_group: insert: %w", err)
				}
				id, err := r.LastInsertId()
				if err != nil {
					return fmt.Errorf("saveGame: cors_mining_group: lastInsertId: %w", err)
				}
				group.Id = int(id)
				for _, unit := range group.Units {
					if unit.Unit.Id == 0 {
						unit.Unit.Id = s.lookupUnitIdByCode(unit.Unit.Code)
					}
					_, err = tx.ExecContext(s.ctx, "insert into cors_mining_group_units (mining_group_id, efftn, endtn, unit_id, qty_operational) values (?, ?, ?, ?, ?)",
						group.Id, unit.EffTurn.String(), unit.EndTurn.String(), unit.Unit.Id, unit.QtyOperational)
					if err != nil {
						return fmt.Errorf("saveGame: cors_mining_group_units: insert: %w", err)
					}
				}
				for _, stage := range group.Stages {
					_, err = tx.ExecContext(s.ctx, "insert into cors_mining_group_stages (mining_group_id, turn, unit_id, qty_stage_1, qty_stage_2, qty_stage_3, qty_stage_4) values (?, ?, ?, ?, ?, ?, ?)",
						group.Id, stage.Turn.String(), group.Deposit.Unit.Id, stage.QtyStage1, stage.QtyStage2, stage.QtyStage3, stage.QtyStage4)
					if err != nil {
						return fmt.Errorf("saveGame: cors_mining_group_stages: insert: %w", err)
					}
				}
			}
			log.Printf("created nation %3d: cors %3d %8d\n", nation.No, cs.MSN, cs.Id)
		}
	}

	// link the players to their nation
	for _, player := range g.Players {
		_, err := tx.ExecContext(s.ctx, "insert into nation_player (nation_id, player_id) values (?, ?)", player.MemberOf.Id, player.Id)
		if err != nil {
			return fmt.Errorf("saveGame: memberOf: insert: %w", err)
		}
	}

	return tx.Commit()
}
