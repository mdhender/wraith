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

package jdb

import (
	"encoding/json"
	"log"
	"os"
)

func Load(filename string) (*Game, error) {
	log.Printf("jdb: loading %s\n", filename)
	var g Game
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (g *Game) Write(filename string) error {
	log.Printf("jdb: saving %s\n", filename)
	b, err := json.MarshalIndent(g, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b, 0666)
	if err != nil {
		return err
	}
	return nil
}
