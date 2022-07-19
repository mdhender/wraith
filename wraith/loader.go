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

import (
	"fmt"
	"math"
	"time"
)

type Engine struct {
	Version string
	Game    struct {
		Id   int
		Code string
		Name string
		Turn struct {
			Year    int
			Quarter int
			StartDt time.Time
			EndDt   time.Time
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

func (u *HullUnit) totalMass() int {
	return int(math.Ceil(float64(u.TotalQty) * u.Unit.MassPerUnit))
}

func (u *HullUnit) totalVolume() int {
	return int(math.Ceil(float64(u.TotalQty) * u.Unit.VolumePerUnit))
}

type InventoryUnit struct {
	Unit      *Unit
	TotalQty  int // number of units
	StowedQty int // number of units that are disassembled for storage
}

func (u *InventoryUnit) totalMass() int {
	return int(math.Ceil(float64(u.TotalQty) * u.Unit.MassPerUnit))
}

func (u *InventoryUnit) totalVolume() int {
	return int(math.Ceil(float64(u.TotalQty-u.StowedQty)*u.Unit.VolumePerUnit)) + int(math.Ceil(float64(u.StowedQty)*u.Unit.StowedVolumePerUnit))
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

// totalPay assumes that the base rates are per unit of population
//  PROFESSIONAL      0.375 CONSUMER GOODS
//  SOLDIER           0.250 CONSUMER GOODS
//  UNSKILLED WORKER  0.125 CONSUMER GOODS
//  UNEMPLOYABLE      0.000 CONSUMER GOODS
func (pay Pay) totalPay(pop Population, code string) int {
	switch code {
	case "PRO":
		return int(math.Ceil((0.375 * pay.ProfessionalPct) * float64(pop.ProfessionalQty)))
	case "SLD":
		return int(math.Ceil((0.250 * pay.SoldierPct) * float64(pop.SoldierQty)))
	case "USK":
		return int(math.Ceil((0.125 * pay.UnskilledPct) * float64(pop.UnskilledQty)))
	case "UEM":
		return 0
	default:
		panic(fmt.Sprintf("assert(pay.totalPay.Code != %q)", code))
	}
}

type Player struct {
	Id        int     // unique id for a player
	UserId    int     // user that controls this player
	Name      string  // unique name for this player
	MemberOf  *Nation // nation the player is aligned with
	ReportsTo *Player // player that this player reports to
	Colonies  CorSs   // colonies controlled by this player
	Ships     CorSs   // ships controlled by this player
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

// totalRations assumes that base rates are per unit of population
//  PROFESSIONAL      0.250 FOOD
//  SOLDIER           0.250 FOOD
//  UNSKILLED WORKER  0.250 FOOD
//  UNEMPLOYABLE      0.250 FOOD
func (ration Rations) totalRations(pop Population, code string) int {
	switch code {
	case "PRO":
		return int(math.Ceil((0.25 * ration.ProfessionalPct) * (float64(pop.ProfessionalQty))))
	case "SLD":
		return int(math.Ceil((0.25 * ration.SoldierPct) * (float64(pop.SoldierQty))))
	case "USK":
		return int(math.Ceil((0.25 * ration.UnskilledPct) * (float64(pop.UnskilledQty))))
	case "UEM":
		return int(math.Ceil((0.25 * ration.UnemployedPct) * (float64(pop.UnemployedQty))))
	default:
		panic(fmt.Sprintf("assert(ration.totalRations.Code != %q)", code))
	}
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

type Stars []*Star

func (s Stars) Len() int {
	return len(s)
}

func (s Stars) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

func (s Stars) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type System struct {
	Id     int // unique identifier
	Coords Coordinates
	Stars  []*Star
}

// Unit is a thing in the game.
type Unit struct {
	Id                    int // unique identifier
	Kind                  string
	Code                  string
	TechLevel             int
	Name                  string
	Description           string
	MassPerUnit           float64 // mass (in metric tonnes) of a single unit
	VolumePerUnit         float64 // volume (in cubic meters) of a single unit
	Hudnut                bool    // if true, unit can be disassembled when stowed
	StowedVolumePerUnit   float64 // volume (in cubic meters) of a single unit when stowed
	FuelPerUnitPerTurn    float64
	MetsPerUnitPerTurn    float64
	NonMetsPerUnitPerTurn float64
}
