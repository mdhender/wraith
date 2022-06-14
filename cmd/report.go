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
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var globalReport struct {
	Game string
}

var cmdReport = &cobra.Command{
	Use:   "report",
	Short: "create status reports",
	Long:  `Create status reports.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		// validate the game name
		globalReport.Game = strings.TrimSpace(globalReport.Game)
		if globalReport.Game == "" {
			return errors.New("missing game name")
		}

		e, err := engine.LoadGame(globalBase.ConfigFile, globalReport.Game)
		if err != nil {
			log.Fatal(err)
		}

		err = e.Report("SP1")
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	cmdReport.PersistentFlags().StringVar(&globalReport.Game, "game", "", "name of game to report on")
	_ = cmdReport.MarkFlagRequired("game")

	cmdBase.AddCommand(cmdReport)
}
