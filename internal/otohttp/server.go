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

// Package otohttp implements an Oto HTTP Server.
package otohttp

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// Server handles oto requests.
type Server struct {
	http.Server
	ctx   context.Context
	debug struct {
		cors bool
	}
	routes map[string]http.Handler

	// Basepath is the path prefix to match.
	// Default: /oto/
	Basepath string

	// NotFound is the http.Handler to use when a resource is not found.
	NotFound http.Handler

	// OnErr is called when there is an error.
	OnErr func(w http.ResponseWriter, r *http.Request, err error)
}

// NewServer makes a new Server.
func NewServer(opts ...func(*Server) error) (*Server, error) {
	s := &Server{
		Basepath: "/oto/",
		OnErr: func(w http.ResponseWriter, r *http.Request, err error) {
			errObj := struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			}
			if err := Encode(w, r, http.StatusInternalServerError, errObj); err != nil {
				log.Printf("failed to encode error: %s\n", err)
			}
		},
		NotFound: http.NotFoundHandler(),
		routes:   make(map[string]http.Handler),
	}

	// update defaults for port, timeouts, input limits, context.
	s.Addr = net.JoinHostPort("", "8080")
	s.BaseContext = func(_ net.Listener) context.Context { return s.ctx }
	s.MaxHeaderBytes = 16 * 1024
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second

	// apply the list of options to update the server configuration
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Register adds a handler for the specified service method.
func (s *Server) Register(service, method string, h http.HandlerFunc) {
	log.Printf("server: registering %s%s.%s\n", s.Basepath, service, method)
	s.routes[fmt.Sprintf("%s%s.%s", s.Basepath, service, method)] = h
}

// ServeHTTP serves the request.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, HEAD, OPTIONS, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == http.MethodOptions {
		if s.debug.cors {
			log.Printf("[cors] %s %q\n", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		s.NotFound.ServeHTTP(w, r)
		return
	}
	h, ok := s.routes[r.URL.Path]
	if !ok {
		s.NotFound.ServeHTTP(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

// Encode writes the response.
func Encode(w http.ResponseWriter, r *http.Request, status int, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "encode json")
	}
	var out io.Writer = w
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		out = gzw
		defer gzw.Close()
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := out.Write(b); err != nil {
		return err
	}
	return nil
}

// Decode unmarshals the object in the request into v.
func Decode(r *http.Request, v interface{}) error {
	bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		return fmt.Errorf("decode: read body: %w", err)
	}
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return fmt.Errorf("decode: json.Unmarshal: %w", err)
	}
	return nil
}
