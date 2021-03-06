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
	"github.com/mdhender/wraith/internal/cheese"
	"github.com/mdhender/wraith/models"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var globalServer struct {
	Host      string
	Port      string
	JwtFile   string
	JwtKey    string
	Templates string
}

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "test server",
	Long:  `Create a web server to test the engine.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		if globalServer.JwtKey = strings.TrimSpace(globalServer.JwtKey); globalServer.JwtKey == "" {
			return errors.New("missing jwt signing key")
		}

		if globalServer.Templates = strings.TrimSpace(globalServer.Templates); globalServer.Templates == "" {
			globalServer.Templates = "."
		}

		cfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfg.Self)
		if cfg.GamesPath == "" {
			return errors.New("missing global config games path")
		}

		s, err := models.Open(cfg)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded store version %q\n", s.Version())

		opts := []cheese.Option{
			cheese.WithHost(globalServer.Host),
			cheese.WithPort(globalServer.Port),
			cheese.WithGamesPath(cfg.GamesPath),
			cheese.WithKey([]byte(globalServer.JwtKey)),
			cheese.WithTemplates(globalServer.Templates),
			cheese.WithStore(s),
		}
		if err := cheese.Serve(opts...); err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	cmdServer.Flags().StringVar(&globalServer.Host, "host", "", "host interface to listen on")
	cmdServer.Flags().StringVar(&globalServer.Port, "port", "3000", "port to listen on")
	cmdServer.Flags().StringVar(&globalServer.JwtFile, "jwt", "", "jwt key data")
	_ = cmdServer.MarkFlagRequired("jwt")
	cmdServer.Flags().StringVar(&globalServer.JwtKey, "jwt-key", "", "jwt signing key")
	_ = cmdServer.MarkFlagRequired("jwt-key")
	cmdServer.Flags().StringVar(&globalServer.Templates, "templates", "", "path to web templates")

	cmdBase.AddCommand(cmdServer)
}
