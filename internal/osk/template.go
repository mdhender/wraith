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

package osk

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Handler struct {
	root     string
	filename string
	t        *template.Template
}

// New returns a new handler
func New(root, filename string) *Handler {
	return &Handler{filename: filepath.Join(root, filename)}
}

func (t *Handler) Handle(w http.ResponseWriter, r *http.Request, data interface{}) {
	x := template.Must(template.ParseFiles(t.filename))
	err := x.Execute(w, data)
	if err != nil {
		log.Printf("osk: %q: %v\n", t.filename, err)
	}
}

// ServeHTTP implements the http.Handler interface
func (t *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.t = template.Must(template.ParseFiles(t.filename))
	w.Header().Set("Content-Type", "text/html")
	err := t.t.Execute(w, nil)
	if err != nil {
		log.Printf("osk: %q: %v\n", t.filename, err)
	}
}
