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
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"unicode"
)

var globalAddUser struct {
	Id     string
	Handle string
}

var cmdAddUser = &cobra.Command{
	Use:   "user",
	Short: "add a new user",
	Long:  `Add a new user to the engine.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		// validate the new user's handle
		globalAddUser.Handle = strings.TrimSpace(globalAddUser.Handle)
		if globalAddUser.Handle == "" {
			return errors.New("missing user handle")
		}
		for _, r := range globalAddUser.Handle {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
				return errors.New("invalid rune in user handle")
			}
		}

		// validate the new user's id or supply default value if missing
		globalAddUser.Id = strings.TrimSpace(globalAddUser.Id)
		if globalAddUser.Id == "" {
			globalAddUser.Id = uuid.New().String()
		}

		// load the base configuration to find the users store
		cfgBase, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfgBase.Path)

		// load the users store
		cfgUsers, err := config.LoadUsers(cfgBase.UsersStore)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded users store %q\n", cfgUsers.Path)

		// check for duplicate handle or id
		for _, u := range cfgUsers.Users {
			if strings.ToLower(u.Handle) == strings.ToLower(globalAddUser.Handle) {
				log.Fatalf("duplicate handle %q", globalAddUser.Handle)
			} else if strings.ToLower(u.Id) == strings.ToLower(globalAddUser.Id) {
				log.Fatalf("duplicate id %q", globalAddUser.Id)
			}
		}

		// add the new user to the users store
		cfgUsers.Users = append(cfgUsers.Users, config.User{
			Id:     globalAddUser.Id,
			Handle: globalAddUser.Handle,
		})

		log.Printf("updating users store %q\n", cfgUsers.Path)
		if err := cfgUsers.Write(); err != nil {
			log.Fatal(err)
		}

		log.Printf("updated users store %q\n", cfgUsers.Path)
		return nil
	},
}

func init() {
	cmdAddUser.Flags().StringVar(&globalAddUser.Handle, "handle", "", "screen name of the new user")
	_ = cmdAddUser.MarkFlagRequired("handle")
	cmdAddUser.Flags().StringVar(&globalAddUser.Id, "id", "", "identifier for the new user (do not use)")

	cmdAdd.AddCommand(cmdAddUser)
}
