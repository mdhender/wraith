////////////////////////////////////////////////////////////////////////////////
// wraith - the wraith game engine and Server
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

package cheese

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/internal/adapters"
	"github.com/mdhender/wraith/storage/jdb"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (s *Server) serve() error {
	// For testing purposes, we hardcode a JWT token with claims here
	tokenAuth := jwtauth.New("HS256", s.key, nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": "mdhender"})
	tokenCookie := "jwt"

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/_auth", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("%s: %s", r.Method, r.URL.Path)))
	})

	// public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		})
		r.Get("/jwt/cookie/clear", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     tokenCookie,
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			_, _ = w.Write([]byte(fmt.Sprintf("cookie: clear %q: ok", tokenCookie)))
		})
		r.Get("/jwt/cookie/get", func(w http.ResponseWriter, r *http.Request) {
			if c, err := r.Cookie(tokenCookie); err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("cookie: get %q: %+v", tokenCookie, err)))
			} else {
				_, _ = w.Write([]byte(c.Value))
			}
		})
		r.Get("/jwt/cookie/set", func(w http.ResponseWriter, r *http.Request) {
			maxAge := 14 * 24 * 60 * 60
			http.SetCookie(w, &http.Cookie{
				Name:     tokenCookie,
				Path:     "/",
				Value:    tokenString,
				MaxAge:   maxAge,
				HttpOnly: true,
			})
			_, _ = w.Write([]byte(fmt.Sprintf("cookie: set %q: %q", tokenCookie, tokenString)))
		})
		r.Get("/jwt/token/get/{user_id}", func(w http.ResponseWriter, r *http.Request) {
			claims := map[string]interface{}{"user_id": chi.URLParam(r, "user_id")}
			jwtauth.SetExpiryIn(claims, time.Second*60*60)
			_, tokenString, _ := tokenAuth.Encode(claims)
			_, _ = w.Write([]byte(tokenString))
		})

		r.Get("/ui/login", s.loginGetHandler(s.templates, tokenCookie, tokenString))
		r.Post("/ui/login", s.loginPostHandler(tokenCookie, tokenString))
		//r.Get("/ui/login/{handle}/{secret}", s.loginGetHandleSecretHandler(tokenCookie, tokenString))

		r.Get("/ui/logout", s.logoutHandler(tokenCookie))
		r.Post("/ui/logout", s.logoutHandler(tokenCookie))
		r.Put("/ui/logout", s.logoutHandler(tokenCookie))
	})

	// protected routes
	r.Group(func(r chi.Router) {
		// pull, verify, and validate JWT tokens from cookie or bearer token
		//r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(s.myVerifier(tokenAuth))

		// Handle valid / invalid tokens.
		// In this example, we use the provided authenticator middleware, but you can write your own very easily.
		// Look at the Authenticator method in jwtauth.go and tweak it; it's not scary.
		r.Use(jwtauth.Authenticator)

		r.Route("/api", func(r chi.Router) {
			r.Get("/claims", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("<body><code><pre>"))
				_, _ = w.Write([]byte(fmt.Sprintf("%s: %s:\n", r.Method, r.URL.Path)))
				_, claims, _ := jwtauth.FromContext(r.Context())
				userId, ok := claims["user_id"].(string)
				if !ok {
					log.Printf("%s: %s: claims[%q] is not a string\n", "user_id")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				claim, ok := s.claims[strings.ToLower(userId)]
				if !ok {
					log.Printf("%s: %s: fetchClaims: %q: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				_, _ = w.Write([]byte(fmt.Sprintf("claim.User %q\n", claim.User)))
				_, _ = w.Write([]byte(fmt.Sprintf("claim.NationNo %d\n", claim.NationNo)))
				_, _ = w.Write([]byte(fmt.Sprintf("claims.Player %q\n", claim.PlayerName)))
				_, _ = w.Write([]byte("</pre></code></body>"))
			})
			r.Get("/panic", func(http.ResponseWriter, *http.Request) {
				panic("panic")
			})
			r.Get("/report/game/{game}/nation/{nation}/turn/{year}/{quarter}", func(w http.ResponseWriter, r *http.Request) {
				game, nation, year, quarter := chi.URLParam(r, "game"), chi.URLParam(r, "nation"), chi.URLParam(r, "year"), chi.URLParam(r, "quarter")
				_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: game %q nation %q year %q quarter %q", r.Method, r.URL.Path, game, nation, year, quarter)))
			})
		})

		r.Route("/ui", func(r chi.Router) {
			r.Get("/", s.homeGetHandler(s.templates))
			r.Get("/games/{game}/cluster", s.clusterGetHandler(s.templates))
			r.Get("/games/{game}/cluster/{x}/{y}/{z}", s.clusterGetHandler(s.templates))
			r.Get("/games/{game}/current-report", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				claim, ok := s.claims[strings.ToLower(claims["user_id"].(string))]
				if !ok {
					log.Printf("%s: %s: fetchClaims: not ok\n", r.Method, r.URL.Path)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				game := chi.URLParam(r, "game")
				year := 0
				quarter := 0
				e, err := engine.Open(s.store)
				if err != nil {
					log.Printf("%s: %s: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				bw := bytes.NewBuffer([]byte(fmt.Sprintf("<body><h1>Nation %d</h1><code><pre>", claim.NationNo)))
				err = e.Report(bw, game, claim.NationNo, year, quarter)
				if err != nil {
					log.Printf("%s: %s: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				bw.Write([]byte("</pre></code></body>"))
				_, _ = w.Write(bw.Bytes())
			})
			r.Get("/games/{game}/nations/{nation}/turn/{year}/{quarter}/report", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				claim, ok := s.claims[strings.ToLower(claims["user_id"].(string))]
				if !ok {
					log.Printf("%s: %s: fetchClaims: not ok\n", r.Method, r.URL.Path)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				pNation := chi.URLParam(r, "nation")
				nationNo, err := strconv.Atoi(pNation)
				if err != nil {
					log.Printf("%s: %s: nation: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				} else if nationNo != claim.NationNo {
					log.Printf("%s: %s: nation: claim.NationNo %d: nationNo %d\n", r.Method, r.URL.Path, claim.NationNo, nationNo)
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}

				game, pYear, pQuarter := chi.URLParam(r, "game"), chi.URLParam(r, "year"), chi.URLParam(r, "quarter")
				year, err := strconv.Atoi(pYear)
				if err != nil {
					log.Printf("%s: %s: year: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				quarter, err := strconv.Atoi(pQuarter)
				if err != nil {
					log.Printf("%s: %s: quarter: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				e, err := engine.Open(s.store)
				if err != nil {
					log.Printf("%s: %s: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				bw := bytes.NewBuffer([]byte(fmt.Sprintf("<body><h1>Nation %d</h1><code><pre>", nationNo)))
				err = e.Report(bw, game, nationNo, year, quarter)
				if err != nil {
					log.Printf("%s: %s: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				bw.Write([]byte("</pre></code></body>"))
				_, _ = w.Write(bw.Bytes())
			})
			r.Get("/reports/{game}/{year}/{quarter}/{player}", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				claim, ok := s.claims[strings.ToLower(claims["user_id"].(string))]
				if !ok {
					log.Printf("%s: %s: fetchClaims: not ok\n", r.Method, r.URL.Path)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				pPlayer := chi.URLParam(r, "player")
				playerId, err := strconv.Atoi(pPlayer)
				if err != nil {
					log.Printf("%s: %s: player: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				} else if playerId != claim.PlayerId {
					log.Printf("%s: %s: player: claim.PlayerName %q: claim.PlayerId %d: playerId %d\n", r.Method, r.URL.Path, claim.PlayerName, claim.PlayerId, playerId)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}

				game, pYear, pQuarter := chi.URLParam(r, "game"), chi.URLParam(r, "year"), chi.URLParam(r, "quarter")
				year, err := strconv.Atoi(pYear)
				if err != nil {
					log.Printf("%s: %s: year: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				quarter, err := strconv.Atoi(pQuarter)
				if err != nil {
					log.Printf("%s: %s: quarter: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				gamePath := filepath.Clean(filepath.Join("D:\\wraith\\testdata\\games", game, fmt.Sprintf("%04d", year), fmt.Sprintf("%d", quarter)))
				log.Printf("hey: gamePath %s\n", gamePath)
				jg, err := jdb.Load(filepath.Join(gamePath, "game.json"))
				if err != nil {
					log.Printf("%s: %s: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				e, err := adapters.JdbGameToWraithEngine(jg)
				if err != nil {
					log.Printf("%s: %s: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				bw := bytes.NewBuffer([]byte(fmt.Sprintf("<body><h1>Player %d</h1><code><pre>", playerId)))
				err = e.Report(bw, playerId)
				if err != nil {
					log.Printf("%s: %s: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				bw.Write([]byte("</pre></code></body>"))
				_, _ = w.Write(bw.Bytes())
			})
			r.Get("/games/{game}/orders", s.ordersGetRedirect())
			r.Get("/games/{game}/orders/{year}/{quarter}", s.ordersGetHandler(s.templates))
			r.Post("/games/{game}/orders/{year}/{quarter}", s.ordersPostHandler())
			r.Get("/units", s.unitsGetHandler(s.templates))
		})
	})

	log.Printf("Server: listening on %s\n", s.addr)
	return http.ListenAndServe(s.addr, r)
}
