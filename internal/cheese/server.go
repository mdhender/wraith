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
	"github.com/go-chi/jwtauth/v5"
	"github.com/mdhender/wraith/models"
	"net/http"
)

type Server struct {
	debug            bool
	addr, host, port string
	gamesPath        string
	key              []byte
	store            *models.Store
	claims           map[string]*models.Claim
	templates        string
}

func Serve(options ...Option) error {
	s := &Server{port: "3000"}
	for _, opt := range options {
		if err := opt(s); err != nil {
			return err
		}
	}
	if len(s.gamesPath) == 0 {
		return fmt.Errorf("missing games path")
	} else if len(s.key) == 0 {
		return fmt.Errorf("missing key")
	} else if s.templates == "" {
		return fmt.Errorf("missing templates path")
	}

	return s.serve()
}

func (s *Server) myVerifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return s.myVerify(ja, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)
}

func (s *Server) myVerify(ja *jwtauth.JWTAuth, findTokenFns ...func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token, err := jwtauth.VerifyRequest(ja, r, findTokenFns...)
			ctx = jwtauth.NewContext(ctx, token, err)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
