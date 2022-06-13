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
)

var globalBootstrap struct {
	// overwrite any existing configuration only if set.
	Force bool
	// location to create game data files in.
	GamesPath string
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

		if globalBootstrap.GamesPath == "" {
			return errors.New("missing games path")
		}

		cfg := config.Global{
			FileName:   globalBase.ConfigFile,
			GamesPath:  globalBootstrap.GamesPath,
			GamesStore: filepath.Join(globalBootstrap.GamesPath, "games.json"),
			UsersStore: filepath.Join(globalBootstrap.GamesPath, "users.json"),
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
		if err := cfg.Write(); err != nil {
			log.Fatal(err)
		}
		log.Printf("created config store %q\n", globalBase.ConfigFile)

		cfgGames := config.Games{
			FileName: cfg.GamesStore,
			Games:    []config.Game{},
		}
		if err := cfgGames.Write(); err != nil {
			log.Fatal(err)
		}
		log.Printf("created games store %q\n", cfgGames.FileName)

		cfgUsers := config.Users{
			FileName: cfg.UsersStore,
			Users:    []config.User{},
		}
		if err := cfgUsers.Write(); err != nil {
			log.Fatal(err)
		}
		log.Printf("created users store %q\n", cfgUsers.FileName)

		return nil
	},
}

func init() {
	cmdBootstrap.Flags().StringVar(&globalBootstrap.GamesPath, "games-path", "", "path to create new game data")
	_ = cmdBootstrap.MarkFlagRequired("games-path")
	cmdBootstrap.Flags().BoolVar(&globalBootstrap.Force, "force", false, "force overwrite of existing configuration")

	cmdBase.AddCommand(cmdBootstrap)
}
