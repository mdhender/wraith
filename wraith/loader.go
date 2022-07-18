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

package wraith

import "fmt"

type Engine struct {
	Version string
	Game    struct {
		Code string
		Turn struct {
			Year    int
			Quarter int
		}
	}
	Colonies      map[string]*CorS
	CorSById      map[int]*CorS
	Deposits      map[int]*Deposit
	FactoryGroups map[int]*FactoryGroup
	FarmGroups    map[int]*FarmGroup
	MineGroups    map[int]*MineGroup
	Nations       map[int]*Nation
	Planets       map[int]*Planet
	Players       map[int]*Player
	Ships         map[string]*CorS
	Stars         map[int]*Star
	Systems       map[int]*System
	Units         map[int]*Unit
}

// CorS is a colony or ship
type CorS struct {
	Id            int     // unique identifier
	Kind          string  // orbital, ship, or surface
	HullId        string  // S or C + MSN
	MSN           int     // manufacturer serial number; in game id for the colony or ship
	BuiltBy       *Nation // nation that originally built the colony or ship
	Name          string  // name of this colony or ship
	TechLevel     int     // tech level of this colony or ship
	ControlledBy  *Player // player that controls this colony or ship
	Planet        *Planet // planet the colony or ship is located at
	Hull          []*HullUnit
	Inventory     []*InventoryUnit
	Population    Population
	Pay           Pay
	Rations       Rations
	FactoryGroups FactoryGroups // list of the factory groups
	FarmGroups    FarmGroups    // list of the farm groups
	MineGroups    MineGroups    // list of the mine groups
}

type CorSs []*CorS

func (c CorSs) Len() int {
	return len(c)
}

func (c CorSs) Less(i, j int) bool {
	return c[i].MSN < c[j].MSN
}

func (c CorSs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Coordinates are the x,y,z coordinates of a system
type Coordinates struct {
	X int
	Y int
	Z int
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%d/%d/%d", c.X, c.Y, c.Z)
}

// Deposit of fuel, gold, metal, or non-metals on the surface of a planet
type Deposit struct {
	Id           int     // unique identifier
	No           int     // number of deposit on planet
	Product      *Unit   // fuel, gold, metallic, non-metallic
	InitialQty   int     // in metric tonnes
	RemainingQty int     // in metric tonnes
	YieldPct     float64 // percentage of each mass unit that yields units
	Planet       *Planet // planet deposit is on
	ControlledBy *CorS   // colony controlling this deposit
}

type Deposits []*Deposit

func (d Deposits) Len() int {
	return len(d)
}

func (d Deposits) Less(i, j int) bool {
	return d[i].Id < d[j].Id
}

func (d Deposits) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// FactoryGroup is a group of factories on a ship or colony.
// Each group is dedicated to manufacturing one type of unit.
type FactoryGroup struct {
	CorS     *CorS                // ship or colony that controls the group
	Id       int                  // unique identifier
	No       int                  // group number, range 1...255
	Product  *Unit                // unit being manufactured
	Units    []*FactoryGroupUnits // factory units in the group
	StageQty [4]int               // assumes four turns to produce a single unit
}

type FactoryGroups []*FactoryGroup

// FactoryGroupUnits is the number of factories working together in the group
type FactoryGroupUnits struct {
	Unit     *Unit // factory unit doing the manufacturing
	TotalQty int   // number of factory units in the group
}

// FarmGroup is a group of farm units on a ship or colony.
type FarmGroup struct {
	CorS     *CorS             // ship or colony that controls the group
	Id       int               // unique identifier
	No       int               // group number, range 1...10
	Units    []*FarmGroupUnits // farm units in the group
	StageQty [4]int            // assumes four turns to produce a single unit
}

type FarmGroups []*FarmGroup

// FarmGroupUnits is the number of farms working together in the group
type FarmGroupUnits struct {
	Unit     *Unit // farm unit growing the food
	TotalQty int   // number of farm units in the group
}

type HullUnit struct {
	Unit     *Unit
	TotalQty int // number of units
}

type InventoryUnit struct {
	Unit      *Unit
	TotalQty  int // number of units
	StowedQty int // number of units that are disassembled for storage
}

// MineGroup is a group of mines working a single deposit.
// All mine units in a group must be the same type and tech level.
type MineGroup struct {
	CorS     *CorS // colony that controls the group
	Id       int   // unique identifier
	No       int
	Deposit  *Deposit // deposit being mined
	Unit     *Unit    // mine units in the group
	TotalQty int      // number of mine units in the group
	StageQty [4]int   // assumes four turns to produce a single unit
}

type MineGroups []*MineGroup

// Nation is a single nation in the game.
// The controller of the nation rules it, and may designate other
// players to control ships and colonies in the nation.
// These players are called viceroys or regents.
type Nation struct {
	Id                 int     // unique id for nation
	No                 int     // nation number, starts at 1
	Name               string  // unique name for this nation
	GovtName           string  // name of the government
	GovtKind           string  // kind of government
	HomePlanet         *Planet // nation's home planet
	ControlledBy       *Player // player controlling this nation
	Speciality         string  // nation's speciality for research
	TechLevel          int     // current tech level of the nation
	ResearchPointsPool int     // points in pool
	// not used currently
	Skills
}

type Pay struct {
	ProfessionalPct float64
	SoldierPct      float64
	UnskilledPct    float64
}

type Player struct {
	Id        int     // unique id for a player
	UserId    int     // user that controls this player
	Name      string  // unique name for this player
	MemberOf  *Nation // nation the player is aligned with
	ReportsTo *Player // player that this player reports to
}

// Planet is an orbit. It may be empty.
type Planet struct {
	Id             int // unique identifier
	System         *System
	Star           *Star
	OrbitNo        int    // 1...10
	Kind           string // asteroid belt, empty, gas giant, terrestrial
	HabitabilityNo int    // 0...25
	Colonies       CorSs
	Deposits       Deposits
	Ships          CorSs
}

type Population struct {
	ProfessionalQty        int
	SoldierQty             int
	UnskilledQty           int
	UnemployedQty          int
	ConstructionCrewQty    int
	SpyTeamQty             int
	RebelPct               float64
	BirthsPriorTurn        int
	NaturalDeathsPriorTurn int
}

type Rations struct {
	ProfessionalPct float64
	SoldierPct      float64
	UnskilledPct    float64
	UnemployedPct   float64
}

type Skills struct {
	Biology       int
	Bureaucracy   int
	Gravitics     int
	LifeSupport   int
	Manufacturing int
	Military      int
	Mining        int
	Shields       int
}

// Star is a stellar system in the game.
// It contains zero or more planets, with each planet assigned to an orbit ranging from 1...10
type Star struct {
	Id       int // unique identifier
	System   *System
	Sequence string // A, B, etc
	Kind     string
	Planets  []*Planet
}

type System struct {
	Id     int // unique identifier
	Coords Coordinates
	Stars  []*Star
}

// Unit is a thing in the game.
type Unit struct {
	Id                  int // unique identifier
	Kind                string
	Code                string
	TechLevel           int
	Name                string
	Description         string
	MassPerUnit         float64 // mass (in metric tonnes) of a single unit
	VolumePerUnit       float64 // volume (in cubic meters) of a single unit
	Hudnut              bool    // if true, unit can be disassembled when stowed
	StowedVolumePerUnit float64 // volume (in cubic meters) of a single unit when stowed
}
