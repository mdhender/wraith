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
	"github.com/mdhender/wraith/internal/adapters"
	"github.com/mdhender/wraith/internal/orders"
	"github.com/mdhender/wraith/storage/config"
	"github.com/mdhender/wraith/storage/jdb"
	"github.com/mdhender/wraith/wraith"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var globalRun struct {
	Root    string
	Year    int
	Quarter int
	Game    string
	Phases  string
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

		if globalRun.Root = strings.TrimSpace(globalRun.Root); globalRun.Root == "" {
			return errors.New("missing path to game files")
		}
		globalRun.Root = filepath.Clean(globalRun.Root)

		if globalRun.Game = strings.TrimSpace(globalRun.Game); globalRun.Game == "" {
			return errors.New("missing game name")
		} else if filepath.Clean(globalRun.Game) != globalRun.Game {
			return errors.New("invalid game name")
		}

		if !(0 <= globalRun.Year && globalRun.Year <= 9999) {
			return errors.New("invalid year")
		}

		if !(1 <= globalRun.Quarter && globalRun.Quarter <= 4) && !(globalRun.Year == 0 && globalRun.Quarter == 0) {
			return errors.New("invalid quarter")
		}

		jg, err := jdb.Load(filepath.Join(globalRun.Root, globalRun.Game, fmt.Sprintf("%04d", globalRun.Year), fmt.Sprintf("%d", globalRun.Quarter), "game.json"))
		if err != nil {
			log.Fatal(err)
		}

		e, err := adapters.JdbGameToWraithEngine(jg)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded engine version %q\n", e.Version)
		log.Printf("loaded game %s: turn %04d/%d\n", e.Game.Code, e.Game.Turn.Year, e.Game.Turn.Quarter)

		var pos []*wraith.PhaseOrders
		for _, user := range e.Players {
			ordersFile := filepath.Join(filepath.Join(globalRun.Root, globalRun.Game, fmt.Sprintf("%04d", globalRun.Year), fmt.Sprintf("%d", globalRun.Quarter), fmt.Sprintf("%d.orders.txt", user.Id)))

			b, err := os.ReadFile(ordersFile)
			if err != nil {
				log.Printf("run: %+v\n", err)
				continue
			}
			log.Printf("run: load %s\n", ordersFile)

			o, err := orders.Parse([]byte(b))
			if err != nil {
				log.Printf("run: nation %3d handle %-32q pid %4d: %+v\n", user.MemberOf.No, user.Name, user.Id, err)
				continue
			}
			//log.Printf("run: nation %3d handle %-32q pid %4d: orders:\n%s\n", user.MemberOf.No, user.Name, user.Id, "...")

			pos = append(pos, adapters.OrdersToPhaseOrders(e.Players[user.Id], o...))
		}

		phases := strings.Split(globalRun.Phases, ",")
		err = e.Execute(pos, phases...)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("wow. executed!\n")

		// bump the turn
		if e.Game.Turn.Quarter = e.Game.Turn.Quarter + 1; e.Game.Turn.Quarter > 4 {
			e.Game.Turn.Year = e.Game.Turn.Year + 1
			e.Game.Turn.Quarter = 1
		}

		// and save the game
		jg = adapters.WraithEngineToJdbGame(e)
		err = jg.Write(filepath.Join(globalRun.Root, globalRun.Game, fmt.Sprintf("%04d", jg.Turn.Year), fmt.Sprintf("%d", jg.Turn.Quarter), "game.json"))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	cmdRun.Flags().StringVar(&globalRun.Root, "root", "", "path to game files")
	_ = cmdRun.MarkFlagRequired("root")
	cmdRun.Flags().StringVar(&globalRun.Game, "game", "", "game to run against")
	_ = cmdRun.MarkFlagRequired("game")
	cmdRun.Flags().IntVar(&globalRun.Year, "year", 0, "turn year")
	_ = cmdRun.MarkFlagRequired("year")
	cmdRun.Flags().IntVar(&globalRun.Quarter, "quarter", 0, "turn quarter")
	_ = cmdRun.MarkFlagRequired("quarter")
	cmdRun.Flags().StringVar(&globalRun.Phases, "phases", "", "comma separated list of phases to process")
	_ = cmdRun.MarkFlagRequired("phases")

	cmdBase.AddCommand(cmdRun)
}
