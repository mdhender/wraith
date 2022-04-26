/*
 * Wraith Game Engine
 * Copyright (c) 2022 Michael D. Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package cmd

import (
	"github.com/mdhender/wraith/internal/config"
	"github.com/mdhender/wraith/internal/server"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
)

var globalServe struct {
}

var cmdServe = &cobra.Command{
	Use:   "serve",
	Short: "start the API server",
	Long:  `Start the API server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{ConfigFile: globalBase.ConfigFile}
		if err := config.Read(cfg); err != nil {
			log.Fatal(err)
		}
		if globalBase.VerboseFlag {
			log.Printf("[serve] %-30s == %q\n", "config", cfg.ConfigFile)
			log.Printf("[serve] %-30s == %q\n", "host", cfg.Server.Host)
			log.Printf("[serve] %-30s == %q\n", "port", cfg.Server.Port)
		}
		s, err := server.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("listening on %q\n", net.JoinHostPort(cfg.Server.Host, cfg.Server.Port))
		return http.ListenAndServe(net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), s)
	},
}

func init() {
	cmdBase.AddCommand(cmdServe)
}
