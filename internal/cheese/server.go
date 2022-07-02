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
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mdhender/wraith/models"
	"io"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

type server struct {
	addr   string
	debug  bool
	key    []byte
	store  *models.Store
	claims map[string]interface{}
}

func Serve(addr string, key []byte, store *models.Store) error {
	s := server{addr: addr, key: key, store: store}
	s.claims = make(map[string]interface{})
	s.claims["mdhender"] = struct{ UserId string }{UserId: "mdhender"}

	//// fetch users
	//store.VerifyUserByCredentials()

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

		r.Get("/ui/login", s.loginGetHandler(tokenCookie, tokenString))
		r.Post("/ui/login", s.loginPostHandler(tokenCookie, tokenString))
		//r.Get("/ui/login/{handle}/{secret}", s.loginGetHandleSecretHandler(tokenCookie, tokenString))

		r.Get("/ui/logout", s.handleLogout(tokenCookie))
		r.Post("/ui/logout", s.handleLogout(tokenCookie))
		r.Put("/ui/logout", s.handleLogout(tokenCookie))
	})

	// protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens.
		// In this example, we use the provided authenticator middleware, but you can write your own very easily.
		// Look at the Authenticator method in jwtauth.go and tweak it; it's not scary.
		r.Use(jwtauth.Authenticator)

		r.Route("/api", func(r chi.Router) {
			r.Get("/report/game/{game}/nation/{nation}/turn/{year}/{quarter}", func(w http.ResponseWriter, r *http.Request) {
				game, nation, year, quarter := chi.URLParam(r, "game"), chi.URLParam(r, "nation"), chi.URLParam(r, "year"), chi.URLParam(r, "quarter")
				_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: game %q nation %q year %q quarter %q", r.Method, r.URL.Path, game, nation, year, quarter)))
			})
		})

		r.Get("/ui", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("yay!"))
		})
		r.Get("/ui/games/{game}/nations/{nation}/turn/{year}/{quarter}/report", func(w http.ResponseWriter, r *http.Request) {
			gameParam := chi.URLParam(r, "game")
			nationParam := chi.URLParam(r, "nation")
			yearParam := chi.URLParam(r, "year")
			quarterParam := chi.URLParam(r, "quarter")
			_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: game %q nation %q year %q quarter %q", r.Method, r.URL.Path, gameParam, nationParam, yearParam, quarterParam)))
		})
		r.Get("/ui/security", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: claims[user_id] %q", r.Method, r.URL.Path, claims["user_id"])))
		})
		r.Get("/ui/panic", func(http.ResponseWriter, *http.Request) {
			panic("foo")
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

func (s *server) loginGetHandler(cookieName string, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("server: %s: %s\n", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "text/html")
		// delete any existing cookie
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			Value:    token,
			MaxAge:   -1,
			HttpOnly: true,
		})
		page := `<body>
				<h1>Wraith Reactor</h1>
				<form action="/ui/login"" method="post">
					<table>
						<tr><td align="right">Username&nbsp;</td><td><input type="text" name="username"></td></tr>
						<tr><td align="right">Password&nbsp;</td><td><input type="password" name="password"></td></tr>
						<tr><td>&nbsp;</td><td align="right"><input type="submit" value="Login"></td></tr>
					</table>
				</form>
			</body>`
		_, _ = w.Write([]byte(page))
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
			response.Data.Token = "value"
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
			return
		case "application/x-www-form-urlencoded":
			//log.Printf("server: %s %q: form: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Path:     "/",
				Value:    token,
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
				Value:    token,
				MaxAge:   14 * 24 * 60 * 60,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/ui", http.StatusSeeOther)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
