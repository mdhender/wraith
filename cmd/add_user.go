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

		globalAddUser.Handle = strings.TrimSpace(globalAddUser.Handle)
		if globalAddUser.Handle == "" {
			return errors.New("missing user handle")
		}
		// validate handle
		for _, r := range globalAddUser.Handle {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
				return errors.New("invalid rune in user handle")
			}
		}

		globalAddUser.Id = strings.TrimSpace(globalAddUser.Id)
		if globalAddUser.Id == "" {
			globalAddUser.Id = uuid.New().String()
		}

		// load the base configuration to find the users store
		globalCfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", globalCfg.FileName)

		// load the users store
		usersCfg, err := config.LoadUsers(globalCfg.UsersStore)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", usersCfg.FileName)

		// create the user to add
		user := config.User{
			Id:     globalAddUser.Id,
			Handle: globalAddUser.Handle,
		}

		// error on duplicate handle or id
		for _, u := range usersCfg.Users {
			if u.Handle == user.Handle {
				log.Fatalf("duplicate handle %q", user.Handle)
			} else if u.Id == user.Id {
				log.Fatalf("duplicate id %q", user.Id)
			}
		}

		usersCfg.Users = append(usersCfg.Users, user)

		if err := usersCfg.Write(); err != nil {
			log.Fatal(err)
		}

		log.Printf("updated %q\n", usersCfg.FileName)
		return nil
	},
}

func init() {
	cmdAddUser.Flags().StringVar(&globalAddUser.Handle, "handle", "", "screen name of the new user")
	_ = cmdAddUser.MarkFlagRequired("handle")
	cmdAddUser.Flags().StringVar(&globalAddUser.Id, "id", "", "identifier for the new user (do not use)")

	cmdAdd.AddCommand(cmdAddUser)
}
