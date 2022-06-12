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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var globalAdd struct {
	Force bool
}

var cmdAdd = &cobra.Command{
	Use:   "add",
	Short: "add a new item to a game store",
	Long:  `Add a new thing to an existing game store.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please specify an item to add")
	},
}

func init() {
	cmdAdd.PersistentFlags().BoolVar(&globalAdd.Force, "force", false, "force overwrite of existing data")

	cmdBase.AddCommand(cmdAdd)
}
