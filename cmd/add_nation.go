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
		// TODO: this has to be updated to use the base configuration and command line flag for game name
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		gamePath, err := os.Getwd()
		if err != nil {
			return err
		}
		gameCfgFile := filepath.Join(gamePath, "game.json")
		gCfg, err := config.LoadGame(gameCfgFile)
		if err != nil {
			return err
		} else if gCfg.FileName != gameCfgFile {
			log.Printf("loaded     %q\n", gameCfgFile)
			log.Printf("but wanted %q\n", gCfg.FileName)
			log.Fatal("error: internal error: path mismatch\n")
		}
		log.Printf("loaded %q\n", gameCfgFile)

		globalAddNation.Name = strings.TrimSpace(globalAddNation.Name)
		if globalAddNation.Name == "" {
			return errors.New("missing player name")
		}
		// validate name
		for _, r := range globalAddNation.Name {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-') {
				return errors.New("invalid rune in player name")
			}
		}

		// load the base configuration to find the games store
		baseCfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", baseCfg.FileName)

		// load the games store to find the game store
		gamesCfg, err := config.LoadGames(baseCfg.GamesStore)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", gamesCfg.FileName)

		// find the game in the store
		var game config.Game
		for _, g := range gamesCfg.Games {
			if g.Name != globalAddNation.Game {
				continue
			}
			game = g
			break
		}
		if game.Name != globalAddNation.Game {
			log.Fatalf("unable to find game %q\n", globalAddNation.Game)
		}

		// load the game store to find the nations store
		gameCfg, err := config.LoadGame(game.FileName)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", gameCfg.FileName)

		// load the nations store
		nationsCfg, err := config.LoadNations(gameCfg.NationsStore)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", gameCfg.FileName)

		nation := config.Nation{
			Name:        globalAddNation.Name,
			Description: globalAddNation.LongName,
		}

		// check for duplicates
		for _, n := range nationsCfg.Nations {
			if n.Name == nation.Name {
				log.Fatalf("error: duplicate nation name %q\n", nation.Name)
			}
		}

		nationsCfg.Nations = append(nationsCfg.Nations, nation)

		log.Printf("updating nations store %q\n", nationsCfg.FileName)
		if err := nationsCfg.Write(); err != nil {
			return err
		}

		log.Printf("updated nations store %q\n", nationsCfg.FileName)
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
