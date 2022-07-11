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
	"github.com/mdhender/wraith/internal/orders"
	"github.com/mdhender/wraith/models"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
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

		ordersFile := filepath.Join(cfg.OrdersPath, fmt.Sprintf("%s.%s.%s.%d.txt", globalRun.Game, "0000", "0", 1))
		b, err := os.ReadFile(ordersFile)
		if err != nil {
			log.Fatal(err)
		}

		p, err := orders.Parse([]byte(b))
		if err != nil {
			log.Fatal(err)
		}

		for _, order := range p {
			if assembleFactoryGroup, ok := order.Command.(*orders.AssembleFactoryGroup); ok {
				log.Printf("line %d: %q %q %d %q\n", order.Line, "assemble-factory-group", assembleFactoryGroup.Id, assembleFactoryGroup.Qty, assembleFactoryGroup.Product)
			} else if assembleMineGroup, ok := order.Command.(*orders.AssembleMineGroup); ok {
				log.Printf("line %d: %q %q %d %q\n", order.Line, "assemble-mine-group", assembleMineGroup.Id, assembleMineGroup.Qty, assembleMineGroup.DepositId)
			} else if name, ok := order.Command.(*orders.NameOrder); ok {
				if name.Id[0] == 'C' {
					log.Printf("line %d: %q %q %q\n", order.Line, "name-colony", name.Id, name.Name)
				} else if name.Id[0] == 'S' {
					log.Printf("line %d: %q %q %q\n", order.Line, "name-ship", name.Id, name.Name)
				}
			}
		}

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
