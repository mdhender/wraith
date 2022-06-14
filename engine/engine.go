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
		base    *config.Global
		games   *config.Games
		game    *config.Game
		nations *config.Nations
	}
	s *Store
}

func LoadGame(baseStore string, game string) (e *Engine, err error) {
	e = &Engine{}

	if baseStore == "" {
		return nil, errors.New("missing base config")
	} else if game == "" {
		return nil, errors.New("missing game name")
	}

	// load the base configuration to find the games store
	e.config.base, err = config.LoadGlobal(baseStore)
	if err != nil {
		return nil, err
	}
	log.Printf("loaded %q\n", e.config.base.Path)

	// load the games store to find the game store
	e.config.games, err = config.LoadGames(filepath.Join(e.config.base.GamesPath, "games.json"))
	if err != nil {
		return nil, err
	}
	log.Printf("loaded games store %q\n", e.config.games.Path)

	// find the game in the store
	for _, g := range e.config.games.Games {
		if strings.ToLower(g.Name) == strings.ToLower(game) {
			e.config.game, err = config.LoadGame(g.Path)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if e.config.game == nil {
		log.Fatalf("unable to find game %q\n", game)
	}
	log.Printf("loaded game store %q\n", e.config.game.Path)

	// use the game store to load the nations store
	e.config.nations, err = config.LoadNations(e.config.game.NationsStore)
	if err != nil {
		return nil, err
	}
	log.Printf("loaded nations store %q\n", e.config.nations.Path)

	return e, nil
}
