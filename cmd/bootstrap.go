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
	"github.com/mdhender/wraith/engine"
	"github.com/spf13/cobra"
	"log"
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

		_, err := engine.Bootstrap(globalBase.ConfigFile, globalBootstrap.Store, globalBootstrap.Force)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("bootstrapped new engine\n")

		return nil
	},
}

func init() {
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Store, "store", "", "path to store data files")
	_ = cmdBootstrap.MarkFlagRequired("store")
	cmdBootstrap.Flags().BoolVar(&globalBootstrap.Force, "force", false, "force overwrite of existing configuration")

	cmdBase.AddCommand(cmdBootstrap)
}
