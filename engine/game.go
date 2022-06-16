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
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// Game configuration
type Game struct {
	Store        string         `json:"store"` // path to store data
	Id           string         `json:"id"`
	Description  string         `json:"description"`
	NationsIndex []NationsIndex `json:"nations-index"`
	TurnsIndex   []TurnsIndex   `json:"turns-index"`
}

// AddGame adds a game to the store.
func (e *Engine) AddGame(id, descr string) error {
	// free up any game already in memory
	e.stores.game, e.stores.nations, e.stores.turns = nil, nil, nil

	// validate the id and description
	id = strings.TrimSpace(id)
	if descr = strings.TrimSpace(descr); descr == "" {
		descr = id
	}
	if err := validateGameId(id); err != nil {
		return err
	}
	if err := validateGameDescription(descr); err != nil {
		return err
	}

	// error on duplicate id
	for _, g := range e.stores.games.Index {
		if strings.ToLower(g.Id) == strings.ToLower(id) {
			return errors.New("duplicate id")
		}
	}

	e.stores.game = &Game{
		Store:        e.RootDir("game", id),
		Id:           id,
		Description:  descr,
		NationsIndex: []NationsIndex{},
		TurnsIndex:   []TurnsIndex{},
	}

	e.stores.nations = &Nations{
		Store: e.RootDir("game", id, "nations"),
		Index: []NationsIndex{},
	}

	e.stores.turns = &Turns{
		Store: e.RootDir("game", id, "turns"),
		Index: []TurnsIndex{},
	}

	// create the folders for the game store
	for _, folder := range []string{
		e.stores.game.Store,
		e.stores.nations.Store,
		e.stores.turns.Store,
		filepath.Join(e.stores.game.Store, "nation"),
	} {
		if _, err := os.Stat(folder); err != nil {
			log.Printf("creating folder %q\n", folder)
			if err = os.MkdirAll(folder, 0700); err != nil {
				log.Fatal(err)
			}
			log.Printf("created folder %q\n", folder)
		}
	}

	if err := e.WriteGame(); err != nil {
		return err
	}
	if err := e.WriteNations(); err != nil {
		return err
	}
	if err := e.WriteTurns(); err != nil {
		return err
	}

	panic("!")
}

// ReadGame loads a store from a JSON file.
// It returns any errors.
func (e *Engine) ReadGame(id string) error {
	// free up any game already in memory
	e.stores.game, e.stores.nations, e.stores.turns = nil, nil, nil

	// validate the id
	id = strings.TrimSpace(id)
	if err := validateGameId(id); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(filepath.Join(e.RootDir("game", id), "store.json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, e.stores.game)
}

// WriteGame writes a store to a JSON file.
// It returns any errors.
func (e *Engine) WriteGame() error {
	b, err := json.MarshalIndent(e.stores.game, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(e.stores.game.Store, "store.json"), b, 0600)
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
