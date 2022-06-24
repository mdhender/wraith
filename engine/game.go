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

	e.game = &Game{
		ShortName: shortName,
		Name:      name,
		Descr:     descr,
		Turn:      0,
	}

	// assume two week turns
	effDt := startDt
	endDt := effDt.Add(2 * 7 * 24 * time.Hour)
	for t := 0; t < 10; t++ {
		e.game.Turns = append(e.game.Turns, &Turn{No: t, EffDt: effDt, EndDt: endDt})
		effDt = endDt
		endDt = effDt.Add(2 * 7 * 24 * time.Hour)
	}

	for i := 0; i < numberOfNations; i++ {
		e.game.Nations = append(e.game.Nations, e.createNation(i+1))
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
