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
	"github.com/mdhender/wraith/storage/config"
	"github.com/pkg/errors"
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
		globalCreateGame.Name = strings.TrimSpace(globalCreateGame.Name)
		if globalCreateGame.Name == "" {
			return errors.New("missing config game name")
		}
		// validate name
		for _, r := range globalCreateGame.Name {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-') {
				return errors.New("invalid rune in game name")
			}
		}
		globalCreateGame.LongName = strings.TrimSpace(globalCreateGame.LongName)
		if globalCreateGame.LongName == "" {
			globalCreateGame.LongName = globalCreateGame.Name
		}
		// validate long name
		for _, r := range globalCreateGame.LongName {
			if r == '\'' || r == '"' || r == '`' || r == '&' || r == '<' || r == '>' || unicode.IsControl(r) {
				return errors.New("invalid rune in game long name")
			}
		}

		gCfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			return err
		}
		log.Printf("loaded %q\n", gCfg.FileName)

		gamePath := filepath.Join(gCfg.GamesPath, globalCreateGame.Name)
		if _, err = os.Stat(gamePath); err != nil {
			log.Printf("creating game folder %q\n", gamePath)
			if err = os.MkdirAll(gamePath, 0700); err != nil {
				return err
			}
		}

		cfg := config.Game{
			FileName:    filepath.Join(gamePath, "game.json"),
			Name:        globalCreateGame.Name,
			Description: globalCreateGame.LongName,
			GamePath:    gamePath,
		}

		if _, err := os.Stat(cfg.FileName); err == nil {
			if !globalCreate.Force {
				log.Fatal("cowardly refusing to overwrite existing game store")
			}
			log.Printf("overwriting game store %q\n", cfg.FileName)
		} else {
			log.Printf("creating game store %q\n", cfg.FileName)
		}
		if err := cfg.Write(); err != nil {
			return err
		}

		log.Printf("created game store %q\n", cfg.FileName)
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
