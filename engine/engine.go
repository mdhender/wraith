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
	"github.com/mdhender/wraith/models"
)

/*
type ReportSystem struct {
	ReportGame        *ReportGame
	Coordinates Coordinates
	Stars       []*Star
}

type Coordinates struct {
	X int
	Y int
	Z int
}

type Star struct {
	ReportSystem   *ReportSystem
	Sequence string // A, B, etc
	Kind     string
	Orbits   [11]*Planet // each orbit may or may not contain a planet
}

type Planet struct {
	Star           *Star
	OrbitNo        int    // 1..10
	Kind           string // asteroid belt, gas giant, terrestrial
	HabitabilityNo int
	ControlledBy   *ReportNation
	Deposits       []*Deposit
	Colonies       []*ReportColony
	Ships          []*ReportShip
}

type Deposit struct {
	Planet           *Planet
	ControlledBy     *ReportNation
	ReportUnit             string  // fuel, gold, metallics, non-metallics
	QtyInitial       int     // in mass units
	QtyRemaining     int     // in mass units
	MiningDifficulty float64 // how hard it is to extract each mass unit
	YieldPct         float64 // percentage of each mass unit that yields units
}

type ReportColony struct {
	Id            int
	Location      *Planet
	Kind          string // surface colony, enclosed colony, orbital colony
	TechLevel     int
	BuiltBy       *ReportNation
	ControlledBy  *ReportNation
	Inventory     []*Inventory
	MiningGroups  []*MiningGroup
	FactoryGroups []*FactoryGroup
}

type ReportShip struct {
	Id            int
	Location      *Planet
	TechLevel     int
	BuiltBy       *ReportNation
	ControlledBy  *ReportNation
	Inventory     []*Inventory
	FactoryGroups []*FactoryGroup
}

type FactoryGroup struct {
	ReportColony    *ReportColony
	ReportShip      *ReportShip
	GroupNo   int
	Inventory []*Inventory
	ReportUnit      string
	TechLevel int
}

type MiningGroup struct {
	ReportColony    *ReportColony
	GroupNo   int
	Deposit   *Deposit
	Inventory []*Inventory
}

type Inventory struct {
	ReportUnit           string
	TechLevel      int
	QtyOperational int
	QtyStowed      int
	TotalMass      int
	EnclosedMass   int
}

type ReportNation struct {
	ReportPlayer   *ReportPlayer
	Name     string
	Colonies []*ReportColony
	Ships    []*ReportShip
}

type ReportPlayer struct {
	ReportGame   *ReportGame
	User   *User
	ReportNation *ReportNation
}
*/

type Engine struct {
	r        *models.Store
	Game     *models.Game
	Colonies map[string]*Colony // key is id
	Ships    map[string]*Ship   // key is id
}

func (e *Engine) Version() string {
	return "0.1.0"
}

// reset will free up any game already in memory
func (e *Engine) reset() {
	e.Game = nil
}
