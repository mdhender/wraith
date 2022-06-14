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
	// path to store data.
	Store string
}

var cmdBootstrap = &cobra.Command{
	Use:   "bootstrap",
	Short: "create a new global configuration file",
	Long: `Create the initial system configuration.
This includes the configuration file and starting data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}
		globalBase.ConfigFile = filepath.Clean(globalBase.ConfigFile)

		// validate data path
		globalBootstrap.Store = strings.TrimSpace(globalBootstrap.Store)
		if globalBootstrap.Store == "" {
			return errors.New("missing data path")
		}
		globalBootstrap.Store = filepath.Clean(globalBootstrap.Store)

		log.Printf("intended config file %q\n", globalBase.ConfigFile)
		if _, err := os.Stat(globalBase.ConfigFile); err == nil {
			if !globalBootstrap.Force {
				log.Fatal("cowardly refusing to overwrite existing configuration file")
			}
			log.Printf("overwriting config file %q\n", globalBase.ConfigFile)
		} else {
			log.Printf("creating config file %q\n", globalBase.ConfigFile)
		}
		cfgBase, err := config.CreateGlobal(globalBase.ConfigFile, globalBootstrap.Store, globalBootstrap.Force)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created config file %q\n", cfgBase.Self)

		// create the data store folder
		if _, err := os.Stat(cfgBase.Store); err != nil {
			log.Printf("creating data folder %q\n", cfgBase.Store)
			if err = os.MkdirAll(cfgBase.Store, 0700); err != nil {
				log.Fatal(err)
			}
			log.Printf("created data folder %q\n", cfgBase.Store)
		}

		// create the default users store
		usersFolder := filepath.Join(cfgBase.Store, "users")
		if _, err := os.Stat(usersFolder); err != nil {
			log.Printf("creating users folder %q\n", usersFolder)
			if err = os.MkdirAll(usersFolder, 0700); err != nil {
				log.Fatal(err)
			}
			log.Printf("created users folder %q\n", usersFolder)
		}
		cfgUsers, err := config.CreateUsers(cfgBase.Store, globalBootstrap.Force)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created users store %q\n", cfgUsers.Store)

		// create the default games store
		gamesFolder := filepath.Join(cfgBase.Store, "games")
		if _, err := os.Stat(gamesFolder); err != nil {
			log.Printf("creating games folder %q\n", gamesFolder)
			if err = os.MkdirAll(gamesFolder, 0700); err != nil {
				log.Fatal(err)
			}
			log.Printf("created games folder %q\n", gamesFolder)
		}
		cfgGames, err := config.CreateGames(cfgBase.Store, globalBootstrap.Force)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created games store %q\n", cfgGames.Store)

		return nil
	},
}

func init() {
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Store, "store", "", "path to store data files")
	_ = cmdBootstrap.MarkFlagRequired("store")
	cmdBootstrap.Flags().BoolVar(&globalBootstrap.Force, "force", false, "force overwrite of existing configuration")

	cmdBase.AddCommand(cmdBootstrap)
}
