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
	"github.com/mdhender/wraith/internal/orders"
	"log"
	"strings"
)

type PhaseOrders struct {
	Player *Player
	// orders sorted by phase
	Combat      []*orders.Order
	SetUp       []*orders.Order
	Disassembly []*orders.Order
	Retool      []*RetoolPhaseOrder
	Transfer    []*orders.Order
	Assembly    []*AssemblyPhaseOrder
	Trade       []*orders.Order
	Survey      []*orders.Order
	Espionage   []*orders.Order
	Movement    []*orders.Order
	Draft       []*orders.Order
	Pay         []*orders.Order
	Ration      []*orders.Order
	Control     []*ControlPhaseOrder
}

type AssemblyPhaseOrder struct {
	FactoryGroup *AssembleFactoryGroupOrder
	MiningGroup  *AssembleMiningGroupOrder
}
type AssembleFactoryGroupOrder struct {
	CorS     string // id of ship or colony to assemble in
	Quantity int
	Unit     string
	Product  string
}
type AssembleMiningGroupOrder struct {
	CorS     string // id of ship or colony to assemble in
	Quantity int
	Unit     string
	Product  string
}

type RetoolPhaseOrder struct {
	FactoryGroup *RetoolFactoryGroupOrder
	MiningGroup  *RetoolMiningGroupOrder
}
type RetoolFactoryGroupOrder struct {
	CorS     string // id of ship or colony to assemble in
	Quantity int
	Unit     string
	Product  string
}
type RetoolMiningGroupOrder struct {
	CorS     string // id of ship or colony to assemble in
	Quantity int
	Unit     string
	Product  string
}

type ControlPhaseOrder struct {
	ControlColony *ControlColonyOrder
	ControlShip   *ControlShipOrder
	NameColony    *NameColonyOrder
	NameShip      *NameShipOrder
}
type ControlColonyOrder struct {
	Id string // id of colony to take control of
}
type ControlShipOrder struct {
	Id string // id of ship to take control of
}
type NameColonyOrder struct {
	Id   string // id of colony to name
	Name string // name to assign to colony
}
type NameShipOrder struct {
	Id   string // id of ship to name
	Name string // name to assign to ship
}

// Execute runs all the orders in the list of phases.
// If the list is empty, no phases will run.
func (e *Engine) Execute(pos []*PhaseOrders, phases ...string) error {
	if indexOf("fuel-allocation", phases) != -1 {
		log.Printf("execute: fuel-allocation phase\n")
		e.ExecuteFuelAllocationPhase(pos)
	}
	if indexOf("life-support", phases) != -1 {
		log.Printf("execute: life-support phase\n")
		e.ExecuteLifeSupportPhase(pos)
	}
	if indexOf("farm-production", phases) != -1 {
		log.Printf("execute: farm-production phase: not implemented\n")
	}
	if indexOf("mine-production", phases) != -1 {
		log.Printf("execute: mine-production phase: not implemented\n")
	}
	if indexOf("factory-production", phases) != -1 {
		log.Printf("execute: factory-production phase: not implemented\n")
	}
	if indexOf("combat", phases) != -1 {
		log.Printf("execute: combat phase: not implemented\n")
	}
	if indexOf("setup", phases) != -1 {
		log.Printf("execute: setup phase: not implemented\n")
	}
	if indexOf("disassembly", phases) != -1 {
		log.Printf("execute: disassembly phase: not implemented\n")
	}
	if indexOf("retool", phases) != -1 {
		log.Printf("execute: retool phase\n")
		e.ExecuteRetoolPhase(pos)
	}
	if indexOf("transfer", phases) != -1 {
		log.Printf("execute: transfer phase: not implemented\n")
	}
	if indexOf("assembly", phases) != -1 {
		log.Printf("execute: assembly phase\n")
		e.ExecuteAssemblyPhase(pos)
	}
	if indexOf("trade", phases) != -1 {
		log.Printf("execute: trade phase: not implemented\n")
	}
	if indexOf("survey", phases) != -1 {
		log.Printf("execute: survey phase: not implemented\n")
	}
	if indexOf("espionage", phases) != -1 {
		log.Printf("execute: espionage phase: not implemented\n")
	}
	if indexOf("movement", phases) != -1 {
		log.Printf("execute: movement phase: not implemented\n")
	}
	if indexOf("draft", phases) != -1 {
		log.Printf("execute: draft phase: not implemented\n")
	}
	if indexOf("pay", phases) != -1 {
		log.Printf("execute: pay phase: not implemented\n")
	}
	if indexOf("ration", phases) != -1 {
		log.Printf("execute: ration phase: not implemented\n")
	}
	if indexOf("control", phases) != -1 {
		log.Printf("execute: control phase\n")
		for _, err := range e.ExecuteControlPhase(pos) {
			log.Printf("execute: control: %v\n", err)
		}
	}

	// bookkeeping
	for _, player := range e.Players {
		for _, colony := range player.Colonies {
			colony.Population.BirthsPriorTurn = 0
			colony.Population.NaturalDeathsPriorTurn = colony.nonCombatDeaths
			foundFuel := false
			for _, u := range colony.Inventory {
				if u.Unit.Kind == "fuel" {
					if foundFuel {
						u.TotalQty, u.StowedQty = 0, 0
					} else {
						foundFuel, u.TotalQty, u.StowedQty = true, colony.fuel.available, colony.fuel.available
					}
				}
			}
		}
		for _, ship := range player.Ships {
			ship.Population.BirthsPriorTurn = 0
			ship.Population.NaturalDeathsPriorTurn = ship.nonCombatDeaths
			foundFuel := false
			for _, u := range ship.Inventory {
				if u.Unit.Kind == "fuel" {
					if foundFuel {
						u.TotalQty, u.StowedQty = 0, 0
					} else {
						foundFuel, u.TotalQty, u.StowedQty = true, ship.fuel.available, ship.fuel.available
					}
				}
			}
		}
	}

	return nil
}

// ExecuteFuelAllocationPhase runs all the orders in the fuel allocation phase.
func (e *Engine) ExecuteFuelAllocationPhase(pos []*PhaseOrders) (errs []error) {
	for _, player := range e.Players {
		for _, colony := range player.Colonies {
			for _, u := range colony.Inventory {
				if u.Unit.Kind != "fuel" {
					continue
				}
				colony.fuel.available += u.TotalQty
			}
		}
		for _, ship := range player.Ships {
			for _, u := range ship.Inventory {
				if u.Unit.Kind != "fuel" {
					continue
				}
				ship.fuel.available += u.TotalQty
			}
		}
	}
	return errs
}

// ExecuteLifeSupportPhase runs all the orders in the life support phase.
func (e *Engine) ExecuteLifeSupportPhase(pos []*PhaseOrders) (errs []error) {
	for _, player := range e.Players {
		for _, colony := range player.Colonies {
			if !(colony.Kind == "enclosed" || colony.Kind == "orbital") {
				continue
			}
			lsuPopulation := 0
			for _, u := range colony.Hull {
				if colony.fuel.available <= 0 {
					break
				} else if u.Unit.Kind != "life-support" {
					continue
				}
				u.fuel.needed = u.Unit.fuelUsed(u.TotalQty)
				if colony.fuel.available < u.fuel.needed {
					u.activeQty = colony.fuel.available / u.Unit.fuelUsed(1)
					u.fuel.allocated = u.Unit.fuelUsed(u.activeQty)
				} else {
					u.activeQty = u.TotalQty
					u.fuel.allocated = u.fuel.needed
				}
				colony.fuel.available -= u.fuel.allocated
				lsuPopulation = lsuPopulation + u.activeQty*u.Unit.TechLevel*u.Unit.TechLevel
			}
			if deaths := colony.Population.Total() - lsuPopulation; deaths > 0 {
				log.Printf("execute: life-support: %q: %q: deaths %d\n", player.Name, colony.HullId, deaths)
				colony.Population.KillProportionally(deaths)
				colony.nonCombatDeaths += deaths
			}
		}
		for _, ship := range player.Ships {
			lsuPopulation := 0
			for _, u := range ship.Hull {
				if ship.fuel.available <= 0 {
					break
				} else if u.Unit.Kind != "life-support" {
					continue
				}
				u.fuel.needed = u.Unit.fuelUsed(u.TotalQty)
				if ship.fuel.available < u.fuel.needed {
					u.activeQty = ship.fuel.available / u.Unit.fuelUsed(1)
					u.fuel.allocated = u.Unit.fuelUsed(u.activeQty)
				} else {
					u.activeQty = u.TotalQty
					u.fuel.allocated = u.fuel.needed
				}
				ship.fuel.available -= u.fuel.allocated
				lsuPopulation = lsuPopulation + u.activeQty*u.Unit.TechLevel*u.Unit.TechLevel
			}
		}
	}
	return errs
}

// ExecuteAssemblyPhase runs all the orders in the assembly phase.
func (e *Engine) ExecuteAssemblyPhase(pos []*PhaseOrders) (errs []error) {
	for _, po := range pos {
		if len(po.Assembly) == 0 {
			continue
		}
		log.Printf("execute: %s: assembly\n", po.Player.Name)
	}
	return errs
}

// ExecuteControlPhase runs all the orders in the control phase.
func (e *Engine) ExecuteControlPhase(pos []*PhaseOrders) (errs []error) {
	for _, po := range pos {
		for _, order := range po.Control {
			if err := order.ControlColony.Execute(e, po.Player); err != nil {
				errs = append(errs, err)
			}
			if err := order.ControlShip.Execute(e, po.Player); err != nil {
				errs = append(errs, err)
			}
			if err := order.NameColony.Execute(e, po.Player); err != nil {
				errs = append(errs, err)
			}
			if err := order.NameShip.Execute(e, po.Player); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// Execute changes the controller of a colony to the player.
// Will fail if the colony is controlled by another player.
func (o *ControlColonyOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	log.Printf("execute: %s: control: colony %q\n", p.Name, o.Id)
	// find colony
	c, ok := e.findColony(o.Id)
	if !ok {
		return fmt.Errorf("no such colony %q", o.Id)
	}
	// fail if controlled by another player
	if c.ControlledBy != nil && c.ControlledBy != p {
		return fmt.Errorf("no such colony %q", o.Id)
	}
	// update the controller to the player
	c.ControlledBy = p
	return nil
}

// Execute changes the controller of a ship to the player.
// Will fail if the ship is controlled by another player.
func (o *ControlShipOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	log.Printf("execute: %s: control: ship %q\n", p.Name, o.Id)
	// find ship
	s, ok := e.findShip(o.Id)
	if !ok {
		return fmt.Errorf("no such ship %q", o.Id)
	}
	// fail if controlled by another player
	if s.ControlledBy != nil && s.ControlledBy != p {
		return fmt.Errorf("no such ship %q", o.Id)
	}
	// update the controller to the player
	s.ControlledBy = p
	return nil
}

// Execute changes the name of a colony controlled by a player.
// Will fail if the colony is not controlled by the player.
func (o *NameColonyOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	log.Printf("execute: %s: name: colony %q %s\n", p.Name, o.Id, o.Name)
	// find colony
	c, ok := e.findColony(o.Id)
	if !ok {
		return fmt.Errorf("no such colony %q", o.Id)
	}
	// fail if controlled by another player
	if c.ControlledBy != p {
		return fmt.Errorf("no such colony %q", o.Id)
	}
	// update the name
	c.Name = strings.Trim(o.Name, `"`)
	return nil
}

// Execute changes the name of a ship controlled by a player.
// Will fail if the ship is not controlled by the player.
func (o *NameShipOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	log.Printf("execute: %s: name: ship %q %s\n", p.Name, o.Id, o.Name)
	// find ship
	s, ok := e.findShip(o.Id)
	if !ok {
		return fmt.Errorf("no such ship %q", o.Id)
	}
	// fail if controlled by another player
	if s.ControlledBy != nil && s.ControlledBy != p {
		return fmt.Errorf("no such ship %q", o.Id)
	}
	// update the name
	s.Name = strings.Trim(o.Name, `"`)
	return nil
}

// ExecuteRetoolPhase runs all the orders in the retool phase.
func (e *Engine) ExecuteRetoolPhase(pos []*PhaseOrders) (errs []error) {
	for _, po := range pos {
		if len(po.Retool) == 0 {
			continue
		}
		log.Printf("execute: %s: retool\n", po.Player.Name)
	}
	return errs
}
