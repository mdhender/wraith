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
	"github.com/mdhender/wraith/internal/server"
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
)

var globalServer struct {
	Host      string
	Port      string
	AuthnFile string
	JwtFile   string
}

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "test server",
	Long:  `Create a web server to test the engine.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalBase.ConfigFile == "" {
			return errors.New("missing config file name")
		}

		cfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfg.Self)

		var options []server.Option
		options = append(options, server.WithHost(globalServer.Host))
		options = append(options, server.WithPort(globalServer.Port))
		options = append(options, server.WithAuthenticationData(globalServer.AuthnFile))
		options = append(options, server.WithJwtData(globalServer.JwtFile))
		s, err := server.New(options...)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("server: listening on %s\n", net.JoinHostPort(globalServer.Host, globalServer.Port))
		if http.ListenAndServe(net.JoinHostPort(globalServer.Host, globalServer.Port), s.Mux()) != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	cmdServer.Flags().StringVar(&globalServer.Host, "host", "", "host interface to listen on")
	cmdServer.Flags().StringVar(&globalServer.Port, "port", "8123", "port to listen on")
	cmdServer.Flags().StringVar(&globalServer.AuthnFile, "authn", "", "authentication data")
	_ = cmdServer.MarkFlagRequired("authn")
	cmdServer.Flags().StringVar(&globalServer.JwtFile, "jwt", "", "jwt key data")
	_ = cmdServer.MarkFlagRequired("jwt")

	cmdBase.AddCommand(cmdServer)
}
