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
	"golang.org/x/text/message"
	"io"
	"log"
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
	Colonies        map[string]*CorS
	CorSById        map[int]*CorS
	Deposits        map[int]*Deposit
	FactoryGroups   map[int]*FactoryGroup
	FarmGroups      map[int]*FarmGroup
	MineGroups      map[int]*MineGroup
	Nations         map[int]*Nation
	Planets         map[int]*Planet
	Players         map[int]*Player
	Ships           map[string]*CorS
	Stars           map[int]*Star
	Systems         map[int]*System
	Units           map[int]*Unit
	UnitsFromString map[string]*Unit
	Seq             int
}

func (e *Engine) NextSeq() int {
	e.Seq++
	return e.Seq
}

// CorS is a colony or ship
type CorS struct {
	Id                                  int     // unique identifier
	Kind                                string  // orbital, ship, or surface
	HullId                              string  // S or C + MSN
	MSN                                 int     // manufacturer serial number; in game id for the colony or ship
	BuiltBy                             *Nation // nation that originally built the colony or ship
	Name                                string  // name of this colony or ship
	TechLevel                           int     // tech level of this colony or ship
	ControlledBy                        *Player // player that controls this colony or ship
	Planet                              *Planet // planet the colony or ship is located at
	Hull                                InventoryUnits
	Inventory                           InventoryUnits
	Population                          Population
	Pay                                 Pay
	Rations                             Rations
	FactoryGroups                       FactoryGroups // list of the factory groups
	FarmGroups                          FarmGroups    // list of the farm groups
	MineGroups                          MineGroups    // list of the mine groups
	fuel, pro, sol, uns, uem, cons, spy requisition
	lifeSupportCapacity                 int
	nonCombatDeaths                     int
}

func (cs *CorS) InitializeInventory() {
	for _, u := range cs.Hull {
		u.operational.initial, u.operational.allocated, u.operational.created, u.operational.destroyed = u.ActiveQty, 0, 0, 0
		u.stowed.initial, u.stowed.allocated, u.operational.created, u.stowed.destroyed = u.StowedQty, 0, 0, 0
	}
	for _, u := range cs.Inventory {
		u.operational.initial, u.operational.allocated, u.stowed.created, u.operational.destroyed = u.ActiveQty, 0, 0, 0
		u.stowed.initial, u.stowed.allocated, u.stowed.created, u.stowed.destroyed = u.StowedQty, 0, 0, 0
	}
}

func (cs *CorS) lifeSupportCheck() {
	if !(cs.Kind == "enclosed" || cs.Kind == "orbital" || cs.Kind == "ship") {
		return
	}
	var playerName string
	if cs.ControlledBy != nil {
		playerName = cs.ControlledBy.Name
	} else {
		playerName = "nobody"
	}
	if totalPop(cs) <= cs.lifeSupportCapacity {
		return
	}
	deaths := totalPop(cs) - cs.lifeSupportCapacity
	log.Printf("execute: life-support: %q: %q: deaths %d\n", playerName, cs.HullId, deaths)
	killProportionally(cs, deaths)
	cs.nonCombatDeaths += deaths
}

func (cs *CorS) lifeSupportInitialization(pos []*PhaseOrders) {
	cs.lifeSupportCapacity = 0
	if !(cs.Kind == "enclosed" || cs.Kind == "orbital" || cs.Kind == "ship") {
		return
	}

	// find fuel in inventory.
	// note that fuel must always be stowed.
	var fuel *InventoryUnit
	for _, u := range cs.Inventory {
		if u.Unit.Kind == "fuel" {
			fuel = u
			break
		}
	}

	for _, u := range cs.Hull {
		if u.Unit.Kind != "life-support" {
			continue
		} else if fuel.stowed.available() == 0 {
			cs.Log("  **** no fuel available for life support!\n")
			break
		}
		// allocate fuel to the life support unit
		fuelNeeded := int(math.Ceil(u.Unit.FuelPerUnitPerTurn * float64(u.operational.available())))
		if fuelNeeded == 0 {
			continue
		}
		fuelAllocated := fuel.stowed.allocate(fuelNeeded)
		lsuActivated := int(float64(fuelAllocated) / u.Unit.FuelPerUnitPerTurn)
		u.operational.allocate(lsuActivated)

		// capacity is number of units times the unit's tech level squared
		cs.lifeSupportCapacity += lsuActivated * u.Unit.TechLevel * u.Unit.TechLevel

		cs.Log("  %13d %-20s  %13d activated  %13d fuel allocated\n",
			u.operational.initial, u.Unit.Name, u.operational.allocated, fuelAllocated)
	}

	cs.Log("  %13d population capacity\n", cs.lifeSupportCapacity)
}

func (cs *CorS) Log(format string, args ...interface{}) {
	cs.ControlledBy.Log(format, args...)
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
	CorS     *CorS          // ship or colony that controls the group
	Id       int            // unique identifier
	No       int            // group number, range 1...255
	Product  *Unit          // unit being produced by the group
	Units    InventoryUnits // units assigned to the group
	StageQty [4]int         // assumes four turns to produce a single unit
}

type FactoryGroups []*FactoryGroup

func (f FactoryGroups) Len() int {
	return len(f)
}

func (f FactoryGroups) Less(i, j int) bool {
	return f[i].No < f[j].No
}

func (f FactoryGroups) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// FactoryGroupUnits is the number of factories working together in the group
type FactoryGroupUnits struct {
	Unit     *Unit // factory unit doing the manufacturing
	TotalQty int   // number of factory units in the group
}

// FarmGroup is a group of farm units on a ship or colony.
type FarmGroup struct {
	CorS     *CorS          // ship or colony that controls the group
	Id       int            // unique identifier
	No       int            // group number, range 1...10
	Product  *Unit          // unit being produced by the group
	Units    InventoryUnits // units assigned to the group
	StageQty [4]int         // assumes four turns to produce a single unit
}

type FarmGroups []*FarmGroup

func (f FarmGroups) Len() int {
	return len(f)
}

func (f FarmGroups) Less(i, j int) bool {
	return f[i].No < f[j].No
}

func (f FarmGroups) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type requisition struct {
	initial   int
	allocated int
	created   int
	destroyed int

	needed, operational int
}

type InventoryUnit struct {
	Unit                 *Unit
	ActiveQty            int // number of active/operational units
	StowedQty            int // number of units that are disassembled for storage
	activeQty            int
	fuel, pro, uns       requisition
	available, allocated int
	operational, stowed  requisition
}

// allocate activates as many units as it can with the remaining available units.
// it returns the amount actually allocated.
func (r *requisition) allocate(qty int) int {
	if r.available() < qty {
		qty = r.available()
	}
	r.allocated += qty
	return qty
}

// available returns the number of units available.
// assumes that destroyed units update the allocated amount, too.
func (r *requisition) available() (qty int) {
	if qty = r.initial - r.destroyed - r.allocated; qty < 0 {
		qty = 0
	}
	return qty
}

// create adds new units.
// they will not be available until after the bookkeeping phase.
func (r *requisition) create(qty int) {
	r.created += qty
}

// destroy flags units as being destroyed.
// it updates both the destroyed and allocated quantities.
// destroyed units will be removed from inventory during the bookkeeping phase.
// it returns the number of units that were destroyed, which will be
// less than the requested amount when there aren't enough units left.
func (r *requisition) destroy(qty int) int {
	// should never happen, but check for zero inventory anyway
	if r.initial == 0 {
		r.destroyed, r.allocated = 0, 0
		return 0
	}
	// remaining is the number of un-destroyed units
	remaining := r.initial - r.destroyed
	// if we are destroying all the remaining units, flag them as destroyed and allocated.
	if remaining <= qty {
		r.destroyed, r.allocated = r.initial, 0
		return remaining
	}
	// if we've already allocated all the units, update only the destroyed quantity
	if r.initial <= r.allocated {
		if r.allocated -= qty; r.allocated < 0 {
			r.allocated = 0
		}
		if r.destroyed += qty; r.initial < r.destroyed {
			r.destroyed = r.initial
		}
		return qty
	}
	// otherwise allocate the quantity between allocated and un-allocated units.
	// we do this by moving some of the previously allocated units to destroyed.
	// when doing so, ensure that we don't exceed the initial quantity.
	pctDestroyed := float64(qty) / float64(r.initial)
	amtAllocatedDestroyed := int(math.Ceil(pctDestroyed * float64(r.allocated)))
	if r.allocated -= amtAllocatedDestroyed; r.allocated < 0 {
		r.allocated = 0
	}
	if r.destroyed += qty; r.initial < r.destroyed {
		r.destroyed = r.initial
	}
	return qty
}

//// allocateFuel activates as many units as it can given the amount of fuel available.
//// it returns the amount actually allocated.
//func (u *InventoryUnit) allocateFuel(fuelAvailable int) int {
//	u.fuel.needed, u.fuel.allocated, u.fuel.used = 0, 0, 0
//
//	u.fuel.needed = int(math.Ceil(u.Unit.FuelPerUnitPerTurn * float64(u.ActiveQty)))
//	if u.fuel.needed == 0 {
//		// nothing to do
//		return 0
//	} else if fuelAvailable < u.fuel.needed {
//		// activate as many units as we can
//		u.activeQty = int(float64(fuelAvailable) / u.Unit.FuelPerUnitPerTurn)
//		u.fuel.allocated = int(math.Ceil(u.Unit.FuelPerUnitPerTurn * float64(u.ActiveQty)))
//	} else {
//		u.activeQty = u.ActiveQty
//		u.fuel.allocated = u.fuel.needed
//	}
//
//	return u.fuel.allocated
//}

func (u *InventoryUnit) totalMass() int {
	return int(math.Ceil(float64(u.ActiveQty+u.StowedQty) * u.Unit.MassPerUnit))
}

func (u *InventoryUnit) totalVolume() int {
	return int(math.Ceil(float64(u.ActiveQty)*u.Unit.VolumePerUnit)) + int(math.Ceil(float64(u.StowedQty)*u.Unit.StowedVolumePerUnit))
}

type InventoryUnits []*InventoryUnit

func (u InventoryUnits) Len() int {
	return len(u)
}

func (u InventoryUnits) Less(i, j int) bool {
	return u[i].Unit.Id < u[j].Unit.Id
}

func (u InventoryUnits) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

// MineGroup is a group of mines working a single deposit.
// All mine units in a group must be the same type and tech level.
type MineGroup struct {
	CorS     *CorS // colony that controls the group
	Id       int   // unique identifier
	No       int
	Deposit  *Deposit       // deposit being mined
	Unit     *InventoryUnit // mine units in the group
	StageQty [4]int         // assumes four turns to produce a single unit
}

type MineGroups []*MineGroup

func (m MineGroups) Len() int {
	return len(m)
}

func (m MineGroups) Less(i, j int) bool {
	return m[i].No < m[j].No
}

func (m MineGroups) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

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
	Logger    struct {
		MP *message.Printer
		W  io.Writer
	}
}

func (p *Player) Log(format string, args ...interface{}) {
	if p != nil && p.Logger.MP != nil && p.Logger.W != nil {
		_, _ = p.Logger.MP.Fprintf(p.Logger.W, format, args...)
	}
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

type population struct {
	professional        int
	soldier             int
	unskilled           int
	unemployed          int
	construction        int
	spy                 int
	lifeSupportCapacity int
	nonCombatDeaths     int
}

func (p *population) total() int {
	return p.professional + p.soldier + p.unskilled + p.unemployed + 2*p.construction + 2*p.spy
}

func (p Population) Total() int {
	return p.ProfessionalQty + p.SoldierQty + p.UnskilledQty + p.UnemployedQty
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

func (u *Unit) fuelUsed(qty int) int {
	return int(math.Ceil(u.FuelPerUnitPerTurn * float64(qty)))
}
