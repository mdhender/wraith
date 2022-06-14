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
	"fmt"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"strings"
	"unicode"
)

var globalAddNation struct {
	Game       string
	UserHandle string
	Name       string
	LongName   string
}

var cmdAddNation = &cobra.Command{
	Use:   "nation",
	Short: "add a new nation",
	Long:  `Add a new nation to the game.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		// validate the new nation name
		globalAddNation.Name = strings.TrimSpace(globalAddNation.Name)
		if globalAddNation.Name == "" {
			return errors.New("missing nation name")
		}
		for _, r := range globalAddNation.Name {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == ' ') {
				return errors.New("invalid rune in nation name")
			}
		}

		// validate the new nation long name
		globalAddNation.LongName = strings.TrimSpace(globalAddNation.LongName)
		if globalAddNation.LongName == "" {
			return errors.New("missing nation long name")
		}

		// load the base configuration to find the games store
		cfgBase, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfgBase.Self)

		// load the games store to find the game store
		cfgGames, err := config.LoadGames(cfgBase.Store)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded games store %q\n", cfgGames.Store)

		// find the game in the store
		var cfgGame *config.Game
		for _, g := range cfgGames.Index {
			if strings.ToLower(g.Name) == strings.ToLower(globalAddNation.Game) {
				cfgGame, err = config.LoadGame(g.Store)
				if err != nil {
					log.Fatal(err)
				}
				break
			}
		}
		if cfgGame == nil {
			log.Fatalf("unable to find game %q\n", globalAddNation.Game)
		}
		log.Printf("loaded game store %q\n", cfgGame.Store)

		// use the game store to load the nations store
		cfgNations, err := config.LoadNations(cfgGame.Store)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded nations store %q\n", cfgNations.Store)

		// generate an id for the new nation
		id := len(cfgNations.Index) + 1

		// check for duplicates in the nations store
		for _, n := range cfgNations.Index {
			if strings.ToLower(n.Name) == strings.ToLower(globalAddNation.Name) {
				log.Fatalf("error: duplicate nation name %q\n", globalAddNation.Name)
			} else if n.Id == id {
				log.Fatalf("error: duplicate nation id %q\n", id)
			}
		}

		// add the new nation to the nations store
		nationIndex := config.NationsIndex{
			Id:   id,
			Name: globalAddNation.Name,
			Path: filepath.Clean(filepath.Join(cfgGame.Store, fmt.Sprintf("%d", id))),
		}
		cfgNations.Index = append(cfgNations.Index, nationIndex)

		log.Printf("updating nations store %q\n", cfgNations.Store)
		if err := cfgNations.Write(); err != nil {
			return err
		}

		log.Printf("updated nations store %q\n", cfgNations.Store)

		log.Printf("created nation %d %q\n", nationIndex.Id, nationIndex.Name)
		return nil
	},
}

func init() {
	cmdAddNation.Flags().StringVar(&globalAddNation.Game, "game", "", "name of game to add nation to")
	_ = cmdAddNation.MarkFlagRequired("game")
	cmdAddNation.Flags().StringVar(&globalAddNation.UserHandle, "user", "", "handle of user controlling nation")
	_ = cmdAddNation.MarkFlagRequired("name")
	cmdAddNation.Flags().StringVar(&globalAddNation.Name, "name", "", "name of new nation")
	_ = cmdAddNation.MarkFlagRequired("name")
	cmdAddNation.Flags().StringVar(&globalAddNation.LongName, "long-name", "", "descriptive name of new game")
	//_ = cmdAddPlayer.MarkFlagRequired("long-name")

	cmdAdd.AddCommand(cmdAddNation)
}
