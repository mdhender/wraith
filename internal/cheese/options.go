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
	"fmt"
	"github.com/mdhender/wraith/models"
	"log"
	"net"
	"os"
	"path/filepath"
)

type Option func(*Server) error

func WithHost(host string) func(*Server) error {
	return func(s *Server) error {
		s.host = host
		s.addr = net.JoinHostPort(s.host, s.port)
		return nil
	}
}

func WithKey(key []byte) func(*Server) error {
	return func(s *Server) error {
		if len(key) == 0 {
			return fmt.Errorf("missing key")
		}
		s.key = append(s.key, key...)
		return nil
	}
}

func WithPort(port string) func(*Server) error {
	return func(s *Server) error {
		s.port = port
		s.addr = net.JoinHostPort(s.host, s.port)
		return nil
	}
}

func WithStore(store *models.Store) func(*Server) error {
	return func(s *Server) error {
		// fetch user claims
		log.Printf("cheese.Serve: todo: needs game and date logic\n")
		claims, err := store.FetchClaims("0000/0")
		if err != nil {
			return err
		}
		s.store, s.claims = store, claims
		return nil
	}
}

func WithTemplates(path string) func(*Server) error {
	return func(s *Server) error {
		path = filepath.Clean(path)
		if fi, err := os.Stat(path); err != nil {
			return err
		} else if !fi.IsDir() {
			return fmt.Errorf("%s: not a directory", path)
		}
		s.templates = path
		return nil
	}
}
