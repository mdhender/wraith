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
	"github.com/mdhender/wraith/cmd"
	"github.com/mdhender/wraith/internal/seeder"
	"log"
	"math/rand"
	"time"
)

func main() {
	defer func(started time.Time) {
		elapsed := time.Now().Sub(started)
		log.Printf("wraith: total time %v\n", elapsed)
	}(time.Now())

	//// default log format to UTC
	//log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	// seed the default PRNG source.
	if seed, err := seeder.Seed(); err != nil {
		log.Fatalln(err)
	} else {
		rand.Seed(seed)
	}

	// run the command as given
	cmd.Execute()
}
