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
	"encoding/json"
	"errors"
	"github.com/mdhender/wraith/models"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var globalCreateGame struct {
	ShortName string
	Name      string
	StartDt   string
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

		b, err := os.ReadFile("D:\\wraith\\testdata\\systems.json")
		if err != nil {
			log.Fatal(err)
		}
		var systems []struct {
			X           int  `json:"x"`
			Y           int  `json:"y"`
			Z           int  `json:"z"`
			Stars       int  `json:"stars"`
			Singularity bool `json:"singularity,omitempty"`
		}
		if err := json.Unmarshal(b, &systems); err != nil {
			log.Fatal(err)
		}

		cfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfg.Self)

		s, err := models.Open(cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()

		g, err := s.CreateGame(globalCreateGame.Name, globalCreateGame.ShortName, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("game id %d %v\n", g.Id, g)

		for _, system := range systems {
			s, err := s.AddSystem(g, system.X, system.Y, system.Z)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("game id %d %v\n", g.Id, s)
		}

		return nil
	},
}

func init() {
	cmdCreateGame.Flags().StringVar(&globalCreateGame.ShortName, "short-name", "", "report code for new game (eg PT-1)")
	_ = cmdCreateGame.MarkFlagRequired("short-name")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.Name, "name", "", "descriptive name of new game")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.StartDt, "start-date", time.Now().String(), "start date for game")

	cmdCreate.AddCommand(cmdCreateGame)
}
