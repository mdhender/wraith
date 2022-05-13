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

package seeder

import (
	crand "crypto/rand"
	"encoding/binary"
)

// Seed returns the best seed we can get from the operating system.
// It is suitable to be used directly as a seed or to build a new source.
// For example:
//  seed, err := seeder.Seed()
// Then
//  rand.New(rand.NewSource(seed))
// Or
//  rand.Seed(seed)
// Please don't ignore errors from this function.
func Seed() (int64, error) {
	var seed [8]byte
	if _, err := crand.Read(seed[:]); err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(seed[:])), nil
}
