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

package server

import (
	"fmt"
	"math"
)

type Game struct {
	Id      string    `json:"id"`
	Turn    int       `json:"turn"`
	Nations []*Nation `json:"nations,omitempty"`
	Cluster *Cluster  `json:"cluster,omitempty"`
}

type Nation struct {
	Id         int    `json:"nation-id"`
	Name       string `json:"name"`
	Speciality string `json:"speciality"`
	Government struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
	} `json:"government"`
	HomePlanet struct {
		Name     string `json:"name"`
		Location struct {
			X     int `json:"x"`
			Y     int `json:"y"`
			Z     int `json:"z"`
			Star  int `json:"star,omitempty"`
			Orbit int `json:"orbit"`
		} `json:"location"`
	} `json:"home-planet"`
	Skills   Skills    `json:"skills"`
	Colonies []*Colony `json:"colonies,omitempty"`
}

type Skills struct {
	Biology       int `json:"biology"`
	Bureaucracy   int `json:"bureaucracy"`
	Gravitics     int `json:"gravitics"`
	LifeSupport   int `json:"life-support"`
	Manufacturing int `json:"manufacturing"`
	Military      int `json:"military"`
	Mining        int `json:"mining"`
	Shields       int `json:"shields"`
}

type Colony struct {
	Id       int    `json:"colony-id"`
	Kind     string `json:"kind"`
	Location struct {
		X     int `json:"x"`
		Y     int `json:"y"`
		Z     int `json:"z"`
		Star  int `json:"star,omitempty"`
		Orbit int `json:"orbit"`
	} `json:"location"`
	TechLevel  int `json:"tech-level"`
	Population struct {
		Professional      Population `json:"professional"`
		Soldier           Population `json:"soldier"`
		Unskilled         Population `json:"unskilled"`
		Unemployed        Population `json:"unemployed"`
		ConstructionCrews int        `json:"construction-crews,omitempty"`
		SpyTeams          int        `json:"spy-teams,omitempty"`
		Births            int        `json:"births,omitempty"`
		Deaths            int        `json:"deaths,omitempty"`
	} `json:"population"`
	Inventory     []*Inventory `json:"inventory,omitempty"`
	FactoryGroups []*Group     `json:"factory-groups,omitempty"`
	FarmGroups    []*Group     `json:"farm-groups,omitempty"`
	MiningGroups  []*Group     `json:"mining-groups,omitempty"`
}

func (c *Colony) TotalPay() int {
	if c == nil {
		return 0
	}
	return c.Population.Professional.TotalPay() + c.Population.Soldier.TotalPay() + c.Population.Unskilled.TotalPay() + c.Population.Unemployed.TotalPay()
}

func (c *Colony) TotalPopulation() int {
	if c == nil {
		return 0
	}
	return c.Population.Professional.TotalPopulation() + c.Population.Soldier.TotalPopulation() + c.Population.Unskilled.TotalPopulation() + c.Population.Unemployed.TotalPopulation()
}

func (c *Colony) TotalRation() int {
	if c == nil {
		return 0
	}
	return c.Population.Professional.TotalRation() + c.Population.Soldier.TotalRation() + c.Population.Unskilled.TotalRation() + c.Population.Unemployed.TotalRation()
}

type Inventory struct {
	Name           string `json:"name"`
	Code           string `json:"code,omitempty"`
	TechLevel      int    `json:"tech-level,omitempty"`
	OperationalQty int    `json:"operational-qty,omitempty"`
	StowedQty      int    `json:"stowed-qty,omitempty"`
	//MassUnits      int    `json:"mass-units,omitempty"`
	//EnclosedUnits  int    `json:"enclosed-units,omitempty"`
}

func (u Inventory) EnclosedMassUnits() int {
	mus := u.MassUnits()
	return int(math.Ceil(float64(mus) / 2))
}

func (u Inventory) MassUnits() int {
	switch u.Name {
	case "anti-missile":
		return 4 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "assault-craft":
		return 5 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "assault-weapon":
		return 2 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "automation":
		return 4 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "consumer-goods":
		return int(math.Ceil(0.6 * float64((u.OperationalQty + u.StowedQty))))
	case "energy-shield":
		return 50 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "energy-weapon":
		return 10 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "factory":
		return (12 + 2*u.TechLevel) * (u.OperationalQty + u.StowedQty)
	case "farm":
		return (6 + 2*u.TechLevel) * (u.OperationalQty + u.StowedQty)
	case "food":
		return 6 * (u.OperationalQty + u.StowedQty)
	case "fuel":
		return 1 * (u.OperationalQty + u.StowedQty)
	case "gold":
		return 1 * (u.OperationalQty + u.StowedQty)
	case "hyper-drive":
		return 45 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "life-support":
		return 8 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "light-structural":
		return int(math.Ceil(0.05 * float64((u.OperationalQty + u.StowedQty))))
	case "metallics":
		return 1 * (u.OperationalQty + u.StowedQty)
	case "military-robots":
		return (20 + 2*u.TechLevel) * (u.OperationalQty + u.StowedQty)
	case "military-supplies":
		return int(math.Ceil(0.04 * float64((u.OperationalQty + u.StowedQty))))
	case "mine":
		return (10 + 2*u.TechLevel) * (u.OperationalQty + u.StowedQty)
	case "missile":
		return 4 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "missile-launcher":
		return 25 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "non-metallics":
		return 1 * (u.OperationalQty + u.StowedQty)
	case "orbital-drive":
		return 25 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "sensor", "sensors":
		return 40 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "space-drive":
		return 25 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "star-drive":
		return 25 * u.TechLevel * (u.OperationalQty + u.StowedQty)
	case "structural":
		return int(math.Ceil(0.5 * float64((u.OperationalQty + u.StowedQty))))
	case "super-light-structural":
		return int(math.Ceil(0.005 * float64((u.OperationalQty + u.StowedQty))))
	case "transport":
		return int(math.Ceil(0.1 * float64(u.TechLevel*u.TechLevel*(u.OperationalQty+u.StowedQty))))
	default:
		panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
	}
}

type Population struct {
	Code   string  `json:"code,omitempty"`
	Qty    int     `json:"qty,omitempty"`
	Pay    float64 `json:"pay,omitempty"`
	Ration float64 `json:"ration,omitempty"`
}

// TotalPay assumes that the base rates are per unit of population
//  PROFESSIONAL      0.375 CONSUMER GOODS
//  SOLDIER           0.250 CONSUMER GOODS
//  UNSKILLED WORKER  0.125 CONSUMER GOODS
//  UNEMPLOYABLE      0.000 CONSUMER GOODS
func (p *Population) TotalPay() int {
	if p == nil {
		return 0
	}
	switch p.Code {
	case "PRO":
		return int(math.Ceil((0.375 * p.Pay) * float64(p.Qty)))
	case "SLD":
		return int(math.Ceil((0.250 * p.Pay) * float64(p.Qty)))
	case "USK":
		return int(math.Ceil((0.125 * p.Pay) * float64(p.Qty)))
	case "UEM":
		return 0
	default:
		panic(fmt.Sprintf("assert(PopUnit.Code != %q)", p.Code))
	}
}

func (p *Population) TotalPopulation() int {
	if p == nil {
		return 0
	}
	return p.Qty
}

// TotalRation assumes that base ration is 0.25 food units per unit of population
func (p *Population) TotalRation() int {
	if p == nil {
		return 0
	}
	switch p.Code {
	case "PRO":
		return int(math.Ceil((0.25 * p.Ration) * (float64(p.Qty))))
	case "SLD":
		return int(math.Ceil((0.25 * p.Ration) * (float64(p.Qty))))
	case "USK":
		return int(math.Ceil((0.25 * p.Ration) * (float64(p.Qty))))
	case "UEM":
		return int(math.Ceil((0.25 * p.Ration) * (float64(p.Qty))))
	default:
		panic(fmt.Sprintf("assert(PopUnit.Code != %q)", p.Code))
	}
}

type Group struct {
	Id    int           `json:"group-id,omitempty"`
	Name  string        `json:"name,omitempty"`
	Units []*GroupUnits `json:"units,omitempty"`
}

type GroupUnits struct {
	TechLevel int   `json:"tech-level,omitempty"`
	Qty       int   `json:"qty,omitempty"`
	Stages    []int `json:"stages,omitempty"`
}

type Cluster struct {
	Radius  int       `json:"radius"`
	Systems []*System `json:"systems"`
}

type System struct {
	Id          int     `json:"system-id,omitempty"`
	X           int     `json:"x"`
	Y           int     `json:"y"`
	Z           int     `json:"z"`
	HomeSystem  bool    `json:"home-system,omitempty"`
	Singularity bool    `json:"singularity,omitempty"`
	Ring        int     `json:"ring,omitempty"`
	Stars       []*Star `json:"stars,omitempty"`
}

type Star struct {
	Id       int       `json:"star-id,omitempty"`
	Kind     string    `json:"kind,omitempty"`
	HomeStar bool      `json:"home-star,omitempty"`
	Orbits   []*Planet `json:"orbits,omitempty"`
}

type Planet struct {
	Id                 int                `json:"planet-id,omitempty"`
	Orbit              int                `json:"orbit,omitempty"`
	Kind               string             `json:"kind"`
	HomePlanet         bool               `json:"home-planet,omitempty"`
	HabitabilityNumber int                `json:"habitability-number,omitempty"`
	Deposits           []*NaturalResource `json:"deposits,omitempty"`
}

type NaturalResource struct {
	Id                int     `json:"natural-resource-id,omitempty"`
	Kind              string  `json:"kind,omitempty"`
	Yield             float64 `json:"yield,omitempty"`
	InitialQuantity   int     `json:"initial-quantity,omitempty"`
	QuantityRemaining int     `json:"quantity-remaining,omitempty"`
}
