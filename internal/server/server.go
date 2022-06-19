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

package server

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/jsonwt"
	"github.com/mdhender/jsonwt/signers"
	"github.com/mdhender/wraith/internal/way"
	"net"
	"net/http"
	"os"
	"time"
)

type Server struct {
	http.Server

	router *way.Router
	www    struct {
		host string
		port string
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

func New(opts ...Option) (*Server, error) {
	// create a server with default values
	s := &Server{}
	s.authn.Users = make(map[string]*user)
	s.router = way.NewRouter()
	s.www.host, s.www.port = "", "8080"

	s.Addr = net.JoinHostPort(s.www.host, s.www.port)
	s.MaxHeaderBytes = 1 << 20 // 1mb?
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second

	// apply the list of options to the server
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

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

	//s.router.HandleFunc("GET", "/ui/add-user", s.handleGetAddUser)
	//s.router.HandleFunc("POST", "/ui/add-user", s.handlePostAddUser)
	s.router.HandleFunc("GET", "/ui", s.authenticatedOnly(s.handleGetIndex))
	s.router.HandleFunc("GET", "/ui/login", s.handleGetLogin)
	s.router.HandleFunc("POST", "/ui/login", s.handlePostLogin)
	s.router.HandleFunc("GET", "/ui/logout", s.handleLogout)
	s.router.HandleFunc("POST", "/ui/logout", s.handleLogout)

	return s, nil
}

func (s *Server) Mux() http.Handler {
	return s.router
}

func (s *Server) handleGetIndex(w http.ResponseWriter, r *http.Request) {
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
