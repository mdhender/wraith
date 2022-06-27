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
	"path/filepath"
	"strings"
	"time"
)

var globalCreateGame struct {
	Force     bool
	ShortName string
	Name      string
	Players   string // location of player data
	Radius    int
	StartDate string
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
		if !(2 < globalCreateGame.Radius && globalCreateGame.Radius < 18) {
			log.Fatalf("radius must be 3..17\n")
		}

		globalCreateGame.Players = filepath.Clean(globalCreateGame.Players)
		b, err := os.ReadFile(globalCreateGame.Players)
		if err != nil {
			log.Fatal(err)
		}
		var data struct {
			Game    string `json:"game"`
			Players []struct {
				Id           string `json:"id"`
				UserHandle   string `json:"user"`
				PlayerHandle string `json:"handle"`
				Nation       struct {
					Name       string `json:"name"`
					Speciality string `json:"speciality"`
					HomeWorld  string `json:"home-world"`
					GovtKind   string `json:"govt-kind"`
					GovtName   string `json:"govt-name"`
				} `json:"nation"`
			} `json:"players"`
		}
		err = json.Unmarshal(b, &data)
		if err != nil {
			log.Fatal(err)
		} else if data.Game != globalCreateGame.ShortName {
			log.Fatalf("name in data file does not match command line\n")
		} else if !(0 < len(data.Players) && len(data.Players) < 225) {
			log.Fatalf("number of players must be 1..225\n")
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

		// don't create if the game already exists
		if game, err := s.LookupGameByName(globalCreateGame.ShortName); game != nil {
			log.Printf("short name %q already exists\n", globalCreateGame.ShortName)
			if !globalCreateGame.Force {
				log.Fatal("unable to create game\n")
			}
			err = s.DeleteGameByName(globalCreateGame.ShortName)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("short name %q purged\n", globalCreateGame.ShortName)
		}

		// convert players from data file to positions for game
		var positions []*models.PlayerPosition
		for _, player := range data.Players {
			position := &models.PlayerPosition{
				UserHandle:   player.UserHandle,
				PlayerHandle: player.PlayerHandle,
			}
			position.Nation.Name = player.Nation.Name
			position.Nation.Speciality = player.Nation.Speciality
			position.Nation.HomeWorld = player.Nation.HomeWorld
			position.Nation.GovtKind = player.Nation.GovtKind
			position.Nation.GovtName = player.Nation.GovtName
			positions = append(positions, position)
		}

		game, err := s.GenerateGame(globalCreateGame.ShortName, globalCreateGame.Name, "", globalCreateGame.Radius, time.Now(), positions)
		if err != nil {
			log.Fatal(err)
		}
		err = s.SaveGame(game)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("created game %q\n", game.ShortName)

		return nil
	},
}

func init() {
	cmdCreateGame.Flags().StringVar(&globalCreateGame.ShortName, "short-name", "", "report code for new game (eg PT-1)")
	_ = cmdCreateGame.MarkFlagRequired("short-name")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.Name, "name", "", "descriptive name of new game")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.Players, "players", "", "name of players data file")
	_ = cmdCreateGame.MarkFlagRequired("players")
	cmdCreateGame.Flags().IntVar(&globalCreateGame.Radius, "radius", 8, "radius of cluster")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.StartDate, "start-date", "", "start date for game")
	cmdCreateGame.Flags().BoolVar(&globalCreateGame.Force, "force", false, "delete any existing game")

	cmdCreate.AddCommand(cmdCreateGame)
}
