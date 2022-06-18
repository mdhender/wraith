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

type Game struct {
	Id         string
	Name       string
	TurnNumber int
	Systems    []*System
}

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
