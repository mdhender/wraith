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
)

var globalQuery struct {
	Id int
}

var cmdQuery = &cobra.Command{
	Use:   "query",
	Short: "test query",
	Long:  `Create a database connection and run a test query.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		cfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded config %q\n", cfg.Self)

		s, err := models.Open(cfg)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded store version %q\n", s.Version())

		if err = s.Ping(); err != nil {
			log.Fatal(err)
		}
		log.Printf("query: connection seems ok\n")

		return nil
	},
}

func init() {
	cmdQuery.Flags().IntVar(&globalQuery.Id, "id", 0, "id to query")
	_ = cmdQuery.MarkFlagRequired("id")

	cmdBase.AddCommand(cmdQuery)
}
