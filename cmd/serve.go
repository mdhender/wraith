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
	"context"
	"github.com/mdhender/jsonwt"
	"github.com/mdhender/jsonwt/signers"
	"github.com/mdhender/wraith/internal/config"
	"github.com/mdhender/wraith/internal/otohttp"
	gsvc "github.com/mdhender/wraith/internal/services/greeter"
	isvc "github.com/mdhender/wraith/internal/services/identity"
	vsvc "github.com/mdhender/wraith/internal/services/version"
	"github.com/mdhender/wraith/internal/storage/identity"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var globalServe struct {
	PIDFile bool // create a file containing the PID if set
	pid     int  // current PID
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
			log.Printf("[serve] %-30s == %q\n", "identity.config", cfg.Identity.Repository.JSONFile)
		}

		if globalServe.PIDFile {
			globalServe.pid = os.Getpid()
			//if err := ioutil.WriteFile("/tmp/.fhapp.pid", []byte(fmt.Sprintf("%d", pid)), 0600); err != nil {
			//	log.Printf("unable to create pid file: %+v", err)
			//	os.Exit(2)
			//}
			log.Printf("server: pid %8d: file %q\n", globalServe.pid, "/tmp/.wraith.pid")
		}

		// load the identity store
		hs256, err := signers.NewHS256([]byte(cfg.Secrets.Signing))
		if err != nil {
			log.Fatal(err)
		}
		i, err := identity.Load(cfg.Identity.Repository.JSONFile, jsonwt.NewFactory("る", hs256))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[serve] %-30s == %q\n", "identity.factory", i.Version)

		// start the server with the ability to shut it down gracefully
		// thanks to https://clavinjune.dev/en/blogs/golang-http-server-graceful-shutdown/.
		// TODO: this should be part of the server.Server implementation!
		log.Printf("server: todo: please move the shutdown logic to the server implementation!\n")

		// create a context that we can use to cancel the server
		ctx, cancel := context.WithCancel(context.Background())

		options := []otohttp.Option{
			otohttp.WithContext(ctx),
			otohttp.WithAddr("", "8080"),
		}
		s, err := otohttp.NewServer(otohttp.Options(options...))
		if err != nil {
			log.Fatal(err)
		}
		gsvc.RegisterGreeterService(s, gsvc.Service{})
		isvc.RegisterIdentityService(s, isvc.Service{})
		vsvc.RegisterVersionService(s, vsvc.Service{})

		// run server in a go routine that we can cancel
		go func() {
			log.Printf("server: listening on %q\n", net.JoinHostPort(cfg.Server.Host, cfg.Server.Port))
			err := http.ListenAndServe(net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), s)
			if err != http.ErrServerClosed {
				log.Fatalf("server: %v", err)
			}
		}()

		// catch signals to interrupt the server and shut it down
		chanSignal := make(chan os.Signal, 1)
		signal.Notify(chanSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-chanSignal
		if globalServe.PIDFile {
			log.Printf("server: signal: interrupt: shutting down pid %d...\n", globalServe.pid)
		} else {
			log.Print("server: signal: interrupt: shutting down...\n")
		}
		go func() {
			<-chanSignal // in case the user is spraying us with interrupts...
			log.Fatal("server: signal: kill: terminating...\n")
		}()

		// allow 5 seconds for a graceful shutdown
		ctxWithDelay, cancelNow := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelNow()
		if err := s.Shutdown(ctxWithDelay); err != nil {
			panic(err) // if that fails, panic!
		}
		log.Printf("server: stopped\n")

		// manually cancel context if not using httpServer.RegisterOnShutdown(cancel)
		cancel()

		defer os.Exit(0)
		return nil
	},
}

func init() {
	cmdServe.Flags().BoolVar(&globalServe.PIDFile, "pid-file", false, "create pid file on startup")

	cmdBase.AddCommand(cmdServe)
}
