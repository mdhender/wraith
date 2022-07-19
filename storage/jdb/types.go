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

package jdb

// Coordinates are the x,y,z coordinates of a system
type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// Deposit of fuel, gold, metal, or non-metals on the surface of a planet
type Deposit struct {
	Id                   int     `json:"id"`                                // unique identifier
	PlanetId             int     `json:"planet-id"`                         // id of planet deposit is on
	No                   int     `json:"no"`                                // number of deposit on planet
	UnitId               int     `json:"unit-id"`                           // fuel, gold, metallic, non-metallic
	InitialQty           int     `json:"initial-qty"`                       // in metric tonnes
	RemainingQty         int     `json:"remaining-qty"`                     // in metric tonnes
	YieldPct             float64 `json:"yield-pct"`                         // percentage of each mass unit that yields units
	ControlledByColonyId int     `json:"controlled-by-colony-id,omitempty"` // colony controlling this deposit
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

// EnclosedColony defines an enclosed surface colony.
type EnclosedColony struct {
	Id                   int            `json:"id"`                                // unique identifier
	MSN                  int            `json:"msn"`                               // manufacturer serial number; in game id for the colony
	BuiltByNationId      int            `json:"built-by-nation-id,omitempty"`      // id of the nation that originally built the colony
	Name                 string         `json:"name,omitempty"`                    // name of this colony
	TechLevel            int            `json:"tech-level,omitempty"`              // tech level of this colony
	ControlledByPlayerId int            `json:"controlled-by-player-id,omitempty"` // id of player that controls this colony
	PlanetId             int            `json:"planet-id,omitempty"`               // id of planet the colony is built on
	Hull                 HullUnits      `json:"hull,omitempty"`
	Inventory            InventoryUnits `json:"inventory,omitempty"`
	Population           struct {
		ProfessionalQty        int     `json:"professional-qty,omitempty"`
		SoldierQty             int     `json:"soldier-qty,omitempty"`
		UnskilledQty           int     `json:"unskilled-qty,omitempty"`
		UnemployedQty          int     `json:"unemployed-qty,omitempty"`
		ConstructionCrewQty    int     `json:"construction-crew-qty,omitempty"`
		SpyTeamQty             int     `json:"spy-team-qty,omitempty"`
		RebelPct               float64 `json:"rebel-pct,omitempty"`
		BirthsPriorTurn        int     `json:"births-prior-turn,omitempty"`
		NaturalDeathsPriorTurn int     `json:"natural-deaths-prior-turn,omitempty"`
	} `json:"population"`
	Pay struct {
		ProfessionalPct float64 `json:"professional-pct,omitempty"`
		SoldierPct      float64 `json:"soldier-pct,omitempty"`
		UnskilledPct    float64 `json:"unskilled-pct,omitempty"`
	} `json:"pay"`
	Rations struct {
		ProfessionalPct float64 `json:"professional-pct,omitempty"`
		SoldierPct      float64 `json:"soldier-pct,omitempty"`
		UnskilledPct    float64 `json:"unskilled-pct,omitempty"`
		UnemployedPct   float64 `json:"unemployed-pct,omitempty"`
	} `json:"rations"`
	FactoryGroupIds []int `json:"factory-group-ids,omitempty"` // list of the factory group ids
	FarmGroupIds    []int `json:"farm-group-ids,omitempty"`    // list of the farm group ids
	MineGroupIds    []int `json:"mine-group-ids,omitempty"`    // list of the mine group ids
}

type EnclosedColonies []*EnclosedColony

func (e EnclosedColonies) Len() int {
	return len(e)
}

func (e EnclosedColonies) Less(i, j int) bool {
	return e[i].Id < e[j].Id
}

func (e EnclosedColonies) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// FactoryGroup is a group of factories on the ship or colony.
// Each group is dedicated to manufacturing one type of unit.
type FactoryGroup struct {
	Id        int                  `json:"id"`              // unique identifier
	CorSId    int                  `json:"cors-id"`         // id for the ship or colony that controls the group
	No        int                  `json:"no"`              // group number, range 1...255
	Product   int                  `json:"product"`         // unit being manufactured
	Units     []*FactoryGroupUnits `json:"units,omitempty"` // factory units in the group
	Stage1Qty int                  `json:"stage-1-qty,omitempty"`
	Stage2Qty int                  `json:"stage-2-qty,omitempty"`
	Stage3Qty int                  `json:"stage-3-qty,omitempty"`
	Stage4Qty int                  `json:"stage-4-qty,omitempty"`
}

type FactoryGroups []*FactoryGroup

func (f FactoryGroups) Len() int {
	return len(f)
}

func (f FactoryGroups) Less(i, j int) bool {
	if f[i].CorSId < f[j].CorSId {
		return true
	} else if f[i].CorSId == f[j].CorSId {
		return f[i].No < f[j].No
	}
	return false
}

func (f FactoryGroups) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// FactoryGroupUnits is the number of factories working together in the group
type FactoryGroupUnits struct {
	UnitId   int `json:"unit-id"` // factory unit doing the manufacturing
	TotalQty int `json:"total-qty,omitempty"`
}

// FarmGroup is a group of farm units on the ship or colony.
type FarmGroup struct {
	Id        int               `json:"id"`              // unique identifier
	CorSId    int               `json:"cors-id"`         // id for the ship or colony that controls the group
	No        int               `json:"no"`              // group number, range 1...10
	Units     []*FarmGroupUnits `json:"units,omitempty"` // farm units in the group
	Stage1Qty int               `json:"stage-1-qty,omitempty"`
	Stage2Qty int               `json:"stage-2-qty,omitempty"`
	Stage3Qty int               `json:"stage-3-qty,omitempty"`
	Stage4Qty int               `json:"stage-4-qty,omitempty"`
}

type FarmGroups []*FarmGroup

func (f FarmGroups) Len() int {
	return len(f)
}

func (f FarmGroups) Less(i, j int) bool {
	if f[i].CorSId < f[j].CorSId {
		return true
	} else if f[i].CorSId == f[j].CorSId {
		return f[i].No < f[j].No
	}
	return false
}

func (f FarmGroups) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// FarmGroupUnits is the number of farms working together in the group
type FarmGroupUnits struct {
	UnitId   int `json:"unit-id"` // farm unit doing the manufacturing
	TotalQty int `json:"total-qty,omitempty"`
}

// Game contains the information about the game being played.
type Game struct {
	Id        int    `json:"id"`   // unique identifier for game
	Name      string `json:"name"` // full name of game
	ShortName string `json:"short-name"`
	Turn      struct {
		Year    int    `json:"year"`              // 1...9999
		Quarter int    `json:"quarter"`           // 1...4
		StartDt string `json:"startDt,omitempty"` // moment turn starts, UTC
		EndDt   string `json:"endDt,omitempty"`   // moment just after turn ends, UTC
	} `json:"turn"`

	Deposits         Deposits         `json:"deposits,omitempty"`
	EnclosedColonies EnclosedColonies `json:"enclosed-colonies,omitempty"`
	FactoryGroups    FactoryGroups    `json:"factory-groups,omitempty"`
	FarmGroups       FarmGroups       `json:"farm-groups,omitempty"`
	MineGroups       MineGroups       `json:"mine-groups,omitempty"`
	Nations          Nations          `json:"nations,omitempty"`
	OrbitalColonies  OrbitalColonies  `json:"orbital-colonies,omitempty"`
	Planets          Planets          `json:"planets,omitempty"`
	Players          Players          `json:"players,omitempty"`
	Ships            Ships            `json:"ships,omitempty"`
	SurfaceColonies  SurfaceColonies  `json:"surface-colonies,omitempty"`
	Units            Units            `json:"units,omitempty"`
	Stars            Stars            `json:"stars,omitempty"`
	Systems          Systems          `json:"systems,omitempty"`
}

type HullUnit struct {
	UnitId   int `json:"unit-id"`             // id of unit
	TotalQty int `json:"total-qty,omitempty"` // number of units
}

type HullUnits []*HullUnit

func (u HullUnits) Len() int {
	return len(u)
}

func (u HullUnits) Less(i, j int) bool {
	return u[i].UnitId < u[j].UnitId
}

func (u HullUnits) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type InventoryUnit struct {
	UnitId    int `json:"unit-id"`              // id of unit
	TotalQty  int `json:"total-qty,omitempty"`  // number of units
	StowedQty int `json:"stowed-qty,omitempty"` // number of units that are disassembled for storage
}

type InventoryUnits []*InventoryUnit

func (u InventoryUnits) Len() int {
	return len(u)
}

func (u InventoryUnits) Less(i, j int) bool {
	return u[i].UnitId < u[j].UnitId
}

func (u InventoryUnits) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

// MineGroup is a group of mines working a single deposit.
// All mine units in a group must be the same type and tech level.
type MineGroup struct {
	Id        int `json:"id"`        // unique identifier
	ColonyId  int `json:"colony-id"` // id for the colony that controls the group
	No        int `json:"no"`
	DepositId int `json:"deposit-id"`
	UnitId    int `json:"unit-id"`             // mine units in the group
	TotalQty  int `json:"total-qty,omitempty"` // number of mine units in the group
	Stage1Qty int `json:"stage-1-qty,omitempty"`
	Stage2Qty int `json:"stage-2-qty,omitempty"`
	Stage3Qty int `json:"stage-3-qty,omitempty"`
	Stage4Qty int `json:"stage-4-qty,omitempty"`
}

type MineGroups []*MineGroup

func (m MineGroups) Len() int {
	return len(m)
}

func (m MineGroups) Less(i, j int) bool {
	if m[i].ColonyId < m[j].ColonyId {
		return true
	} else if m[i].ColonyId == m[j].ColonyId {
		return m[i].No < m[j].No
	}
	return false
}

func (m MineGroups) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Nation is a single nation in the game.
// The controller of the nation rules it, and may designate other
// players to control ships and colonies in the nation.
// These players are called viceroys or regents.
type Nation struct {
	Id                   int      `json:"id"`                      // unique id for nation
	No                   int      `json:"no"`                      // nation number, starts at 1
	Name                 string   `json:"name"`                    // unique name for this nation
	GovtName             string   `json:"govt-name"`               // name of the government
	GovtKind             string   `json:"govt-kind"`               // kind of government
	HomePlanetId         int      `json:"home-planet-id"`          // id of nation's home planet
	ControlledByPlayerId int      `json:"controlled-by-player-id"` // id of player controlling this nation
	Speciality           string   `json:"speciality"`              // nation's speciality for research
	TechLevel            int      `json:"tech-level"`              // current tech level of the nation
	ResearchPointsPool   int      `json:"research-points-pool"`    // points in pool
	Skills               struct { // not used currently
		Biology       int `json:"biology,omitempty"`       // not used currently
		Bureaucracy   int `json:"bureaucracy,omitempty"`   // not used currently
		Gravitics     int `json:"gravitics,omitempty"`     // not used currently
		LifeSupport   int `json:"life-support,omitempty"`  // not used currently
		Manufacturing int `json:"manufacturing,omitempty"` // not used currently
		Military      int `json:"military,omitempty"`      // not used currently
		Mining        int `json:"mining,omitempty"`        // not used currently
		Shields       int `json:"shields,omitempty"`       // not used currently
	} `json:"skills"` // not used currently
}

type Nations []*Nation

// Len implements sort.Sort interface
func (n Nations) Len() int {
	return len(n)
}

// Less implements sort.Sort interface
func (n Nations) Less(i, j int) bool {
	return n[i].No < n[j].No
}

// Swap implements sort.Sort interface
func (n Nations) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// OrbitalColony defines an orbital colony.
type OrbitalColony struct {
	Id                   int            `json:"id"`                                // unique identifier
	MSN                  int            `json:"msn"`                               // manufacturer serial number; in game id for the colony
	BuiltByNationId      int            `json:"built-by-nation-id,omitempty"`      // id of the nation that originally built the colony
	Name                 string         `json:"name,omitempty"`                    // name of this colony
	TechLevel            int            `json:"tech-level,omitempty"`              // tech level of this colony
	ControlledByPlayerId int            `json:"controlled-by-player-id,omitempty"` // id of player that controls this colony
	PlanetId             int            `json:"planet-id,omitempty"`               // id of planet the colony is orbiting
	Hull                 HullUnits      `json:"hull,omitempty"`
	Inventory            InventoryUnits `json:"inventory,omitempty"`
	Population           struct {
		ProfessionalQty        int     `json:"professional-qty,omitempty"`
		SoldierQty             int     `json:"soldier-qty,omitempty"`
		UnskilledQty           int     `json:"unskilled-qty,omitempty"`
		UnemployedQty          int     `json:"unemployed-qty,omitempty"`
		ConstructionCrewQty    int     `json:"construction-crew-qty,omitempty"`
		SpyTeamQty             int     `json:"spy-team-qty,omitempty"`
		RebelPct               float64 `json:"rebel-pct,omitempty"`
		BirthsPriorTurn        int     `json:"births-prior-turn,omitempty"`
		NaturalDeathsPriorTurn int     `json:"natural-deaths-prior-turn,omitempty"`
	} `json:"population"`
	Pay struct {
		ProfessionalPct float64 `json:"professional-pct,omitempty"`
		SoldierPct      float64 `json:"soldier-pct,omitempty"`
		UnskilledPct    float64 `json:"unskilled-pct,omitempty"`
	} `json:"pay"`
	Rations struct {
		ProfessionalPct float64 `json:"professional-pct,omitempty"`
		SoldierPct      float64 `json:"soldier-pct,omitempty"`
		UnskilledPct    float64 `json:"unskilled-pct,omitempty"`
		UnemployedPct   float64 `json:"unemployed-pct,omitempty"`
	} `json:"rations"`
	FactoryGroupIds []int `json:"factory-group-ids,omitempty"` // list of the factory group ids
	FarmGroupIds    []int `json:"farm-group-ids,omitempty"`    // list of the farm group ids
}

type OrbitalColonies []*OrbitalColony

func (o OrbitalColonies) Len() int {
	return len(o)
}

func (o OrbitalColonies) Less(i, j int) bool {
	return o[i].Id < o[j].Id
}

func (o OrbitalColonies) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Planet is an orbit. It may be empty.
type Planet struct {
	Id                int    `json:"id"` // unique identifier
	SystemId          int    `json:"system-id"`
	StarId            int    `json:"star-id"`
	OrbitNo           int    `json:"orbitNo"` // 1..10
	Kind              string `json:"kind"`    // asteroid belt, gas giant, terrestrial
	HabitabilityNo    int    `json:"habitabilityNo,omitempty"`
	DepositIds        []int  `json:"deposit-ids,omitempty"`
	SurfaceColonyIds  []int  `json:"surface-colony-ids,omitempty"`
	EnclosedColonyIds []int  `json:"enclosed-colony-ids,omitempty"`
	OrbitalColonyIds  []int  `json:"orbital-colony-ids,omitempty"`
	ShipIds           []int  `json:"ship-ids,omitempty"`
}

type Planets []*Planet

func (p Planets) Len() int {
	return len(p)
}

func (p Planets) Less(i, j int) bool {
	return p[i].Id < p[j].Id
}

func (p Planets) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Player is a position in the game.
type Player struct {
	Id                int    `json:"id"`                          // unique id for a player, starts at 1
	UserId            int    `json:"user-id"`                     // user that controls this player
	Name              string `json:"name"`                        // unique name for this player
	MemberOf          int    `json:"member-of"`                   // nation the player is aligned with
	ReportsToPlayerId int    `json:"reports-to-player,omitempty"` // player that this player reports to
}

type Players []*Player

// Len implements sort.Sort interface
func (p Players) Len() int {
	return len(p)
}

// Less implements sort.Sort interface
func (p Players) Less(i, j int) bool {
	return p[i].Id < p[j].Id
}

// Swap implements sort.Sort interface
func (p Players) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Ship defines a ship.
type Ship struct {
	Id                   int            `json:"id"`                                // unique identifier
	MSN                  int            `json:"msn"`                               // manufacturer serial number; in game id for the ship
	BuiltByNationId      int            `json:"built-by-nation-id,omitempty"`      // id of the nation that originally built the ship
	Name                 string         `json:"name,omitempty"`                    // name of this ship
	TechLevel            int            `json:"tech-level,omitempty"`              // tech level of this ship
	ControlledByPlayerId int            `json:"controlled-by-player-id,omitempty"` // id of player that controls this ship
	PlanetId             int            `json:"planet-id,omitempty"`               // id of planet the ship is orbiting
	Hull                 HullUnits      `json:"hull,omitempty"`
	Inventory            InventoryUnits `json:"inventory,omitempty"`
	Population           struct {
		ProfessionalQty        int     `json:"professional-qty,omitempty"`
		SoldierQty             int     `json:"soldier-qty,omitempty"`
		UnskilledQty           int     `json:"unskilled-qty,omitempty"`
		UnemployedQty          int     `json:"unemployed-qty,omitempty"`
		ConstructionCrewQty    int     `json:"construction-crew-qty,omitempty"`
		SpyTeamQty             int     `json:"spy-team-qty,omitempty"`
		RebelPct               float64 `json:"rebel-pct,omitempty"`
		NaturalDeathsPriorTurn int     `json:"natural-deaths-prior-turn,omitempty"`
	} `json:"population"`
	Pay struct {
		ProfessionalPct float64 `json:"professional-pct,omitempty"`
		SoldierPct      float64 `json:"soldier-pct,omitempty"`
		UnskilledPct    float64 `json:"unskilled-pct,omitempty"`
		UnemployedPct   float64 `json:"unemployed-pct,omitempty"`
	} `json:"pay"`
	Rations struct {
		ProfessionalPct float64 `json:"professional-pct,omitempty"`
		SoldierPct      float64 `json:"soldier-pct,omitempty"`
		UnskilledPct    float64 `json:"unskilled-pct,omitempty"`
		UnemployedPct   float64 `json:"unemployed-pct,omitempty"`
	} `json:"rations"`
	FactoryGroupIds []int `json:"factory-group-ids,omitempty"` // list of the factory group ids
	FarmGroupIds    []int `json:"farm-group-ids,omitempty"`    // list of the farm group ids
}

type Ships []*Ship

func (s Ships) Len() int {
	return len(s)
}

func (s Ships) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

func (s Ships) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Star is a stellar system in the game.
// It contains zero or more planets, with each planet assigned to an orbit ranging from 1...10
type Star struct {
	Id        int    `json:"id"` // unique identifier
	SystemId  int    `json:"system-id"`
	Sequence  string `json:"sequence"` // A, B, etc
	Kind      string `json:"kind"`
	PlanetIds []int  `json:"planet-ids"`
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

// SurfaceColony defines a surface colony.
type SurfaceColony struct {
	Id                   int            `json:"id"`                                // unique identifier
	MSN                  int            `json:"msn"`                               // manufacturer serial number; in game id for the colony
	BuiltByNationId      int            `json:"built-by-nation-id,omitempty"`      // id of the nation that originally built the colony
	Name                 string         `json:"name,omitempty"`                    // name of this colony
	TechLevel            int            `json:"tech-level,omitempty"`              // tech level of this colony
	ControlledByPlayerId int            `json:"controlled-by-player-id,omitempty"` // id of player that controls this colony
	PlanetId             int            `json:"planet-id,omitempty"`               // id of planet the colony is built on
	Hull                 HullUnits      `json:"hull,omitempty"`
	Inventory            InventoryUnits `json:"inventory,omitempty"`
	Population           struct {
		ProfessionalQty        int     `json:"professional-qty,omitempty"`
		SoldierQty             int     `json:"soldier-qty,omitempty"`
		UnskilledQty           int     `json:"unskilled-qty,omitempty"`
		UnemployedQty          int     `json:"unemployed-qty,omitempty"`
		ConstructionCrewQty    int     `json:"construction-crew-qty,omitempty"`
		SpyTeamQty             int     `json:"spy-team-qty,omitempty"`
		RebelPct               float64 `json:"rebel-pct,omitempty"`
		BirthsPriorTurn        int     `json:"births-prior-turn,omitempty"`
		NaturalDeathsPriorTurn int     `json:"natural-deaths-prior-turn,omitempty"`
	} `json:"population"`
	Pay struct {
		ProfessionalPct float64 `json:"professional-pct,omitempty"`
		SoldierPct      float64 `json:"soldier-pct,omitempty"`
		UnskilledPct    float64 `json:"unskilled-pct,omitempty"`
	} `json:"pay"`
	Rations struct {
		ProfessionalPct float64 `json:"professional-pct,omitempty"`
		SoldierPct      float64 `json:"soldier-pct,omitempty"`
		UnskilledPct    float64 `json:"unskilled-pct,omitempty"`
		UnemployedPct   float64 `json:"unemployed-pct,omitempty"`
	} `json:"rations"`
	FactoryGroupIds []int `json:"factory-group-ids,omitempty"` // list of the factory group ids
	FarmGroupIds    []int `json:"farm-group-ids,omitempty"`    // list of the farm group ids
	MineGroupIds    []int `json:"mine-group-ids,omitempty"`    // list of the mine group ids
}

type SurfaceColonies []*SurfaceColony

func (s SurfaceColonies) Len() int {
	return len(s)
}

func (s SurfaceColonies) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

func (s SurfaceColonies) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// System is a system in the game.
// It contains zero or more stars.
type System struct {
	Id      int         `json:"id,omitempty"` // unique identifier
	Coords  Coordinates `json:"coords"`
	StarIds []int       `json:"star-ids,omitempty"`
}

type Systems []*System

func (s Systems) Len() int {
	return len(s)
}

func (s Systems) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

func (s Systems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Unit is a thing in the game.
type Unit struct {
	Id                  int     `json:"id"` // unique identifier
	Kind                string  `json:"kind"`
	Code                string  `json:"code"`
	TechLevel           int     `json:"tech-level,omitempty"`
	Name                string  `json:"name"`
	Description         string  `json:"description,omitempty"`
	MassPerUnit         float64 `json:"mass-per-unit"`          // mass (in metric tonnes) of a single unit
	VolumePerUnit       float64 `json:"volume-per-unit"`        // volume (in cubic meters) of a single unit
	Hudnut              bool    `json:"hudnut,omitempty"`       // if true, unit can be disassembled when stowed
	StowedVolumePerUnit float64 `json:"stowed-volume-per-unit"` // volume (in cubic meters) of a single unit when stowed
	FuelPerUnitPerTurn  float64 `json:"fuel-per-unit-per-turn,omitempty"`
	MetsPerUnit         float64 `json:"mets-per-unit,omitempty"`
	NonMetsPerUnit      float64 `json:"non-mets-per-unit,omitempty"`
}

type Units []*Unit

func (u Units) Len() int {
	return len(u)
}

func (u Units) Less(i, j int) bool {
	if u[i].Kind < u[j].Kind {
		return true
	} else if u[i].Kind == u[j].Kind {
		return u[i].TechLevel < u[j].TechLevel
	}
	return false
}

func (u Units) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
