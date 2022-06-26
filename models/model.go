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

package models

import (
	"time"
)

// User is the person who owns the account used to access the game
type User struct {
	Id           int    // unique identifier
	Handle       string // unique handle, forced to lower-case
	HashedSecret string
	Profiles     []*UserProfile
}

// UserProfile is the handle and email of the user
type UserProfile struct {
	User   *User
	EffDt  time.Time // moment record becomes effective
	EndDt  time.Time // moment record ceases to be effective
	Email  string    // email, forced to lower-case
	Handle string    // display version of handle
}

// Game is a single game
type Game struct {
	Id          int // unique identifier
	ShortName   string
	Name        string
	Description string
	CurrentTurn *Turn
	Nations     []*Nation
	Players     []*Player
	Systems     []*System
	Turns       []*Turn
}

// Turn is a single turn in the game
type Turn struct {
	No      int       // 0+
	Year    int       // 1...9999
	Quarter int       // 1...4
	StartDt time.Time // moment turn starts
	EndDt   time.Time // moment turn ends
}

// Player is a User's position in a Game
type Player struct {
	Game         *Game
	Id           int
	ControlledBy *User
	SubjectOf    *Player // set if the player is a viceroy or regent
	Details      []*PlayerDetail
}

// PlayerDetail is the handle of the player in the game
type PlayerDetail struct {
	Player  *Player
	EffTurn *Turn // turn record becomes active
	EndTurn *Turn // turn record ceases to be active
	Handle  string
}

// Nation is a single nation in the game.
// The "ruler" of the nation controls it, but may create viceroys
// or regents to control ships and colonies in the nation.
type Nation struct {
	Game        *Game
	Id          int // unique identifier
	No          int // nation number in the game, starts at 1
	Speciality  string
	Description string
	Details     []*NationDetail
	Skills      []*NationSkills
	Player      *Player // player currently controlling the nation
	Colonies    []*ColonyOrShip
	Ships       []*ColonyOrShip
}

// NationDetail is items that can change value during the game.
type NationDetail struct {
	Nation       *Nation
	EffTurn      *Turn // turn record becomes active
	EndTurn      *Turn // turn record ceases to be active
	Name         string
	GovtName     string
	GovtKind     string
	ControlledBy *Player
}

// NationSkills are the skills and tech levels of the nation.
type NationSkills struct {
	Nation             *Nation
	EffTurn            *Turn // turn record becomes active
	EndTurn            *Turn // turn record ceases to be active
	ResearchPointsPool int   //
	Biology            int   // not used currently
	Bureaucracy        int   // not used currently
	Gravitics          int   // not used currently
	LifeSupport        int   // not used currently
	Manufacturing      int   // not used currently
	Military           int   // not used currently
	Mining             int   // not used currently
	Shields            int   // not used currently
}

// NationResearch is the tech level and research level of the nation.
type NationResearch struct {
	Nation             *Nation
	EffTurn            *Turn // turn record becomes active
	EndTurn            *Turn // turn record ceases to be active
	TechLevel          int   //
	ResearchPointsPool int   //
}

// System is a system in the game.
// It contains zero or more stars.
type System struct {
	Game     *Game
	Id       int // unique identifier
	X        int
	Y        int
	Z        int
	QtyStars int
}

// Star is a stellar system in the game.
// It contains zero or more planets, with each planet assigned to an orbit ranging from 1...10
type Star struct {
	System   *System
	Id       int    // unique identifier
	Sequence string // A, B, etc
	Kind     string
	Orbits   [11]*Planet // each orbit may or may not contain a planet
}

// Planet is a non-empty orbit.
type Planet struct {
	Star           *Star
	Id             int    // unique identifier
	OrbitNo        int    // 1..10
	Kind           string // asteroid belt, gas giant, terrestrial
	HabitabilityNo int
	HomePlanet     bool
	Deposits       []*NaturalResource
	Colonies       []*ColonyOrShip // colonies on or orbiting the planet
	Ships          []*ColonyOrShip // ships orbiting the planet
}

// PlanetDetail contains items that change from turn to turn
type PlanetDetail struct {
	Planet       *Planet
	EffTurn      *Turn   // turn record becomes active
	EndTurn      *Turn   // turn record ceases to be active
	ControlledBy *Nation // nation currently controlling the planet
}

// NaturalResource is a deposit of fuel, gold, metal, or non-metals on a planet
type NaturalResource struct {
	Planet     *Planet
	Id         int     // unique identifier
	DepositNo  int     // number of deposit on planet
	Kind       string  // fuel, gold, metallics, non-metallics
	QtyInitial int     // in mass units
	YieldPct   float64 // percentage of each mass unit that yields units
	Details    []*NaturalResourceDetail
}

// NaturalResourceDetail contains items that change from turn to turn
type NaturalResourceDetail struct {
	NaturalResource *NaturalResource
	EffTurn         *Turn   // turn record becomes active
	EndTurn         *Turn   // turn record ceases to be active
	ControlledBy    *Nation // nation currently controlling the resource
	QtyRemaining    int     // in mass units
}

// ColonyOrShip is either a colony or a ship.
// Ships may change orbits; colonies may not.
type ColonyOrShip struct {
	Id         int    // unique identifier
	MSN        int    // manufacturer serial number; in game id for the colony or ship
	Kind       string // surface colony, enclosed colony, orbital colony, ship
	BuiltBy    *Nation
	Details    []*CSDetail
	Hull       []*CSHull
	Inventory  []*CSInventory
	Population []*CSPopulation
	Rations    []*CSRations
	Pay        []*CSPay
	Factories  []*FactoryGroup
	Farms      []*FarmGroup
	Mines      []*MiningGroup
}

// CSDetail contains items that may change from turn to turn
type CSDetail struct {
	CS           *ColonyOrShip
	EffTurn      *Turn // turn record becomes active
	EndTurn      *Turn // turn record ceases to be active
	Location     *Planet
	Name         string
	TechLevel    int
	ControlledBy *Player
}

// CSHull is the infrastructure and components of the ship or colony.
// Colonies are not allowed to add engines or drives to their hull.
type CSHull struct {
	CS              *ColonyOrShip
	EffTurn         *Turn // turn record becomes active
	EndTurn         *Turn // turn record ceases to be active
	Unit            *Unit
	QtyOperational  int
	MassOperational int
	TotalMass       int
}

// CSInventory is the contents of the ship or colony.
type CSInventory struct {
	CS              *ColonyOrShip
	EffTurn         *Turn // turn record becomes active
	EndTurn         *Turn // turn record ceases to be active
	Unit            *Unit
	QtyOperational  int
	MassOperational int
	QtyStowed       int
	MassStowed      int
	TotalMass       int
	EnclosedMass    int
}

// CSPay is the pay rate for the ship or colony.
type CSPay struct {
	CS              *ColonyOrShip
	EffTurn         *Turn // turn record becomes active
	EndTurn         *Turn // turn record ceases to be active
	ProfessionalPct float64
	SoldierPct      float64
	UnskilledPct    float64
	UnemployedPct   float64
}

// CSPopulation is the population of the ship or colony.
type CSPopulation struct {
	CS                  *ColonyOrShip
	EffTurn             *Turn // turn record becomes active
	EndTurn             *Turn // turn record ceases to be active
	QtyProfessional     int
	QtySoldier          int
	QtyUnskilled        int
	QtyUnemployed       int
	QtyConstructionCrew int
	QtySpyTeam          int
	RebelPct            float64
}

// CSPay is the rations rate for the ship or colony.
type CSRations struct {
	CS              *ColonyOrShip
	EffTurn         *Turn // turn record becomes active
	EndTurn         *Turn // turn record ceases to be active
	ProfessionalPct float64
	SoldierPct      float64
	UnskilledPct    float64
	UnemployedPct   float64
}

// FactoryGroup is a group of factories on the ship or colony.
// Each group is dedicated to manufacturing one type of unit.
type FactoryGroup struct {
	CS      *ColonyOrShip
	Id      int // unique identifier
	GroupNo int
	EffTurn *Turn // turn record becomes active
	EndTurn *Turn // turn record ceases to be active
	Unit    *Unit // unit being manufactured
}

// FactoryGroupUnits is the number of factories working together in the group
type FactoryGroupUnits struct {
	Group          *FactoryGroup
	EffTurn        *Turn // turn record becomes active
	EndTurn        *Turn // turn record ceases to be active
	Unit           *Unit // unit being manufactured
	QtyOperational int
}

// FactoryGroupStages is the number units in each stage of the group
type FactoryGroupStages struct {
	Group     *FactoryGroup
	Turn      *Turn
	QtyStage1 int
	QtyStage2 int
	QtyStage3 int
	QtyStage4 int
}

// FarmGroup is a group of farms on the ship or colony.
type FarmGroup struct {
	CS      *ColonyOrShip
	Id      int // unique identifier
	GroupNo int
	EffTurn *Turn // turn record becomes active
	EndTurn *Turn // turn record ceases to be active
}

// FarmGroupUnits is the number of farms working together in the group
type FarmGroupUnits struct {
	Group          *FactoryGroup
	EffTurn        *Turn // turn record becomes active
	EndTurn        *Turn // turn record ceases to be active
	Unit           *Unit // farm unit
	QtyOperational int
}

// FarmGroupStages is the number units in each stage of the group
type FarmGroupStages struct {
	Group     *FarmGroup
	Turn      *Turn
	QtyStage1 int
	QtyStage2 int
	QtyStage3 int
	QtyStage4 int
}

// MiningGroup is a group of mines working a single deposit.
type MiningGroup struct {
	CS      *ColonyOrShip
	Id      int // unique identifier
	GroupNo int
	EffTurn *Turn // turn record becomes active
	EndTurn *Turn // turn record ceases to be active
	Deposit *NaturalResource
}

// MiningGroupUnits is the number of mines working together in the group
type MiningGroupUnits struct {
	Group          *MiningGroup
	EffTurn        *Turn // turn record becomes active
	EndTurn        *Turn // turn record ceases to be active
	Unit           *Unit // unit mining
	QtyOperational int
}

// MiningGroupStages is the number units in each stage of the group
type MiningGroupStages struct {
	Group     *MiningGroup
	Turn      *Turn
	QtyStage1 int
	QtyStage2 int
	QtyStage3 int
	QtyStage4 int
}

// Unit is a thing in the game.
type Unit struct {
	Id           int // unique identifier
	Code         string
	TechLevel    int
	Name         string
	Description  string
	Hudnut       bool    // true if unit can be disassembled for storage
	Mass         float64 // mass (in mass units) of a single unit
	EnclosedMass float64 // volume (in enclosed mass units) of a single unit
}
