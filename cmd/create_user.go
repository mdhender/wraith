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
	"github.com/google/uuid"
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"unicode"
)

var globalCreateUser struct {
	Id     string
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
			return errors.New("missing user handle")
		}
		for _, r := range globalCreateUser.Handle {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
				return errors.New("invalid rune in user handle")
			}
		}

		// validate the new user's id or supply default value if missing
		globalCreateUser.Id = strings.TrimSpace(globalCreateUser.Id)
		if globalCreateUser.Id == "" {
			globalCreateUser.Id = uuid.New().String()
		}

		// load the base configuration to find the users store
		cfgBase, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfgBase.Self)

		// load the users store
		cfgUsers, err := engine.LoadUsers(cfgBase.Store)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded users store %q\n", cfgUsers.Store)

		// check for duplicate handle or id
		for _, u := range cfgUsers.Index {
			if strings.ToLower(u.Handle) == strings.ToLower(globalCreateUser.Handle) {
				log.Fatalf("duplicate handle %q", globalCreateUser.Handle)
			} else if strings.ToLower(u.Id) == strings.ToLower(globalCreateUser.Id) {
				log.Fatalf("duplicate id %q", globalCreateUser.Id)
			}
		}

		// add the new user to the users store
		cfgUsers.Index = append(cfgUsers.Index, engine.UsersIndex{
			Id:     globalCreateUser.Id,
			Handle: globalCreateUser.Handle,
		})

		log.Printf("updating users store %q\n", cfgUsers.Store)
		if err := cfgUsers.Write(); err != nil {
			log.Fatal(err)
		}

		log.Printf("updated users store %q\n", cfgUsers.Store)
		return nil
	},
}

func init() {
	cmdCreateUser.Flags().StringVar(&globalCreateUser.Handle, "handle", "", "screen name of the new user")
	_ = cmdCreateUser.MarkFlagRequired("handle")
	cmdCreateUser.Flags().StringVar(&globalCreateUser.Id, "id", "", "identifier for the new user (do not use)")

	cmdCreate.AddCommand(cmdCreateUser)
}
