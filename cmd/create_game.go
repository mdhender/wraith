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
	"context"
	"errors"
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

var globalCreateGame struct {
	Force           bool
	ShortName       string
	Name            string
	NumberOfNations int
	StartDate       string
}

var cmdCreateGame = &cobra.Command{
	Use:   "game",
	Short: "create new game",
	Long:  `Create a new game.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		// validate the new game information
		globalCreateGame.ShortName = strings.ToUpper(strings.TrimSpace(globalCreateGame.ShortName))
		if globalCreateGame.ShortName == "" {
			return errors.New("missing short name")
		}
		globalCreateGame.Name = strings.TrimSpace(globalCreateGame.Name)
		if globalCreateGame.Name == "" {
			globalCreateGame.Name = globalCreateGame.ShortName
		}
		if !(0 < globalCreateGame.NumberOfNations && globalCreateGame.NumberOfNations < 225) {
			log.Fatalf("number of nations must be 1..225\n")
		}

		cfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded config %q\n", cfg.Self)

		e, err := engine.Open(cfg, context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// don't create if the game already exists
		if e.LookupGameByName(globalCreateGame.ShortName) != nil {
			log.Printf("short name %q already exists\n", globalCreateGame.ShortName)
			if !globalCreateGame.Force {
				log.Fatal("unable to create game\n")
			}
			err = e.DeleteGameByName(globalCreateGame.ShortName)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("short name %q purged\n", globalCreateGame.ShortName)
		}

		err = e.CreateGame(globalCreateGame.ShortName, globalCreateGame.Name, globalCreateGame.Name, 8, 14, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created game %q\n", globalCreateGame.ShortName)

		return nil
	},
}

func init() {
	cmdCreateGame.Flags().StringVar(&globalCreateGame.ShortName, "short-name", "", "report code for new game (eg PT-1)")
	_ = cmdCreateGame.MarkFlagRequired("short-name")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.Name, "name", "", "descriptive name of new game")
	cmdCreateGame.Flags().IntVar(&globalCreateGame.NumberOfNations, "nations", 20, "number of nations in game")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.StartDate, "start-date", "", "start date for game")
	cmdCreateGame.Flags().BoolVar(&globalCreateGame.Force, "force", false, "delete any existing game")

	cmdCreate.AddCommand(cmdCreateGame)
}
