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
	"github.com/mdhender/wraith/storage/config"
	"log"
	"path/filepath"
	"strings"
)

type Engine struct {
	config struct {
		base *config.Global
	}
	stores struct {
		games   *Games
		game    *Game
		nations *Nations
	}
	s *Store
}

// New returns an initialized engine with the base configuration
// and the games store loaded.
func New(baseConfigFile string) (e *Engine, err error) {
	if baseConfigFile == "" {
		return nil, errors.New("missing base config")
	}

	e = &Engine{}

	// load the base configuration and the games store
	e.config.base, err = config.LoadGlobal(baseConfigFile)
	if err != nil {
		return nil, err
	}
	log.Printf("loaded base config %q\n", e.config.base.Self)

	e.stores.games, err = LoadGames(filepath.Join(e.config.base.Store, "games"))
	if err != nil {
		return nil, err
	}
	log.Printf("loaded games store %q\n", e.stores.games.Store)

	return e, nil
}

func (e *Engine) LoadGame(game string) (err error) {
	// free up any game already in memory
	e.stores.game, e.stores.nations = nil, nil

	if game == "" {
		return errors.New("missing game name")
	}

	// find the game in the store
	for _, g := range e.stores.games.Index {
		if strings.ToLower(g.Name) == strings.ToLower(game) {
			e.stores.game, err = LoadGame(g.Store)
			if err != nil {
				return err
			}
			break
		}
	}
	if e.stores.game == nil {
		log.Fatalf("unable to find game %q\n", game)
	}
	log.Printf("loaded game store %q\n", e.stores.game.Store)

	// use the game store to load the nations store
	e.stores.nations, err = LoadNations(e.stores.game.Store)
	if err != nil {
		return err
	}
	log.Printf("loaded nations store %q\n", e.stores.nations.Store)

	return nil
}
