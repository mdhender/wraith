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
	"github.com/mdhender/jsonwt"
	"github.com/mdhender/jsonwt/signers"
	"github.com/mdhender/wraith/internal/config"
	"github.com/mdhender/wraith/internal/storage/identity"
	"github.com/mdhender/wraith/internal/storage/words"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"path/filepath"
)

var globalBootstrap struct {
	ConfigFile string
	BCrypt     struct {
		Cost int // bcrypt hashing cost
	}
	Force    bool // overwrite any existing configuration only if set
	Identity struct {
		Filename string // path to create identity store
	}
	Secrets struct {
		Signing string // plain-text
		Sysop   string // plain-text
	}
	Server struct {
		Host string
		Port string
	}
}

var cmdBootstrap = &cobra.Command{
	Use:   "bootstrap",
	Short: "create a new global configuration file",
	Long: `Create the initial system configuration.
This includes the configuration file, sysop account, and starting data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		psg, err := words.New(" ")
		if err != nil {
			log.Fatal(err)
		}

		globalBootstrap.ConfigFile = globalBase.ConfigFile
		if globalBootstrap.ConfigFile == "" {
			return fmt.Errorf("missing config file name")
		}
		log.Printf("creating config store %q\n", globalBootstrap.ConfigFile)

		if globalBootstrap.Identity.Filename == "" {
			return fmt.Errorf("missing identity file name")
		}

		if len(globalBootstrap.Secrets.Sysop) == 0 {
			globalBootstrap.Secrets.Sysop = psg.Generate(5)
		}
		if len(globalBootstrap.Secrets.Sysop) < 12 {
			return fmt.Errorf("sysop-password must be at least 12 characters long")
		}

		if len(globalBootstrap.Secrets.Signing) == 0 {
			globalBootstrap.Secrets.Signing = psg.Generate(12)
		}
		if len(globalBootstrap.Secrets.Signing) < 32 {
			return fmt.Errorf("signing-key must be at least 32 characters long")
		}
		cfg := config.Config{
			ConfigFile: globalBootstrap.ConfigFile,
			//Accounts: make(map[string]config.Account),
			//Users: make(map[string]config.User),
		}

		cfg.Identity.Repository.JSONFile = filepath.Clean(globalBootstrap.Identity.Filename)
		cfg.Identity.Cost = globalBootstrap.BCrypt.Cost
		if cfg.Identity.Cost < bcrypt.DefaultCost {
			cfg.Identity.Cost = bcrypt.DefaultCost
		}
		if cfg.Identity.Cost < bcrypt.MinCost {
			cfg.Identity.Cost = bcrypt.MinCost
		}

		cfg.Server.Host = globalBootstrap.Server.Host
		cfg.Server.Port = globalBootstrap.Server.Port

		cfg.Secrets.Signing = globalBootstrap.Secrets.Signing
		cfg.Secrets.Sysop = globalBootstrap.Secrets.Sysop

		// bootstrap the identity store data file
		hs256, err := signers.NewHS256([]byte(cfg.Secrets.Signing))
		if err != nil {
			log.Fatal(err)
		}
		i, err := identity.Bootstrap(cfg.Identity.Repository.JSONFile, cfg.Identity.Cost, jsonwt.NewFactory("る", hs256))
		if err != nil {
			log.Fatal(err)
		}

		// create the sysop account used by the command line interface
		if err := i.Create(identity.Identity{
			Email:  "sysop",
			Handle: "sysop",
			Secret: cfg.Secrets.Sysop,
			Roles:  []string{"sysop"},
		}); err != nil {
			log.Fatal(err)
		}
		sysops := i.Fetch(func(u identity.Identity) bool {
			return u.Handle == "sysop"
		})
		if len(sysops) != 1 {
			log.Fatal("failed to create sysop in identity store")
		}
		sysop := sysops[0]
		log.Printf("sysop %+v\n", sysop)
		log.Printf("created identity store %q\n", i.Filename)

		if _, err := os.Stat(globalBootstrap.ConfigFile); err == nil {
			log.Fatal("cowardly refusing to overwrite existing configuration store")
		}
		if err := config.Write(globalBootstrap.ConfigFile, &cfg); err != nil {
			return err
		}

		log.Printf("***** these values will stored as plain-text in the config store.\nsysop passphrase: %q\nsigning key: %q\n",
			globalBootstrap.Secrets.Sysop, globalBootstrap.Secrets.Signing)
		log.Printf("***** please protect the config file.\n")

		log.Printf("created config store %q\n", globalBootstrap.ConfigFile)
		return nil
	},
}

func init() {
	cmdBootstrap.Flags().IntVar(&globalBootstrap.BCrypt.Cost, "bcrypt-cost", bcrypt.DefaultCost, "bcrypt hash setting")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Secrets.Signing, "signing-key", "", "key for signing tokens")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Secrets.Sysop, "sysop-secret", "", "passphrase for sysop account")
	_ = cmdBootstrap.MarkFlagRequired("sysop-password")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Identity.Filename, "identity-file", "", "path to create identity store")
	_ = cmdBootstrap.MarkFlagRequired("identity-file")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Server.Host, "host", "", "binding for server")
	cmdBootstrap.Flags().StringVar(&globalBootstrap.Server.Port, "port", "8080", "port for server")

	cmdBase.AddCommand(cmdBootstrap)
}
