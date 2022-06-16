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

var globalCreateUser struct {
	Handle string
}

var cmdCreateUser = &cobra.Command{
	Use:   "user",
	Short: "create a new user",
	Long:  `Create a new user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		// validate the new user's handle
		globalCreateUser.Handle = strings.TrimSpace(globalCreateUser.Handle)
		if globalCreateUser.Handle == "" {
			return errors.New("missing handle")
		}

		// load the engine
		e, err := engine.New(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}

		err = e.AddUser(globalCreateUser.Handle)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	cmdCreateUser.Flags().StringVar(&globalCreateUser.Handle, "handle", "", "screen name of the new user")
	_ = cmdCreateUser.MarkFlagRequired("handle")

	cmdCreate.AddCommand(cmdCreateUser)
}
