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
	"context"
	"encoding/json"
	"errors"
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var globalCreateUsers struct {
	Filename string
}

var cmdCreateUsers = &cobra.Command{
	Use:   "users",
	Short: "create new users from a data file",
	Long:  `Create new users from a JSON data file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		globalCreateUsers.Filename = strings.TrimSpace(globalCreateUsers.Filename)
		if globalCreateUsers.Filename == "" {
			return errors.New("missing data file name")
		}

		cfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded config %q\n", cfg.Self)

		data, err := os.ReadFile(globalCreateUsers.Filename)
		if err != nil {
			log.Fatal(err)
		}

		e, err := engine.Open(cfg, context.Background())
		if err != nil {
			log.Fatal(err)
		}

		var users []struct {
			Handle string `json:"handle"`
			Email  string `json:"email"`
			Secret string `json:"secret"`
		}
		err = json.Unmarshal(data, &users)
		if err != nil {
			log.Fatal(err)
		}

		for _, user := range users {
			user.Handle = strings.ToLower(strings.TrimSpace(user.Handle))
			user.Email = strings.ToLower(strings.TrimSpace(user.Email))

			err := e.CreateUser(user.Handle, user.Email, user.Secret)
			if err != nil {
				log.Printf("user %q %q: %+v\n", user.Handle, user.Email, err)
			} else {
				log.Printf("user %q %q: created\n", user.Handle, user.Email)
			}
		}

		return nil
	},
}

func init() {
	cmdCreateUsers.Flags().StringVar(&globalCreateUsers.Filename, "data", "", "name of json data file to load")
	_ = cmdCreateUsers.MarkFlagRequired("data")

	cmdCreate.AddCommand(cmdCreateUsers)
}
