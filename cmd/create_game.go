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
	"unicode"
)

var globalCreateGame struct {
	Name     string
	LongName string
}

var cmdCreateGame = &cobra.Command{
	Use:   "game",
	Short: "create new game",
	Long:  `Create a new game.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		// validate the new game name
		globalCreateGame.Name = strings.TrimSpace(globalCreateGame.Name)
		if globalCreateGame.Name == "" {
			return errors.New("missing config game name")
		}
		for _, r := range globalCreateGame.Name {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-') {
				return errors.New("invalid rune in game name")
			}
		}

		// validate the new game long name
		globalCreateGame.LongName = strings.TrimSpace(globalCreateGame.LongName)
		if globalCreateGame.LongName == "" {
			globalCreateGame.LongName = globalCreateGame.Name
		}
		for _, r := range globalCreateGame.LongName {
			if r == '\'' || r == '"' || r == '`' || r == '&' || r == '<' || r == '>' || unicode.IsControl(r) {
				return errors.New("invalid rune in game long name")
			}
		}

		// load the base configuration to get the games store
		gCfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", gCfg.FileName)

		// load the games store
		gamesCfg, err := config.LoadGames(gCfg.GamesStore)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", gamesCfg.FileName)

		// create the game to add
		game := config.Game{
			FileName:    filepath.Join(gCfg.GamesPath, globalCreateGame.Name, "game.json"),
			GamePath:    filepath.Join(gCfg.GamesPath, globalCreateGame.Name),
			Name:        globalCreateGame.Name,
			Description: globalCreateGame.LongName,
		}

		// error on duplicate name
		for _, g := range gamesCfg.Games {
			if g.Name == game.Name {
				log.Fatalf("duplicate game name %q", g.Name)
			}
		}

		// create the folder for the new game store
		if _, err = os.Stat(game.GamePath); err != nil {
			log.Printf("creating game folder %q\n", game.GamePath)
			if err = os.MkdirAll(game.GamePath, 0700); err != nil {
				log.Fatal(err)
			}
		}

		// create the default players store
		players := config.Nations{
			FileName: filepath.Join(game.GamePath, "players.json"),
			Nations:  []config.Nation{},
		}
		if err := players.Write(); err != nil {
			log.Fatal(err)
		}
		log.Printf("created players store %q\n", players.FileName)

		// and update the games store
		gamesCfg.Games = append(gamesCfg.Games, game)
		if err := gamesCfg.Write(); err != nil {
			log.Printf("internal error - games store corrupted\n")
			log.Fatalf("%+v", err)
		}
		log.Printf("updated games store %q\n", gamesCfg.FileName)

		return nil
	},
}

func init() {
	cmdCreateGame.Flags().StringVar(&globalCreateGame.Name, "name", "", "name of new game (eg PT-1)")
	_ = cmdCreateGame.MarkFlagRequired("name")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.LongName, "long-name", "", "descriptive name of new game")
	_ = cmdCreateGame.MarkFlagRequired("long-name")

	cmdCreate.AddCommand(cmdCreateGame)
}
