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
	"github.com/mdhender/wraith/storage/config"
	"github.com/spf13/cobra"
	"log"
	"net"
	"strings"
)

var globalServer struct {
	Host      string
	Port      string
	AuthnFile string
	GameFile  string
	JwtFile   string
	JwtKey    string
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

		cfg, err := config.LoadGlobal(globalBase.ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("loaded %q\n", cfg.Self)

		key := []byte(globalServer.JwtKey)
		if err := cheese.Serve(net.JoinHostPort(globalServer.Host, globalServer.Port), key); err != nil {
			log.Fatal(err)
		}

		////var options []server.Option
		////options = append(options, server.WithGame(globalServer.GameFile))
		////options = append(options, server.WithHost(globalServer.Host))
		////options = append(options, server.WithPort(globalServer.Port))
		////options = append(options, server.WithAuthenticationData(globalServer.AuthnFile))
		////options = append(options, server.WithJwtData(globalServer.JwtFile))
		////s, err := server.New(options...)
		////if err != nil {
		////	log.Fatal(err)
		////}
		//
		//// For testing purposes, we hardcode a JWT token with claims here
		//tokenAuth := jwtauth.New("HS256", []byte("secret"), nil) // replace with secret key
		//_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": "mdhender"})
		//tokenCookie := "jwt"
		//
		//r := chi.NewRouter()
		//r.Use(middleware.CleanPath)
		//r.Use(middleware.Logger)
		//r.Use(middleware.Recoverer)
		//r.Use(middleware.Heartbeat("/ping"))
		//
		//r.Get("/_auth", func(w http.ResponseWriter, r *http.Request) {
		//	_, _ = w.Write([]byte(fmt.Sprintf("%s: %s", r.Method, r.URL.Path)))
		//})
		//
		//// public routes
		//r.Group(func(r chi.Router) {
		//	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		//		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		//	})
		//	r.Get("/ui/", func(w http.ResponseWriter, r *http.Request) {
		//		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		//	})
		//	r.Get("/jwt/cookie/clear", func(w http.ResponseWriter, r *http.Request) {
		//		http.SetCookie(w, &http.Cookie{
		//			Name:     tokenCookie,
		//			Path:     "/",
		//			MaxAge:   -1,
		//			HttpOnly: true,
		//		})
		//		_, _ = w.Write([]byte(fmt.Sprintf("cookie: clear %q: ok", tokenCookie)))
		//	})
		//	r.Get("/jwt/cookie/get", func(w http.ResponseWriter, r *http.Request) {
		//		if c, err := r.Cookie(tokenCookie); err != nil {
		//			_, _ = w.Write([]byte(fmt.Sprintf("cookie: get %q: %+v", tokenCookie, err)))
		//		} else {
		//			_, _ = w.Write([]byte(c.Value))
		//		}
		//	})
		//	r.Get("/jwt/cookie/set", func(w http.ResponseWriter, r *http.Request) {
		//		maxAge := 14 * 24 * 60 * 60
		//		http.SetCookie(w, &http.Cookie{
		//			Name:     tokenCookie,
		//			Path:     "/",
		//			Value:    tokenString,
		//			MaxAge:   maxAge,
		//			HttpOnly: true,
		//		})
		//		_, _ = w.Write([]byte(fmt.Sprintf("cookie: set %q: %q", tokenCookie, tokenString)))
		//	})
		//	r.Get("/jwt/token/get/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		//		claims := map[string]interface{}{"user_id": chi.URLParam(r, "user_id")}
		//		jwtauth.SetExpiryIn(claims, time.Second*60*60)
		//		_, tokenString, _ := tokenAuth.Encode(claims)
		//		_, _ = w.Write([]byte(tokenString))
		//	})
		//
		//	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		//		http.SetCookie(w, &http.Cookie{
		//			Name:     tokenCookie,
		//			Path:     "/",
		//			MaxAge:   -1,
		//			HttpOnly: true,
		//		})
		//		w.WriteHeader(http.StatusNoContent)
		//	})
		//})
		//
		//// protected routes
		//r.Group(func(r chi.Router) {
		//	// Seek, verify and validate JWT tokens
		//	r.Use(jwtauth.Verifier(tokenAuth))
		//
		//	// Handle valid / invalid tokens.
		//	// In this example, we use the provided authenticator middleware, but you can write your own very easily.
		//	// Look at the Authenticator method in jwtauth.go and tweak it; it's not scary.
		//	r.Use(jwtauth.Authenticator)
		//
		//	r.Route("/api", func(r chi.Router) {
		//		r.Get("/report/game/{game}/nation/{nation}/turn/{year}/{quarter}", func(w http.ResponseWriter, r *http.Request) {
		//			game, nation, year, quarter := chi.URLParam(r, "game"), chi.URLParam(r, "nation"), chi.URLParam(r, "year"), chi.URLParam(r, "quarter")
		//			_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: game %q nation %q year %q quarter %q", r.Method, r.URL.Path, game, nation, year, quarter)))
		//		})
		//	})
		//
		//	r.Get("/games/{game}/nations/{nation}/turn/{year}/{quarter}/report", func(w http.ResponseWriter, r *http.Request) {
		//		gameParam := chi.URLParam(r, "game")
		//		nationParam := chi.URLParam(r, "nation")
		//		yearParam := chi.URLParam(r, "year")
		//		quarterParam := chi.URLParam(r, "quarter")
		//		_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: game %q nation %q year %q quarter %q", r.Method, r.URL.Path, gameParam, nationParam, yearParam, quarterParam)))
		//	})
		//	r.Get("/security", func(w http.ResponseWriter, r *http.Request) {
		//		_, claims, _ := jwtauth.FromContext(r.Context())
		//		_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: claims[user_id] %q", r.Method, r.URL.Path, claims["user_id"])))
		//	})
		//	r.Get("/panic", func(http.ResponseWriter, *http.Request) {
		//		panic("foo")
		//	})
		//})
		//
		//log.Printf("server: listening on %s\n", net.JoinHostPort(globalServer.Host, globalServer.Port))
		//if http.ListenAndServe(net.JoinHostPort(globalServer.Host, globalServer.Port), r) != nil {
		//	log.Fatal(err)
		//}

		return nil
	},
}

func init() {
	cmdServer.Flags().StringVar(&globalServer.Host, "host", "", "host interface to listen on")
	cmdServer.Flags().StringVar(&globalServer.Port, "port", "3000", "port to listen on")
	cmdServer.Flags().StringVar(&globalServer.AuthnFile, "authn", "", "authentication data")
	_ = cmdServer.MarkFlagRequired("authn")
	cmdServer.Flags().StringVar(&globalServer.GameFile, "game", "", "game data")
	_ = cmdServer.MarkFlagRequired("game")
	cmdServer.Flags().StringVar(&globalServer.JwtFile, "jwt", "", "jwt key data")
	_ = cmdServer.MarkFlagRequired("jwt")
	cmdServer.Flags().StringVar(&globalServer.JwtKey, "jwt-key", "", "jwt signing key")
	_ = cmdServer.MarkFlagRequired("jwt-key")

	cmdBase.AddCommand(cmdServer)
}
