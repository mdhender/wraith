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
	"net"
	"path/filepath"
	"time"
)

type Option func(*Server) error

// Options turns a list of Option instances into an Option.
func Options(opts ...Option) Option {
	return func(s *Server) error {
		for _, opt := range opts {
			if err := opt(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithAuthenticationData(root string) Option {
	return func(s *Server) (err error) {
		s.authn.source = filepath.Clean(root)
		return nil
	}
}

func WithHost(host string) Option {
	return func(s *Server) (err error) {
		s.www.host = host
		s.Addr = net.JoinHostPort(s.www.host, s.www.port)
		return nil
	}
}

func WithJwtData(root string) Option {
	return func(s *Server) (err error) {
		s.jwt.source = filepath.Clean(root)
		return nil
	}
}

func WithMaxBodyLength(l int) Option {
	return func(s *Server) (err error) {
		s.MaxHeaderBytes = l
		return nil
	}
}

func WithPort(port string) Option {
	return func(s *Server) (err error) {
		s.www.port = port
		s.Addr = net.JoinHostPort(s.www.host, s.www.port)
		return nil
	}
}

func WithReadTimeout(d time.Duration) Option {
	return func(s *Server) (err error) {
		s.ReadTimeout = d
		return nil
	}
}

func WithWriteTimeout(d time.Duration) Option {
	return func(s *Server) (err error) {
		s.WriteTimeout = d
		return nil
	}
}
