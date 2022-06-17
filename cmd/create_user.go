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
	"github.com/mdhender/wraith/models"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

var globalCreateUser struct {
	Email  string
	Handle string
	Secret string
}

var cmdCreateUser = &cobra.Command{
	Use:   "user",
	Short: "create a new user",
	Long:  `Create a new user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		// validate the new user's information
		globalCreateUser.Email = strings.ToLower(strings.TrimSpace(globalCreateUser.Email))
		if globalCreateUser.Email == "" {
			return errors.New("missing email")
		}
		globalCreateUser.Handle = strings.ToLower(strings.TrimSpace(globalCreateUser.Handle))
		if globalCreateUser.Handle == "" {
			return errors.New("missing handle")
		}
		if globalCreateUser.Secret == "" {
			return errors.New("missing secret")
		} else if len(globalCreateUser.Secret) < 8 {
			return errors.New("secret too short")
		} else {
			hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(globalCreateUser.Secret), bcrypt.MinCost)
			if err != nil {
				log.Fatal(err)
			}
			globalCreateUser.Secret = string(hashedPasswordBytes)
		}

		cfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfg.Self)

		s, err := models.Open(cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()

		u, err := s.CreateUser(models.User{
			Email:        globalCreateUser.Email,
			Handle:       globalCreateUser.Handle,
			HashedSecret: globalCreateUser.Secret,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("user %v\n", u)

		return nil
	},
}

func init() {
	cmdCreateUser.Flags().StringVar(&globalCreateUser.Email, "email", "", "e-mail account for the new user")
	_ = cmdCreateUser.MarkFlagRequired("email")
	cmdCreateUser.Flags().StringVar(&globalCreateUser.Handle, "handle", "", "screen name of the new user")
	_ = cmdCreateUser.MarkFlagRequired("handle")
	cmdCreateUser.Flags().StringVar(&globalCreateUser.Secret, "secret", "", "secret password for the new user")
	_ = cmdCreateUser.MarkFlagRequired("secret")

	cmdCreate.AddCommand(cmdCreateUser)
}
