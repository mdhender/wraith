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
	"net/http"
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

		r.Get("/ui/logout", s.handleLogout(tokenCookie))
		r.Post("/ui/logout", s.handleLogout(tokenCookie))
		r.Put("/ui/logout", s.handleLogout(tokenCookie))
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
			r.Get("/games/{game}/cluster", func(w http.ResponseWriter, r *http.Request) {
				pGameName, x, y, z := chi.URLParam(r, "game"), 0, 0, 0
				game, err := s.store.LookupGameByName(pGameName)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				b, err := s.store.FetchClusterByGameOrigin(game.Id)
				if err != nil {
					log.Printf("%s: %s: fetchCluster: %d: %v\n", r.Method, r.URL.Path, game.Id, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				_, _ = w.Write([]byte(fmt.Sprintf("<body><h1>Cluster Map: Origin %d/%d/%d</h1>", x, y, z)))
				_, _ = w.Write(b)
				_, _ = w.Write([]byte("</body>"))
			})
			r.Get("/games/{game}/cluster/{x}/{y}/{z}", func(w http.ResponseWriter, r *http.Request) {
				game, err := s.store.LookupGameByName(chi.URLParam(r, "game"))
				if err != nil {
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				x, err := strconv.Atoi(chi.URLParam(r, "x"))
				if err != nil {
					log.Printf("%s: %s: x: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				y, err := strconv.Atoi(chi.URLParam(r, "y"))
				if err != nil {
					log.Printf("%s: %s: y: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				z, err := strconv.Atoi(chi.URLParam(r, "z"))
				if err != nil {
					log.Printf("%s: %s: z: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				b, err := s.store.FetchClusterByGame(game.Id, x, y, z)
				if err != nil {
					log.Printf("%s: %s: fetchCluster: %d: %v\n", r.Method, r.URL.Path, game.Id, err)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				_, _ = w.Write([]byte(fmt.Sprintf("<body><h1>Cluster Map: Origin %d/%d/%d</h1>", x, y, z)))
				_, _ = w.Write(b)
				_, _ = w.Write([]byte("</body>"))
			})
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
			r.Get("/units", s.unitsGetHandler(s.templates))
		})
	})

	log.Printf("server: listening on %s\n", s.addr)
	return http.ListenAndServe(s.addr, r)
}

func (s *server) handleLogout(cookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusNoContent)
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
