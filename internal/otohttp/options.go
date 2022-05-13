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

package otohttp

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"
)

// Option allows us to pass in options when creating a new server.
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

// WithAddr sets the http Addr field to the host and port.
func WithAddr(host, port string) Option {
	return func(s *Server) (err error) {
		s.Addr = net.JoinHostPort(host, port)
		return nil
	}
}

// WithBasepath changes the default from `/oto/`.
func WithBasepath(path string) Option {
	return func(s *Server) (err error) {
		s.Basepath = path
		return nil
	}
}

// WithContext adds a context to the server so that we can shut it down gracefully.
func WithContext(ctx context.Context) Option {
	return func(s *Server) (err error) {
		s.ctx = ctx
		return nil
	}
}

// WithDebugCors changes the default CORS debug flag from false to the given value.
func WithDebugCors(debug bool) Option {
	return func(s *Server) (err error) {
		s.debug.cors = debug
		return nil
	}
}

// WithMaxHeaderBytes updates the limit on header size.
func WithMaxHeaderBytes(n int) Option {
	return func(s *Server) (err error) {
		s.MaxHeaderBytes = n
		return nil
	}
}

// WithNotFound changes the default not found handler.
func WithNotFound(h http.Handler) Option {
	return func(s *Server) (err error) {
		s.NotFound = h
		return nil
	}
}

// WithOnErr changes the default error handler.
func WithOnErr(fn func(http.ResponseWriter, *http.Request, error)) Option {
	return func(s *Server) (err error) {
		s.OnErr = fn
		return nil
	}
}

// WithReadTimeout updates the timeout.
// Example: 5 * time.Second
func WithReadTimeout(t time.Duration) Option {
	return func(s *Server) (err error) {
		s.ReadTimeout = t
		return nil
	}
}

// WithWorkingDir changes the working directory for the server.
func WithWorkingDir(path string) Option {
	return func(s *Server) (err error) {
		return os.Chdir(path)
	}
}

// WithWriteTimeout updates the timeout.
// Example: 10 * time.Second
func WithWriteTimeout(t time.Duration) Option {
	return func(s *Server) (err error) {
		s.WriteTimeout = t
		return nil
	}
}
