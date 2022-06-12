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
	"fmt"
	"github.com/spf13/cobra"
)

var globalCreate struct {
	Force bool
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "create new items",
	Long:  `Create a new thing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return fmt.Errorf("missing config file name")
		}
		return nil
	},
}

var globalCreateGame struct {
	Name      string
	ShortName string
}

var cmdCreateGame = &cobra.Command{
	Use:   "game",
	Short: "create new game",
	Long:  `Create a new game.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return fmt.Errorf("missing config file name")
		}
		return nil
	},
}

func init() {
	cmdCreate.PersistentFlags().BoolVar(&globalCreate.Force, "force", false, "force overwrite of existing data")

	cmdCreateGame.Flags().StringVar(&globalCreateGame.Name, "name", "", "full name of new game")
	_ = cmdCreate.MarkFlagRequired("name")
	cmdCreateGame.Flags().StringVar(&globalCreateGame.ShortName, "short-name", "", "short name of new game")
	_ = cmdCreate.MarkFlagRequired("short-name")

	cmdBase.AddCommand(cmdCreate)
	cmdCreate.AddCommand(cmdCreateGame)
}
