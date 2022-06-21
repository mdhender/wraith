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

// Package main generates playtest game data.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/wraith/internal/seeder"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	started := time.Now()

	// default log format to UTC
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	// use the crand package to seed the default PRNG source.
	if seed, err := seeder.Seed(); err != nil {
		log.Fatalln(err)
	} else {
		rand.Seed(seed)
	}

	if err := run(); err != nil {
		log.Println(err)
	}

	elapsed := time.Now().Sub(started)
	fmt.Printf("total time: %+v\n", elapsed)
}

func run() error {
	numberOfNations, radius := 14, 8
	game := GenGame(numberOfNations, radius)

	if buf, err := json.MarshalIndent(game, "", "  "); err != nil {
		return err
	} else if err := os.WriteFile("ptest.json", buf, 0600); err != nil {
		return err
	}

	fmt.Printf("asteroid belts   %8d\n", numAsteroidBelts)
	fmt.Printf("gas giants       %8d\n", numGasGiants)
	fmt.Printf("terrestrials     %8d\n", numTerrestrials)
	fmt.Printf("                 %8d\n", numPlanets)

	//var systems []*System
	//if data, err := os.ReadFile("systems.json"); err != nil {
	//	return err
	//} else if err = json.Unmarshal(data, &systems); err != nil {
	//	return err
	//}
	//for _, s := range systems {
	//	if s.NumberOfStars > 0 {
	//		s.Stars = make([]*Star, s.NumberOfStars, s.NumberOfStars)
	//	}
	//}
	//
	//var planets []*Planet
	//if data, err := os.ReadFile("planets.json"); err != nil {
	//	return err
	//} else if err = json.Unmarshal(data, &planets); err != nil {
	//	return err
	//}
	//
	//for _, p := range planets {
	//	// assign habitability number to terrestrials
	//	if p.Kind == "terrestrial" && (1 <= p.Orbit && p.Orbit <= 5) {
	//		if p.HabitabilityNumber == 0 {
	//			if p.Special {
	//				p.HabitabilityNumber = 25
	//			} else {
	//				p.HabitabilityNumber = roll(2, 10) + roll(1, 5)
	//			}
	//		}
	//	}
	//	// create deposits
	//	if len(p.Deposits) == 0 {
	//		n := roll(4, 10)
	//		if p.Special {
	//			n = 35
	//		}
	//		if n > 0 {
	//			p.Deposits = make([]*NaturalResource, n, n)
	//		}
	//	}
	//	for i := range p.Deposits {
	//		if p.Deposits[i] != nil {
	//			continue
	//		}
	//		d := &NaturalResource{Id: i + 1}
	//		if p.Special && i == 0 {
	//			d.Kind = "gold"
	//			d.Yield = 0.07
	//			d.Quantity = 234_567
	//		} else if p.Special && i == 1 {
	//			d.Kind = "metallic"
	//			d.Yield = 0.99
	//			d.Quantity = 99_999_999
	//		} else if p.Special && i == 2 {
	//			d.Kind = "non-metallic"
	//			d.Yield = 0.89
	//			d.Quantity = 89_999_999
	//		} else if p.Special && i == 3 {
	//			d.Kind = "fuel"
	//			d.Yield = 0.79
	//			d.Quantity = 79_999_999
	//		} else {
	//			switch roll(1, 6) {
	//			case 0, 1, 2:
	//				d.Kind = "metallic"
	//				d.Yield = float64(roll(10, 10)+1) / 100
	//			case 3, 4:
	//				d.Kind = "non-metallic"
	//				d.Yield = float64(roll(8, 10)+1) / 100
	//			case 5:
	//				d.Kind = "fuel"
	//				d.Yield = float64(roll(5, 10)+1) / 100
	//			}
	//			d.Quantity = roll(99, 1_000_000)
	//		}
	//		p.Deposits[i] = d
	//	}
	//}
	//
	//for _, s := range systems {
	//	for i := range s.Stars {
	//		if i == 0 && s.NumberOfStars == 1 {
	//			s.Stars[i] = &Star{
	//				Kind:   "A",
	//				Orbits: planets,
	//			}
	//		} else {
	//			s.Stars[i] = &Star{
	//				Kind: "B",
	//			}
	//		}
	//	}
	//}

	return nil
}

func roll(n, d int) int {
	total := 0
	for i := 0; i < n; i++ {
		total += rand.Intn(d)
	}
	return total
}
