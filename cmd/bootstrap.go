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

package cmd

import (
	"errors"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var globalBootstrap struct {
	// overwrite any existing configuration only if set.
	Force bool
	// location to create game data files in.
	GamesPath string
	// location to create user data files in.
	UsersPath string
}

var cmdBootstrap = &cobra.Command{
	Use:   "bootstrap",
	Short: "create a new global configuration file",
	Long: `Create the initial system configuration.
This includes the configuration file, games path, and starting data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		// validate games path
		globalBootstrap.GamesPath = strings.TrimSpace(globalBootstrap.GamesPath)
		if globalBootstrap.GamesPath == "" {
			return errors.New("missing games path")
		}

		// validate users path
		globalBootstrap.UsersPath = strings.TrimSpace(globalBootstrap.UsersPath)
		if globalBootstrap.UsersPath == "" {
			return errors.New("missing users path")
		}

		cfgBase := config.Global{
			Path:       filepath.Clean(globalBase.ConfigFile),
			GamesPath:  filepath.Clean(globalBootstrap.GamesPath),
			GamesStore: filepath.Clean(filepath.Join(globalBootstrap.GamesPath, "games.json")),
			UsersPath:  filepath.Clean(globalBootstrap.UsersPath),
			UsersStore: filepath.Clean(filepath.Join(globalBootstrap.UsersPath, "users.json")),
		}

		log.Printf("intended config store %q\n", globalBase.ConfigFile)
		if _, err := os.Stat(globalBase.ConfigFile); err == nil {
			if !globalBootstrap.Force {
				log.Fatal("cowardly refusing to overwrite existing configuration store")
			}
			log.Printf("overwriting config store %q\n", globalBase.ConfigFile)
		} else {
			log.Printf("creating config store %q\n", globalBase.ConfigFile)
		}
		if err := cfgBase.Write(); err != nil {
			log.Fatal(err)
		}
		log.Printf("created config store %q\n", globalBase.ConfigFile)

		cfgGames := config.Games{
			GamesPath: cfgBase.GamesPath,
			Games:     []config.GamesIndex{},
		}
		if err := cfgGames.Create(cfgBase.GamesPath, true); err != nil {
			log.Fatal(err)
		}
		log.Printf("created games store %q\n", cfgGames.Path)

		cfgUsers := config.Users{
			Users: []config.User{},
		}
		if err := cfgUsers.Create(cfgBase.UsersPath, true); err != nil {
			log.Fatal(err)
		}
		log.Printf("created users store %q\n", cfgUsers.Path)

		return nil
	},
}

func init() {
	cmdBootstrap.Flags().StringVar(&globalBootstrap.GamesPath, "games-path", "", "path to create new games data")
	_ = cmdBootstrap.MarkFlagRequired("games-path")
	cmdBootstrap.Flags().BoolVar(&globalBootstrap.Force, "force", false, "force overwrite of existing configuration")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.UsersPath, "users-path", "", "path to create new users data")
	_ = cmdBootstrap.MarkFlagRequired("users-path")

	cmdBase.AddCommand(cmdBootstrap)
}
