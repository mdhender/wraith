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

package main

import (
	"github.com/mdhender/wraith/internal/otohttp"
	"github.com/mdhender/wraith/internal/services/greeter"
	"github.com/mdhender/wraith/internal/services/identity"
	"log"
	"net/http"
)

func main() {
	otoServer, _ := otohttp.NewServer()

	identity.RegisterIdentityService(otoServer, identity.Service{})
	greeter.RegisterGreeterService(otoServer, greeter.Service{})

	http.Handle("/oto/", otoServer)

	log.Println("server listening on :8080")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
