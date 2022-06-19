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
	"github.com/mdhender/jsonwt"
	"log"
	"net/http"
)

// currentUser extracts data for the user making the request.
// It always returns a user struct, even if the request does not have a valid token.
func (s *Server) currentUser(r *http.Request) (user struct {
	IsAdmin         bool
	IsAuthenticated bool
	User            *user
}) {
	log.Printf("server: currentUser: entered\n")
	c, err := r.Cookie("jsonwt")
	if err != nil {
		log.Printf("server: currentUser: cookie: %+v\n", err)
		return user
	}
	token, err := jsonwt.Decode(c.Value)
	if err != nil {
		log.Printf("server: currentUser: token: %+v\n", err)
		return user
	} else if err = s.jwt.factory.Validate(token); err != nil {
		log.Printf("server: currentUser: validateToken %+v\n", err)
	} else if token.IsValid() {
		var claims []string
		if err := token.Claim(&claims); err != nil {
			log.Printf("server: currentUser: claims %+v\n", err)
		} else {
			for _, role := range claims {
				switch role {
				case "authenticated":
					user.IsAuthenticated = true
				default:
					// todo: get other roles from the claim
				}
			}
		}
	}
	return user
}
