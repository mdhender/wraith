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
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/internal/adapters"
	"github.com/mdhender/wraith/internal/orders"
	"github.com/mdhender/wraith/models"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"time"
)

var globalRun struct {
	Game  string
	Phase string
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "run a phase",
	Long:  `Run a phase of the game.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
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

		game, err := s.LookupGameByName(globalRun.Game)
		if err != nil {
			log.Fatal(err)
		} else if game, err = s.FetchGameByNameAsOf(game.ShortName, game.CurrentTurn.String()); err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded game %q: turn %q\n", game.ShortName, game.CurrentTurn.String())

		e, err := engine.Open(s)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded engine version %q\n", e.Version())

		now := time.Now()
		users, err := s.FetchUserClaimsFromGameAsOf(game.Id, now)
		if err != nil {
			log.Fatal(err)
		}
		var pos []*engine.PlayerOrders
		for _, user := range users {
			ordersFile := filepath.Join(cfg.OrdersPath, fmt.Sprintf("%s.%04d.%d.%d.txt", user.Games[0].ShortName, user.Games[0].EffTurn.Year, user.Games[0].EffTurn.Quarter, user.Games[0].PlayerId))
			b, err := os.ReadFile(ordersFile)
			if err != nil {
				log.Printf("run: nation %3d handle %-32q pid %4d turn %q: %+v\n", user.Games[0].NationNo, user.Games[0].PlayerHandle, user.Games[0].PlayerId, user.Games[0].EffTurn, err)
				continue
			}

			o, err := orders.Parse([]byte(b))
			if err != nil {
				log.Printf("run: nation %3d handle %-32q pid %4d turn %q: %+v\n", user.Games[0].NationNo, user.Games[0].PlayerHandle, user.Games[0].PlayerId, user.Games[0].EffTurn, err)
				continue
			}

			log.Printf("run: nation %3d handle %-32q pid %4d turn %q: orders:\n%s\n", user.Games[0].NationNo, user.Games[0].PlayerHandle, user.Games[0].PlayerId, user.Games[0].EffTurn, "...")
			pos = append(pos, e.PlayerOrders(adapters.ModelsPlayerToEnginePlayer(game.Players[user.Games[0].PlayerId]), o))
		}

		err = e.Execute(pos, "control", "retool")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("wow. executed!\n")

		return nil
	},
}

func init() {
	cmdRun.Flags().StringVar(&globalRun.Game, "game", "", "game to run against")
	_ = cmdRun.MarkFlagRequired("game")
	cmdRun.Flags().StringVar(&globalRun.Phase, "phase", "", "phase to process")
	_ = cmdRun.MarkFlagRequired("phase")

	cmdBase.AddCommand(cmdRun)
}
