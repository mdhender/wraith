/*
 * Wraith Game Engine
 * Copyright (c) 2022 Michael D. Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package server

import (
	"context"
	"github.com/mdhender/wraith/internal/config"
	"github.com/mdhender/wraith/internal/way"
	"net"
	"net/http"
	"time"
)

func New(cfg *config.Config, opts ...func(*Server) error) (*Server, error) {
	s := &Server{
		router: way.NewRouter(),
	}
	s.Addr = net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	s.BaseContext = func(_ net.Listener) context.Context { return s.ctx }
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20 // 1mb?

	// apply the list of options to Store
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	s.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	return s, nil
}

type Server struct {
	http.Server
	ctx    context.Context
	router *way.Router
}

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

func WithContext(ctx context.Context) Option {
	return func(s *Server) (err error) {
		s.ctx = ctx
		return nil
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
