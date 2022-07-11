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

package cheese

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/internal/osk"
	"github.com/mdhender/wraith/models"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type server struct {
	addr      string
	debug     bool
	key       []byte
	store     *models.Store
	claims    map[string]*models.Claim
	templates string
}

func Serve(addr string, key []byte, store *models.Store, templates string) error {
	s := server{addr: addr, key: key, store: store, templates: templates}

	// fetch user claims
	log.Printf("cheese.Serve: todo: needs game and date logic\n")
	claims, err := s.store.FetchClaims("0000/0")
	if err != nil {
		return err
	}
	s.claims = claims

	return s.serve()
}

func (s *server) serve() error {
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
				_, _ = w.Write([]byte(fmt.Sprintf("claims.Player %q\n", claim.Player)))
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
			r.Get("/games/{game}/orders", s.ordersGetRedirect())
			r.Get("/games/{game}/orders/{year}/{quarter}", s.ordersGetHandler(s.templates))
			r.Post("/games/{game}/orders/{year}/{quarter}", s.ordersPostHandler())
			r.Get("/units", s.unitsGetHandler(s.templates))
		})
	})

	log.Printf("server: listening on %s\n", s.addr)
	return http.ListenAndServe(s.addr, r)
}

func (s *server) clusterGetHandler(templates string) http.HandlerFunc {
	t := osk.New(templates, "cluster_list.html")
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		if _, ok := s.claims[strings.ToLower(claims["user_id"].(string))]; !ok {
			log.Printf("%s: %s: fetchClusterListByGame: not ok\n", r.Method, r.URL.Path)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// todo: use claim to fetch game

		pGameName := chi.URLParam(r, "game")
		game, err := s.store.LookupGameByName(pGameName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		x, _ := strconv.Atoi(chi.URLParam(r, "x"))
		y, _ := strconv.Atoi(chi.URLParam(r, "y"))
		z, _ := strconv.Atoi(chi.URLParam(r, "z"))

		systems, err := s.store.FetchClusterListByGame(game.Id)
		if err != nil {
			log.Printf("%s: %s: game %d: %v\n", r.Method, r.URL.Path, game.Id, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		scans := make(ClusterList, len(systems), len(systems))
		for i, sys := range systems {
			dx, dy, dz := sys.X-x, sys.Y-y, sys.Z-z
			scans[i] = &ClusterListItem{
				Distance: math.Sqrt(float64(dx*dx + dy*dy + dz*dz)),
				X:        sys.X,
				Y:        sys.Y,
				Z:        sys.Z,
				QtyStars: sys.QtyStars,
			}
		}
		sort.Sort(scans)

		type row struct {
			X, Y, Z, QtyStars int
			Distance          string
			URL               string
		}
		clusterList := struct {
			X, Y, Z int // origin
			Systems []row
		}{
			X:       x,
			Y:       y,
			Z:       z,
			Systems: make([]row, len(scans), len(scans)),
		}
		for i, scan := range scans {
			clusterList.Systems[i] = row{
				X:        scan.X,
				Y:        scan.Y,
				Z:        scan.Z,
				QtyStars: scan.QtyStars,
				Distance: fmt.Sprintf("%.3f ly", scan.Distance),
				URL:      fmt.Sprintf("/ui/games/PT-1/cluster/%d/%d/%d", scan.X, scan.Y, scan.Z),
			}
		}

		t.Handle(w, r, clusterList)
	}
}

func (s *server) homeGetHandler(templates string) http.HandlerFunc {
	t := osk.New(templates, "home.html")
	return func(w http.ResponseWriter, r *http.Request) {
		var claim *models.Claim
		_, claims, _ := jwtauth.FromContext(r.Context())
		if userId, ok := claims["user_id"].(string); !ok {
			log.Printf("%s: %s: claims[%q]: not a string\n", r.Method, r.URL.Path, "user_id")
		} else if claim, ok = s.claims[strings.ToLower(userId)]; !ok {
			log.Printf("%s: %s: claims[%q]: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
		}

		t.Handle(w, r, claim)
	}
}

func (s *server) loginGetHandler(templates, cookieName, token string) http.HandlerFunc {
	t := osk.New(templates, "login.html")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s: %s\n", r.Method, r.URL.Path)
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			MaxAge:   -1, // delete any existing cookie
			HttpOnly: true,
		})
		t.ServeHTTP(w, r)
	}
}

func (s *server) loginGetHandleSecretHandler(cookieName string, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("server: %s: %s: %q %q\n", r.Method, r.URL.Path, chi.URLParam(r, "handle"), chi.URLParam(r, "secret"))
		u, err := s.store.FetchUserByCredentials(chi.URLParam(r, "handle"), chi.URLParam(r, "secret"))
		if err != nil {
			log.Printf("server: %s: %s: fetchUsersByCredentials: %v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		log.Printf("server: %s: %s: fetchUsersByCredentials: %q\n", r.Method, r.URL.Path, u.Handle)

		maxAge := 14 * 24 * 60 * 60
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			Value:    token,
			MaxAge:   maxAge,
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: handle %q handle %q email %q", r.Method, r.URL.Path, u.Handle, u.Profiles[0].Handle, u.Profiles[0].Email)))
		//_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: claims[user_id] %q", r.Method, r.URL.Path, claims["user_id"])))
	}
}

func (s *server) loginPostHandler(cookieName string, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("server: %s: %s\n", r.Method, r.URL.Path)

		var input struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		}

		contentType := r.Header.Get("Content-type")
		switch contentType {
		case "application/json":
			r.Body = http.MaxBytesReader(w, r.Body, 1024) // enforce a maximum read of 1kb from the response body
			dec := json.NewDecoder(r.Body)                // create a json decoder that will accept only our specific fields
			dec.DisallowUnknownFields()
			if err := dec.Decode(&input); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			// call decode again to confirm that the request contained only a single JSON object
			if err := dec.Decode(&struct{}{}); err != io.EOF {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			//log.Printf("server: %s %q: json: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		case "application/x-www-form-urlencoded":
			if err := r.ParseForm(); err != nil {
				log.Printf("server: %s %q: form: %+v\n", r.Method, r.URL.Path, err)
				http.SetCookie(w, &http.Cookie{Name: cookieName, Path: "/", MaxAge: -1, HttpOnly: true})
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			//log.Printf("server: %s %q: form: %v\n", r.Method, r.URL.Path, r.PostForm)
			for k, v := range r.Form {
				switch k {
				case "username":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Username = v[0]
				case "password":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Password = v[0]
				}
			}
			//log.Printf("server: %s %q: form: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		case "text/html":
			if err := r.ParseForm(); err != nil {
				log.Printf("server: %s %q: html: %+v\n", r.Method, r.URL.Path, err)
				http.SetCookie(w, &http.Cookie{Name: cookieName, Path: "/", MaxAge: -1, HttpOnly: true})
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			//log.Printf("server: %s %q: html: %v\n", r.Method, r.URL.Path, r.PostForm)
			for k, v := range r.Form {
				switch k {
				case "username":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Username = v[0]
				case "password":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Password = v[0]
				}
			}
			//log.Printf("server: %s %q: html: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		default:
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		}

		if input.Username == "" || input.Password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		//log.Printf("server: %s %q: %v\n", r.Method, r.URL.Path, input)

		u, err := s.store.FetchUserByCredentials(input.Username, input.Password)
		if err != nil {
			log.Printf("server: %s: %s: fetchUsersByCredentials: %v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		log.Printf("server: %s: %s: fetchUsersByCredentials: %q\n", r.Method, r.URL.Path, u.Handle)

		claims := map[string]interface{}{"user_id": strings.ToLower(u.Handle)}

		jwtauth.SetExpiryIn(claims, time.Second*7*24*60*60)
		_, tokenString, _ := jwtauth.New("HS256", s.key, nil).Encode(claims)

		switch contentType {
		case "application/json":
			//log.Printf("server: %s %q: json: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
			var response struct {
				Links struct {
					Self string `json:"self"`
				} `json:"links"`
				Data struct {
					Token string `json:"token"`
				} `json:"data,omitempty"`
			}
			response.Links.Self = r.URL.Path
			response.Data.Token = tokenString
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
			return
		case "application/x-www-form-urlencoded":
			//log.Printf("server: %s %q: form: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Path:     "/",
				Value:    tokenString,
				MaxAge:   14 * 24 * 60 * 60,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/ui", http.StatusSeeOther)
			return
		case "text/html":
			//log.Printf("server: %s %q: html: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Path:     "/",
				Value:    tokenString,
				MaxAge:   14 * 24 * 60 * 60,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/ui", http.StatusSeeOther)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *server) logoutHandler(cookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			MaxAge:   -1, // delete cookie
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) ordersGetHandler(templates string) http.HandlerFunc {
	t := osk.New(templates, "order_entry.html")

	type orderEntry struct {
		Game       string
		Year       string
		Quarter    string
		NationNo   int
		Rows, Cols int
		Orders     string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s: ordersPath %q\n", r.Method, r.URL.Path, s.store.OrdersPath())

		var claim *models.Claim
		_, claims, _ := jwtauth.FromContext(r.Context())
		if userId, ok := claims["user_id"].(string); !ok {
			log.Printf("%s: %s: claims[%q]: not a string\n", r.Method, r.URL.Path, "user_id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if claim, ok = s.claims[strings.ToLower(userId)]; !ok {
			log.Printf("%s: %s: claims[%q]: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		pGameName := chi.URLParam(r, "game")
		game, err := s.store.LookupGameByName(pGameName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		oe := orderEntry{
			Game:     game.ShortName,
			Year:     chi.URLParam(r, "year"),
			Quarter:  chi.URLParam(r, "quarter"),
			NationNo: claim.NationNo,
			Rows:     5,
			Cols:     80,
		}

		if oe.Year == "" || oe.Quarter == "" {
			oe.Year = "0000" //fmt.Sprintf("%04d", game.CurrentTurn.Year)
			oe.Quarter = "0" //fmt.Sprintf("%d", game.CurrentTurn.Quarter)
			http.Redirect(w, r, "/ui/games/PT-1/orders/0000/0", http.StatusTemporaryRedirect)
		} else if oe.Quarter == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ordersFile := filepath.Join(s.store.OrdersPath(), fmt.Sprintf("%s.%s.%s.%d.txt", game.ShortName, chi.URLParam(r, "year"), chi.URLParam(r, "quarter"), claim.NationNo))
		if b, err := os.ReadFile(ordersFile); err == nil {
			oe.Orders = string(b)
		}

		t.Handle(w, r, oe)
	}
}

func (s *server) ordersGetRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		userId, ok := claims["user_id"].(string)
		if !ok {
			log.Printf("%s: %s: claims[%q]: not a string\n", r.Method, r.URL.Path, "user_id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if _, ok = s.claims[strings.ToLower(userId)]; !ok {
			log.Printf("%s: %s: claims[%q]: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		pGameName := chi.URLParam(r, "game")
		t, err := s.store.FetchCurrentTurn(userId, pGameName)
		if err != nil {
			log.Printf("%s: %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/ui/games/%s/orders/%04d/%d", pGameName, t.Year, t.Quarter), http.StatusTemporaryRedirect)
	}
}

func (s *server) ordersPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s: entered\n", r.Method, r.URL.Path)

		_, claims, _ := jwtauth.FromContext(r.Context())
		userId, ok := claims["user_id"].(string)
		if !ok {
			log.Printf("%s: %s: claims[%q]: not a string\n", r.Method, r.URL.Path, "user_id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		claim, ok := s.claims[strings.ToLower(userId)]
		if !ok {
			log.Printf("%s: %s: claims[%q]: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		pGameName, pYear, pQuarter := chi.URLParam(r, "game"), chi.URLParam(r, "year"), chi.URLParam(r, "quarter")
		t, err := s.store.FetchCurrentTurn(userId, pGameName)
		if err != nil {
			log.Printf("%s: %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		currentTurn := pYear + "/" + pQuarter
		if currentTurn != t.String() {
			log.Printf("%s: %s: not current turn: %q %q\n", r.Method, r.URL.Path, currentTurn, t.String())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("%s: %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		//log.Printf("server: %s %q: %v\n", r.Method, r.URL.Path, r.PostForm)
		var input struct {
			orders string
		}
		for k, v := range r.Form {
			switch k {
			case "orders":
				if len(v) != 1 {
					log.Printf("%s: %s: too many forms.orders: %d\n", r.Method, r.URL.Path, len(v))
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				} else if len(v[0]) < 1 || len(v[0]) > 64*1024 {
					log.Printf("%s: %s: invalid orders length: %d\n", r.Method, r.URL.Path, len(input.orders))
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				} else if !utf8.ValidString(v[0]) || len(v[0]) < 1 || len(v[0]) > 64*1024 {
					log.Printf("%s: %s: invalid utf-8 string\n", r.Method, r.URL.Path)
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.orders = v[0]
			}
		}

		if len(input.orders) == 0 {
			log.Printf("%s: %s: invalid orders length: %d\n", r.Method, r.URL.Path, len(input.orders))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// try to replace characters we know the parser doesn't like
		o := input.orders
		for _, pair := range [][]string{
			{"\r\n", "\n"},
			{"\t", " "},
			{"\u2013", `-`},
			{"\u2014", `-`},
			{"\u2015", `-`},
			{"\u2017", `_`},
			{"\u2018", `'`},
			{"\u2019", `'`},
			{"\u201a", `,`},
			{"\u201b", `'`},
			{"\u201c", `"`},
			{"\u201d", `"`},
			{"\u201e", `"`},
			{"\u201f", `"`},
			{"\u2026", `...`},
			{"\u2032", `'`},
			{"\u2033", `"`},
		} {
			o = strings.ReplaceAll(o, pair[0], pair[1])
		}

		//log.Printf("%s: %s: ordersPath %q\n", r.Method, r.URL.Path, s.store.OrdersPath())

		ordersFile := filepath.Join(s.store.OrdersPath(), fmt.Sprintf("%s.%s.%s.%d.txt", pGameName, chi.URLParam(r, "year"), chi.URLParam(r, "quarter"), claim.NationNo))
		date := time.Now().UTC().Format(time.RFC3339)
		o = fmt.Sprintf(";; %s %d %s %s\n\n", pGameName, claim.NationNo, currentTurn, date) + o + "\n"
		if err := os.WriteFile(ordersFile, []byte(o), 0644); err != nil {
			log.Printf("%s: %s: writeFile %q: %v\n", r.Method, r.URL.Path, ordersFile, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
}

func (s *server) unitsGetHandler(templates string) http.HandlerFunc {
	units := s.store.FetchUnits()
	t := osk.New(templates, "units.html")
	return func(w http.ResponseWriter, r *http.Request) {
		t.Handle(w, r, units)
	}
}

func (s *server) myVerifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return s.myVerify(ja, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)
}

func (s *server) myVerify(ja *jwtauth.JWTAuth, findTokenFns ...func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := jwtauth.VerifyRequest(ja, r, findTokenFns...)
			ctx = jwtauth.NewContext(ctx, token, err)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
