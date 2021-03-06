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
	"github.com/mdhender/wraith/models"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

var globalBootstrap struct {
	// overwrite any existing configuration only if set.
	Force bool

	// database configuration
	User       string
	Password   string
	OrdersPath string
	Schema     string
	SchemaFile string
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

		if globalBootstrap.User == "" {
			return errors.New("missing database user")
		}
		if globalBootstrap.Password == "" {
			return errors.New("missing database password")
		}
		if globalBootstrap.OrdersPath == "" {
			return errors.New("missing orders path")
		}
		if globalBootstrap.Schema == "" {
			return errors.New("missing database schema name")
		}
		if globalBootstrap.SchemaFile == "" {
			return errors.New("missing schema generation file name")
		}

		cfg, err := config.CreateGlobal(globalBase.ConfigFile, globalBootstrap.User, globalBootstrap.Password, globalBootstrap.Schema, globalBootstrap.SchemaFile, globalBootstrap.OrdersPath, globalBootstrap.Force)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created %q\n", globalBase.ConfigFile)

		s, err := models.Bootstrap(cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()
		log.Printf("bootstrapped models version %s\n", s.Version())

		return nil
	},
}

func init() {
	cmdBootstrap.Flags().StringVar(&globalBootstrap.User, "user", "", "database user name")
	_ = cmdBootstrap.MarkFlagRequired("user")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Password, "password", "", "database password for user")
	_ = cmdBootstrap.MarkFlagRequired("password")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.OrdersPath, "orders-path", "", "path to orders files")
	_ = cmdBootstrap.MarkFlagRequired("orders-path")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Schema, "schema", "", "schema name in database")
	_ = cmdBootstrap.MarkFlagRequired("schema")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.SchemaFile, "schema-file", "", "path to schema generation file")
	_ = cmdBootstrap.MarkFlagRequired("schema-file")
	cmdBootstrap.Flags().BoolVar(&globalBootstrap.Force, "force", false, "force overwrite of existing configuration")

	cmdBase.AddCommand(cmdBootstrap)
}
