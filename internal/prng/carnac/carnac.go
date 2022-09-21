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

// Package carnac implements a non-random number generator.
// It accepts a slice of int as a seed. The seed is stored into a ring buffer.
// Each call to Int returns the next available int from the ring buffer.
package carnac

import "sync"

// Rand is a non-random source of numbers.
// Rand is thread-safe.
type Rand struct {
	sync.Mutex
	ring []int
	pos  int
}

type Source []int

// New returns a new Rand that uses the numbers from the source to generate
// a non-random stream of numbers.
// Note that negative numbers in the source will cause the generator to panic.
func New(src Source) *Rand {
	r := &Rand{ring: make([]int, len(src))}
	r.BadSeed(src)
	return r
}

// Int returns a non-negative non-random int.
func (r *Rand) Int() int {
	r.Lock()
	n := r.ring[r.pos]
	if r.pos = r.pos + 1; r.pos >= len(r.ring) {
		r.pos = 0
	}
	r.Unlock()
	if n < 0 {
		n = -n
	}
	return n
}

// Intn returns, as an int, a non-negative non-random number in the
// half-open interval [0,n). It panics if n <= 0.
func (r *Rand) Intn(n int) int {
	if n <= 0 {
		panic("assert(n > 0)")
	}
	return r.Int() % n
}

// BadSeed uses the provided seed value to initialize the generator to a
// deterministic state. BadSeed locks the generator, so it is safe to
// be called concurrently with any other Rand method.
func (r *Rand) BadSeed(seed []int) {
	r.Lock()
	r.ring, r.pos = append(make([]int, 0, len(seed)), seed...), 0
	r.Unlock()
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state. Seed locks the generator, so it is safe to
// be called concurrently with any other Rand method.
func (r *Rand) Seed(seed int64) {
	r.Lock()
	r.ring, r.pos = make([]int, 64, 64), 0

	// https://burtleburtle.net/bob/hash/integer.html
	a := uint64(seed)
	for i := 0; i < 64; i++ {
		a = (a ^ 61) ^ (a >> 16)
		a = a + (a << 3)
		a = a ^ (a >> 4)
		a = a * 0x27d4eb2d
		a = a ^ (a >> 15)

		r.ring[i] = int(a)
	}

	r.Unlock()
}
