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
	"encoding/json"
	"errors"
	"github.com/mdhender/wraith/models"
	"github.com/mdhender/wraith/storage/config"
	"github.com/mdhender/wraith/storage/jdb"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var globalExport struct {
	File string
	Game string
	Turn string
}

var cmdExport = &cobra.Command{
	Use:   "export",
	Short: "export game",
	Long:  `Export a turn for a game.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}
		if globalExport.File = strings.TrimSpace(globalExport.File); globalExport.File == "" {
			return errors.New("missing export file name")
		}
		globalExport.File = filepath.Clean(globalExport.File)

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

		game, err := s.LookupGameByName(globalExport.Game)
		if err != nil {
			log.Fatal(err)
		} else if game, err = s.FetchGameByNameAsOf(game.ShortName, game.CurrentTurn.String()); err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded game %q: turn %q\n", game.ShortName, game.CurrentTurn.String())
		if globalExport.Turn == "" {
			globalExport.Turn = game.CurrentTurn.String()
		}

		if gj, err := jdb.Extract(s.GetDB(), context.Background(), game.Id); err != nil {
			log.Fatal(err)
		} else if b, err := json.MarshalIndent(gj, "", "\t"); err != nil {
			log.Fatal(err)
		} else if err := os.WriteFile(globalExport.File, b, 0666); err != nil {
			log.Fatal(err)
		}

		log.Printf("export: created %q\n", globalExport.File)

		return nil
	},
}

func init() {
	cmdExport.Flags().StringVar(&globalExport.File, "filename", "", "file name to create")
	_ = cmdExport.MarkFlagRequired("file")
	cmdExport.Flags().StringVar(&globalExport.Game, "game", "", "game to export")
	_ = cmdExport.MarkFlagRequired("game")
	cmdExport.Flags().StringVar(&globalExport.Turn, "turn", "", "turn to export")

	cmdBase.AddCommand(cmdExport)
}
