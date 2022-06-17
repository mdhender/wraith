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
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var globalCreateGame struct {
	Name        string
	Description string
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

		// validate the new game description
		globalCreateGame.Description = strings.TrimSpace(globalCreateGame.Description)
		if globalCreateGame.Description == "" {
			globalCreateGame.Description = globalCreateGame.Name
		}
		for _, r := range globalCreateGame.Description {
			if r == '\'' || r == '"' || r == '`' || r == '&' || r == '<' || r == '>' || unicode.IsControl(r) {
				return errors.New("invalid rune in game long name")
			}
		}

		// load the base configuration to find the games store
		cfgBase, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfgBase.Self)

		// load the games store
		cfgGames, err := engine.LoadGames("D:\\wraith\\testdata")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded games store %q\n", cfgGames.Store)

		// error on duplicate name
		for _, g := range cfgGames.Index {
			if strings.ToLower(g.Id) == strings.ToLower(globalCreateGame.Name) {
				log.Fatalf("duplicate game name %q", globalCreateGame.Name)
			}
		}

		// create the folders for the new game store
		gameFolder := filepath.Clean(filepath.Join("D:\\wraith\\testdata", "game", globalCreateGame.Name))
		if _, err = os.Stat(gameFolder); err != nil {
			log.Printf("creating game folder %q\n", gameFolder)
			if err = os.MkdirAll(gameFolder, 0700); err != nil {
				log.Fatal(err)
			}
			log.Printf("created game folder %q\n", gameFolder)
		}
		for _, dir := range []string{"nation", "nations", "turns"} {
			folder := filepath.Join(gameFolder, dir)
			if _, err = os.Stat(folder); err != nil {
				log.Printf("creating game %s folder %q\n", dir, folder)
				if err = os.MkdirAll(folder, 0700); err != nil {
					log.Fatal(err)
				}
				log.Printf("created game %s folder %q\n", dir, folder)
			}
		}

		// create the game store
		cfgGame, err := engine.CreateGame(strings.ToUpper(globalCreateGame.Name), globalCreateGame.Description, gameFolder, false)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created game store %q\n", cfgGame.Store)

		// create the nations store
		cfgNations, err := engine.CreateNations(cfgGame.Store, false)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created nations store %q\n", cfgNations.Store)

		// create the turns store
		cfgTurns, err := engine.CreateTurns(cfgGame.Store, false)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created turns store %q\n", cfgTurns.Store)

		// and update the games store
		cfgGames.Index = append(cfgGames.Index, engine.GamesIndex{
			Id:    cfgGame.Id,
			Store: cfgGame.Store,
		})
		//if err := cfgGames.Write(); err != nil {
		//	log.Printf("internal error - games store corrupted\n")
		//	log.Fatalf("%+v\n", err)
		//}
		log.Printf("updated games store %q\n", cfgGames.Store)

		return nil
	},
}

func init() {
	cmdCreateGame.Flags().StringVar(&globalCreateGame.Name, "name", "", "name of new game (eg PT-1)")
	_ = cmdCreateGame.MarkFlagRequired("name")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.Description, "descr", "", "descriptive name of new game")
	_ = cmdCreateGame.MarkFlagRequired("descr")

	cmdCreate.AddCommand(cmdCreateGame)
}
