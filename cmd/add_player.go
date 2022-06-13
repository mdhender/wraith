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

var globalAddPlayer struct {
	Name     string
	LongName string
}

var cmdAddPlayer = &cobra.Command{
	Use:   "player",
	Short: "add a new player",
	Long:  `Add a new player to the game.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		globalAddPlayer.Name = strings.TrimSpace(globalAddPlayer.Name)
		if globalAddPlayer.Name == "" {
			return errors.New("missing player name")
		}
		// validate name
		for _, r := range globalAddPlayer.Name {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-') {
				return errors.New("invalid rune in player name")
			}
		}

		cfg := config.Players{
			FileName: filepath.Join(gamePath, "players.json"),
			Players: []config.Player{
				{Name: globalAddPlayer.Name, Description: globalAddPlayer.LongName},
			},
		}

		if _, err := os.Stat(cfg.FileName); err == nil {
			log.Printf("overwriting players store %q\n", cfg.FileName)
			if !globalAdd.Force {
				log.Fatal("cowardly refusing to overwrite existing players store")
			}
			log.Printf("overwriting players store %q\n", cfg.FileName)
		} else {
			log.Printf("creating players store %q\n", cfg.FileName)
		}
		if err := cfg.Write(); err != nil {
			return err
		}

		log.Printf("created players store %q\n", cfg.FileName)
		return nil
	},
}

func init() {
	cmdAddPlayer.Flags().StringVar(&globalAddPlayer.Name, "name", "", "name of new game (eg PT-1)")
	_ = cmdAddPlayer.MarkFlagRequired("name")
	cmdAddPlayer.Flags().StringVar(&globalAddPlayer.LongName, "long-name", "", "descriptive name of new game")
	//_ = cmdAddPlayer.MarkFlagRequired("long-name")

	cmdAdd.AddCommand(cmdAddPlayer)
}
