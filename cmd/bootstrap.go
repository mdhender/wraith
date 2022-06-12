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
	"fmt"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var globalBootstrap struct {
	// overwrite any existing configuration only if set.
	Force bool
	// location of global configuration file.
	ConfigFile string
	// location to create game data files in.
	GamesPath string
}

var cmdBootstrap = &cobra.Command{
	Use:   "bootstrap",
	Short: "create a new global configuration file",
	Long: `Create the initial system configuration.
This includes the configuration file, games path, and starting data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		globalBootstrap.ConfigFile = globalBase.ConfigFile
		if globalBootstrap.ConfigFile == "" {
			return fmt.Errorf("missing config file name")
		}

		if globalBootstrap.GamesPath == "" {
			return fmt.Errorf("missing games path")
		}

		cfg := config.Config{
			ConfigFile: globalBootstrap.ConfigFile,
			GamesPath:  globalBootstrap.GamesPath,
		}

		if _, err := os.Stat(globalBootstrap.ConfigFile); err == nil {
			if !globalBootstrap.Force {
				log.Fatal("cowardly refusing to overwrite existing configuration store")
			}
			log.Printf("overwriting config store %q\n", globalBootstrap.ConfigFile)
		} else {
			log.Printf("creating config store %q\n", globalBootstrap.ConfigFile)
		}
		if err := config.Write(globalBootstrap.ConfigFile, &cfg); err != nil {
			return err
		}

		log.Printf("created config store %q\n", globalBootstrap.ConfigFile)
		return nil
	},
}

func init() {
	cmdBootstrap.Flags().StringVar(&globalBootstrap.GamesPath, "games-path", "", "path to create new game data")
	_ = cmdBootstrap.MarkFlagRequired("games-path")
	cmdBootstrap.Flags().BoolVar(&globalBootstrap.Force, "force", false, "force overwrite of existing configuration")

	cmdBase.AddCommand(cmdBootstrap)
}
