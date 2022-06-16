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
	"os"
	"path/filepath"
	"strings"
)

// Bootstrap creates a new engine.
func Bootstrap(configFile, data string, overwrite bool) (*Engine, error) {
	// validate parameters
	if configFile == "" {
		return nil, errors.New("missing base config")
	}
	configFile = filepath.Clean(configFile)
	//
	data = strings.TrimSpace(data)
	if data == "" {
		return nil, errors.New("missing data path")
	}
	data = filepath.Clean(data)

	// create the default folders for a new engine
	folders := []string{
		data,
		filepath.Join(data, "game"),
		filepath.Join(data, "games"),
		filepath.Join(data, "user"),
		filepath.Join(data, "users"),
	}
	for _, folder := range folders {
		if _, err := os.Stat(folder); err != nil {
			log.Printf("creating folder %q\n", folder)
			if err = os.MkdirAll(folder, 0700); err != nil {
				log.Fatal(err)
			}
			log.Printf("created folder %q\n", folder)
		}
	}

	e := &Engine{}
	var err error

	// create the base configuration file
	if _, err := os.Stat(configFile); err == nil {
		if !overwrite {
			return nil, errors.New("configuration file exists")
		}
		log.Printf("overwriting config file %q\n", configFile)
	} else {
		log.Printf("creating config file %q\n", configFile)
	}
	e.config.base, err = config.CreateGlobal(configFile, data, overwrite)
	if err != nil {
		return nil, err
	}
	log.Printf("created config file %q\n", e.config.base.Self)

	// create a new games store
	e.stores.games = &Games{
		Store: e.RootDir("games"),
		Index: []GamesIndex{},
	}
	if err = e.WriteGames(); err != nil {
		log.Fatal(err)
	}
	log.Printf("created games store %q\n", e.stores.games.Store)

	// create a new users store
	e.stores.users = &Users{
		Store: e.RootDir("users"),
		Index: []UsersIndex{},
	}
	if err = e.WriteUsers(); err != nil {
		log.Fatal(err)
	}
	log.Printf("created users store %q\n", e.stores.users.Store)

	return e, nil
}
