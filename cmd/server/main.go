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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/wraith/internal/way"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"unicode/utf8"

	//"github.com/mdhender/wraith/internal/otohttp"
	//"github.com/mdhender/wraith/internal/services/greeter"
	//"github.com/mdhender/wraith/internal/services/identity"
	"github.com/mdhender/jsonwt"
	"github.com/mdhender/jsonwt/signers"
	"log"
	"net/http"
)

func main() {
	//otoServer, _ := otohttp.NewServer()
	//
	//identity.RegisterIdentityService(otoServer, identity.Service{})
	//greeter.RegisterGreeterService(otoServer, greeter.Service{})
	//
	//http.Handle("/oto/", otoServer)
	//
	log.Fatalln(run())
}

func run() error {
	s, err := newServer("D:\\wraith\\testdata\\authn.json", "D:\\wraith\\testdata\\jsonwt.json")
	if err != nil {
		return err
	}
	log.Printf("running on :3030\n")
	return http.ListenAndServe(":3030", s.router)
}

type server struct {
	router *way.Router
	www    struct {
		root string
	}
	authn struct {
		source  string
		Version string `json:"version"`
		Bcrypt  struct {
			MinCost int `json:"min_cost"`
		} `json:"bcrypt"`
		Users map[string]*user `json:"users"` // key is user handle
	}
	jwt struct {
		source   string
		Version  string `json:"version"`
		TTLHours int    `json:"ttl-hours"`
		Key      struct {
			Name   string `json:"name"`
			Secret string `json:"secret"`
		} `json:"key"`
		factory *jsonwt.Factory
		ttl     time.Duration
	}
}

type user struct {
	Id           string `json:"id"`
	Handle       string `json:"handle"`
	Secret       string `json:"secret"`
	HashedSecret string `json:"hashed-secret"`
}

func newServer(authenticationFilename, jsonwtFilename string) (*server, error) {
	var s server
	s.authn.Users = make(map[string]*user)

	s.authn.source = authenticationFilename
	if b, err := os.ReadFile(s.authn.source); err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &s.authn); err != nil {
		return nil, err
	} else {
		for id, u := range s.authn.Users {
			if u.Id != id {
				u.Id = id
			}
			if u.Handle == "" {
				u.Handle = u.Id
			}
		}
	}

	s.jwt.source = jsonwtFilename
	if b, err := os.ReadFile(s.jwt.source); err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &s.jwt); err != nil {
		return nil, err
	} else if hsSigner, err := signers.NewHS256([]byte(s.jwt.Key.Secret)); err != nil {
		return nil, err
	} else {
		s.jwt.factory = jsonwt.NewFactory(s.jwt.Key.Name, hsSigner)
		s.jwt.ttl = time.Hour * time.Duration(s.jwt.TTLHours)
	}

	s.router = way.NewRouter()
	//s.router.HandleFunc("GET", "/ui/add-user", s.handleGetAddUser)
	//s.router.HandleFunc("POST", "/ui/add-user", s.handlePostAddUser)
	s.router.HandleFunc("GET", "/ui", s.authenticatedOnly(s.handleGetIndex))
	s.router.HandleFunc("GET", "/ui/login", s.handleGetLogin)
	s.router.HandleFunc("POST", "/ui/login", s.handlePostLogin)
	s.router.HandleFunc("GET", "/ui/logout", s.handleLogout)
	s.router.HandleFunc("POST", "/ui/logout", s.handleLogout)

	return &s, nil
}

func (s *server) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	page := fmt.Sprintf(`<body>
				<h1>Wraith UI</h1>
			</body>`)
	_, _ = w.Write([]byte(page))
}

func (s *server) handleGetLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	jsonwt.DeleteCookie(w)
	page := fmt.Sprintf(`<body>
				<h1>Wraith Login</h1>
				<form action="/ui/login"" method="post">
					<table>
						<tr><td align="right">Handle&nbsp;</td><td><input type="text" name="handle"></td></tr>
						<tr><td align="right">Password&nbsp;</td><td><input type="password" name="password"></td></tr>
						<tr><td>&nbsp;</td><td align="right"><input type="submit" value="Login"></td></tr>
					</table>
				</form>
			</body>`)
	_, _ = w.Write([]byte(page))
}

func (s *server) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf("server: %s %q: handlePostLogin\n", r.Method, r.URL.Path)
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("server: %s %q: form: %v\n", r.Method, r.URL.Path, r.PostForm)
	var input struct {
		handle   string
		password string
	}
	for k, v := range r.Form {
		switch k {
		case "handle":
			if len(v) != 1 || !utf8.ValidString(v[0]) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			input.handle = v[0]
		case "password":
			if len(v) != 1 || !utf8.ValidString(v[0]) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			input.password = v[0]
		}
	}

	if b, err := os.ReadFile(s.authn.source); err != nil {
		log.Printf("server: %s %q: read: %+v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if err = json.Unmarshal(b, &s.authn); err != nil {
		log.Printf("server: %s %q: json: %+v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	acct, ok := s.authn.Users[input.handle]
	if !ok {
		log.Printf("server: %s %q: handle %q: not found\n", r.Method, r.URL.Path, input.handle)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	log.Printf("server: %s %q: handle %q: password %q hashed %q\n", r.Method, r.URL.Path, input.handle, input.password, acct.HashedSecret)

	if err := bcrypt.CompareHashAndPassword([]byte(acct.HashedSecret), []byte(input.password)); err != nil {
		log.Printf("server: %s %q: bcrypt: %+v\n", r.Method, r.URL.Path, err)
		log.Printf("server: %s %q: handle %q: invalid password\n", r.Method, r.URL.Path, acct.Handle)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	j, err := s.jwt.factory.Token(s.jwt.ttl, []string{"authenticated"})
	if err != nil {
		log.Printf("server: %s %q: token: %+v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	j.SetCookie(w)
	http.Redirect(w, r, "/ui", http.StatusSeeOther)
	return
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
	jsonwt.DeleteCookie(w)
	http.Redirect(w, r, "/ui", http.StatusSeeOther)
}
