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
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"
)

// Game configuration
type Game struct {
	Id        int
	ShortName string
	Name      string
	Descr     string
	Turn      int // index to current game turn
	Turns     []*Turn
	Nations   []*Nation
	Systems   []*System
}

// CreateGame creates a new game in the engine.
// It will replace any game currently in the engine.
func (e *Engine) CreateGame(shortName, name, descr string, numberOfNations, radius int, startDt time.Time) error {
	if e == nil {
		return ErrNoEngine
	} else if e.db == nil {
		return ErrNoStore
	}

	e.reset()

	shortName = strings.ToUpper(strings.TrimSpace(shortName))
	if shortName == "" {
		return fmt.Errorf("short name: %w", ErrMissingField)
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
	err := e.deleteGameByName(shortName)
	if err != nil {
		return fmt.Errorf("createGame: %w", err)
	}

	systemsPerRing := numberOfNations
	totalSystems := radius * systemsPerRing
	log.Printf("createGame: systems per ring %3d estimated systems %6d\n", systemsPerRing, totalSystems)
	rings := mkrings(radius, systemsPerRing)
	numPoints := 0
	for d := 0; d <= radius; d++ {
		numPoints += len(rings[d])
		log.Printf("createGame: ring %2d: %5d\n", d, len(rings[d]))
	}
	log.Printf("createGame:   total: %5d\n", numPoints)

	e.game = &Game{
		ShortName: shortName,
		Name:      name,
		Descr:     descr,
		Turn:      0,
	}

	turnDuration := 2 * 7 * 24 * time.Hour // assume two-week turns
	effDt := startDt
	endDt := effDt.Add(turnDuration)
	for t := 0; t < 10; t++ {
		e.game.Turns = append(e.game.Turns, &Turn{No: t, EffDt: effDt, EndDt: endDt})
		effDt = endDt
		endDt = effDt.Add(turnDuration)
	}

	systemId, ring, colonyNo := 0, 5, 0

	// generate nations and their home systems
	for i := 0; i < numberOfNations; i++ {
		systemId++
		coords := rings[ring][0]
		rings[ring] = rings[ring][1:]

		system := e.genHomeSystem(systemId)
		system.Ring, system.X, system.Y, system.Z = ring, coords.X, coords.Y, coords.Z
		e.game.Systems = append(e.game.Systems, system)

		nation := e.createNation(i+1, system.Stars[0].Orbits[3])
		nation.HomePlanet.Location = system.Stars[0].Orbits[3]
		colonyNo++
		nation.Colonies[0].No = colonyNo
		nation.Colonies[0].Location = nation.HomePlanet.Location
		colonyNo++
		nation.Colonies[1].No = colonyNo
		nation.Colonies[1].Location = nation.HomePlanet.Location
		e.game.Nations = append(e.game.Nations, nation)
	}

	// generate the remainder of the systems
	for ring := 0; ring < len(rings); ring++ {
		for _, coords := range rings[ring] {
			systemId++
			system := e.genSystem(systemId)
			system.Ring, system.X, system.Y, system.Z = ring, coords.X, coords.Y, coords.Z
			e.game.Systems = append(e.game.Systems, system)
		}
	}

	return e.saveGame()
}

func (e *Engine) DeleteGameByName(shortName string) error {
	return e.deleteGameByName(shortName)
}

func (e *Engine) LookupGameByName(shortName string) *Game {
	return e.lookupGameByName(shortName)
}

// ReadGame loads a store from a JSON file.
// It returns any errors.
func (e *Engine) ReadGame(id string) error {
	panic("!")
}

// WriteGame writes a store to a JSON file.
// It returns any errors.
func (e *Engine) WriteGame() error {
	panic("!")
}

// validateGameDescription validates the game description
func validateGameDescription(descr string) error {
	for _, r := range descr {
		if r == '\'' || r == '"' || r == '`' || r == '&' || r == '<' || r == '>' || r == '/' || r == '\\' || r == '$' || r == ';' || r == '{' || r == '}' || r == '!' || unicode.IsControl(r) {
			return errors.New("invalid rune in description")
		}
	}
	return nil
}

// validateGameId validates the game id
func validateGameId(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("missing id")
	}
	for _, r := range id {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-') {
			return errors.New("invalid rune in id")
		}
	}
	return nil
}
