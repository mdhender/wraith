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
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"unicode/utf8"
)

func (s *Server) handleGetAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	jsonwt.DeleteCookie(w)
	page := fmt.Sprintf(`<body>
				<h1>Wraith Create Account</h1>
				<form action="/ui/add-user"" method="post">
					<table>
						<tr><td align="right">Handle&nbsp;</td><td><input type="text" name="handle"></td></tr>
						<tr><td align="right">Password&nbsp;</td><td><input type="password" name="password"></td></tr>
						<tr><td>&nbsp;</td><td align="right"><input type="submit" value="Create"></td></tr>
					</table>
				</form>
			</body>`)
	_, _ = w.Write([]byte(page))
}

func (s *Server) handlePostAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf("server: %s %q: entered\n", r.Method, r.URL.Path)
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
		log.Printf("server: read: %+v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if err = json.Unmarshal(b, &s.authn); err != nil {
		log.Printf("server: json: %+v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	acct, ok := s.authn.Users[input.handle]
	if !ok {
		log.Printf("server: %s %q: user %q not found\n", r.Method, r.URL.Path, input.handle)
		acct = &user{
			Id:           input.handle,
			Handle:       input.handle,
			Secret:       input.password,
			HashedSecret: "",
		}
		s.authn.Users[input.handle] = acct
		log.Printf("server: %s %q: user %q created %q\n", r.Method, r.URL.Path, acct.Handle, *acct)
	}
	if acct.HashedSecret == "" {
		hashedSecretBytes, err := bcrypt.GenerateFromPassword([]byte(input.password), s.authn.Bcrypt.MinCost)
		if err != nil {
			log.Printf("server: bcrypt: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		acct.HashedSecret = string(hashedSecretBytes)
	}
	log.Printf("server: %s %q: handle %q password %q hashed %q\n", r.Method, r.URL.Path, input.handle, input.password, acct.HashedSecret)

	if data, err := json.MarshalIndent(s.authn, "", "  "); err != nil {
		log.Printf("server: json: %+v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if err = os.WriteFile(s.authn.source, data, 0600); err != nil {
		log.Printf("server: write: %+v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ui/add-user", http.StatusOK)
	return
}
