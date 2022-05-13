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

// Package words implements a passphrase generator inspired by https://xkcd.com/936/.
package words

import (
	"encoding/json"
	"github.com/mdhender/wraith/internal/seeder"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"sync"
)

// New returns an initialized passphrase generator using the default word list.
func New(separators string) (*Generator, error) {
	return NewFromList(defaultWordList(), separators)
}

// NewFromFile returns an initialized passphrase generator after loading
// the list of words from a JSON array in the given file.
func NewFromFile(file, separators string) (*Generator, error) {
	var list []string
	if b, err := os.ReadFile(file); err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	return NewFromList(list, separators)
}

// NewFromList returns an initialized passphrase generator using the supplied
// list of words.
func NewFromList(list []string, separators string) (*Generator, error) {
	g := &Generator{words: list}
	if len(g.words) < 2048 {
		return nil, errors.New("too few words")
	}
	g.setSeparators(separators)
	if err := g.seed(); err != nil {
		return nil, err
	}
	return g, nil
}

// Generator implements a passphrase generator inspired by https://xkcd.com/936/
type Generator struct {
	sync.Mutex
	rnd        *rand.Rand
	words      []string
	separators []string
}

// Generate returns a passphrase with the given number of words.
func (g *Generator) Generate(n int) (passphrase string) {
	g.Lock()
	defer g.Unlock()
	for ; n > 0; n-- {
		if len(passphrase) != 0 {
			passphrase += g.separators[g.rnd.Int()%len(g.separators)]
		}
		passphrase += g.words[g.rnd.Int()%len(g.words)]
	}
	return passphrase
}

func (g *Generator) seed() error {
	g.Lock()
	defer g.Unlock()
	g.rnd = nil
	if seed, err := seeder.Seed(); err != nil {
		return err
	} else {
		g.rnd = rand.New(rand.NewSource(seed))
	}
	return nil
}

func (g *Generator) setSeparators(separators string) {
	for _, r := range separators {
		g.separators = append(g.separators, string(r))
	}
	if len(g.separators) == 0 {
		switch g.rnd.Int() % 5 {
		case 0:
			g.separators = append(g.separators, " ")
		case 1:
			g.separators = append(g.separators, ".")
		case 2:
			g.separators = append(g.separators, "+")
		case 3:
			g.separators = append(g.separators, "-")
		default:
			g.separators = append(g.separators, " ", ".", "+", "-")
		}
	}
}
