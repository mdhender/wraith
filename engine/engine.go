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

package engine

import (
	"errors"
	"github.com/mdhender/wraith/storage/config"
	"log"
	"path/filepath"
)

/*
type System struct {
	Game        *Game
	Coordinates Coordinates
	Stars       []*Star
}

type Coordinates struct {
	X int
	Y int
	Z int
}

type Star struct {
	System   *System
	Sequence string // A, B, etc
	Kind     string
	Orbits   [11]*Planet // each orbit may or may not contain a planet
}

type Planet struct {
	Star           *Star
	OrbitNo        int    // 1..10
	Kind           string // asteroid belt, gas giant, terrestrial
	HabitabilityNo int
	ControlledBy   *Nation
	Deposits       []*Deposit
	Colonies       []*Colony
	Ships          []*Ship
}

type Deposit struct {
	Planet           *Planet
	ControlledBy     *Nation
	Unit             string  // fuel, gold, metallics, non-metallics
	QtyInitial       int     // in mass units
	QtyRemaining     int     // in mass units
	MiningDifficulty float64 // how hard it is to extract each mass unit
	YieldPct         float64 // percentage of each mass unit that yields units
}

type Colony struct {
	Id            int
	Location      *Planet
	Kind          string // surface colony, enclosed colony, orbital colony
	TechLevel     int
	BuiltBy       *Nation
	ControlledBy  *Nation
	Inventory     []*Inventory
	MiningGroups  []*MiningGroup
	FactoryGroups []*FactoryGroup
}

type Ship struct {
	Id            int
	Location      *Planet
	TechLevel     int
	BuiltBy       *Nation
	ControlledBy  *Nation
	Inventory     []*Inventory
	FactoryGroups []*FactoryGroup
}

type FactoryGroup struct {
	Colony    *Colony
	Ship      *Ship
	GroupNo   int
	Inventory []*Inventory
	Unit      string
	TechLevel int
}

type MiningGroup struct {
	Colony    *Colony
	GroupNo   int
	Deposit   *Deposit
	Inventory []*Inventory
}

type Inventory struct {
	Unit           string
	TechLevel      int
	QtyOperational int
	QtyStowed      int
	TotalMass      int
	EnclosedMass   int
}

type Nation struct {
	Player   *Player
	Name     string
	Colonies []*Colony
	Ships    []*Ship
}

type Player struct {
	Game   *Game
	User   *User
	Nation *Nation
}
*/

type Engine struct {
	config struct {
		base *config.Global
	}
	stores struct {
		games   *Games
		game    *Game
		nations *Nations
		turns   *Turns
		users   *Users
	}
	s *Store
}

// New returns an initialized engine with the base configuration
// and the games store loaded.
func New(baseConfigFile string) (e *Engine, err error) {
	if baseConfigFile == "" {
		return nil, errors.New("missing base config")
	}

	e = &Engine{}

	// load the base configuration
	e.config.base, err = config.LoadGlobal(baseConfigFile)
	if err != nil {
		return nil, err
	}
	log.Printf("loaded base config %q\n", e.config.base.Self)

	// load the users store
	e.stores.users = &Users{
		Store: e.RootDir("users"),
		Index: []UsersIndex{},
	}
	if err = e.ReadUsers(); err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded users store %q\n", e.stores.users.Store)

	// load the games store
	e.stores.games = &Games{
		Store: e.RootDir("games"),
		Index: []GamesIndex{},
	}
	if err = e.ReadGames(); err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded games store %q\n", e.stores.games.Store)

	return e, nil
}

func (e *Engine) RootDir(stores ...string) string {
	if len(stores) == 0 {
		return "D:\\wraith\\testdata"
	}
	return filepath.Join(append([]string{"D:\\wraith\\testdata"}, stores...)...)
}

func CreateGame(string, string, string, bool) (*Game, error) {
	panic("!")
}
func CreateNations(string, bool) (*Nation, error) {
	panic("!")
}
func LoadGames(string) (*Games, error) {
	panic("!")
}
func LoadNations(string) (*Nations, error) {
	panic("!")
}

func (e *Engine) LoadGame(game string) (err error) {
	// free up any game already in memory
	e.stores.game, e.stores.nations, e.stores.turns = nil, nil, nil

	if game == "" {
		return errors.New("missing game name")
	}

	//// find the game in the store
	//for _, g := range e.stores.games.Index {
	//	if strings.ToLower(g.Name) == strings.ToLower(game) {
	//		e.stores.game, err = LoadGame(g.Store)
	//		if err != nil {
	//			return err
	//		}
	//		break
	//	}
	//}
	if e.stores.game == nil {
		log.Fatalf("unable to find game %q\n", game)
	}
	log.Printf("loaded game store %q\n", e.stores.game.Store)

	//// use the game store to load the nations store
	//e.stores.nations, err = LoadNations(e.stores.game.Store)
	//if err != nil {
	//	return err
	//}
	log.Printf("loaded nations store %q\n", e.stores.nations.Store)

	return nil
}

func (e *Engine) Version() string {
	return "0.1.0"
}
