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
	"fmt"
	"github.com/mdhender/wraith/internal/orders"
	"log"
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
	if indexOf("combat", phases) != -1 {
		// not yet implemented
	}
	if indexOf("setup", phases) != -1 {
		// not yet implemented
	}
	if indexOf("disassembly", phases) != -1 {
		// not yet implemented
	}
	if indexOf("retool", phases) != -1 {
		e.ExecuteRetoolPhase(pos)
	}
	if indexOf("transfer", phases) != -1 {
		// not yet implemented
	}
	if indexOf("assembly", phases) != -1 {
		e.ExecuteAssemblyPhase(pos)
	}
	if indexOf("trade", phases) != -1 {
		// not yet implemented
	}
	if indexOf("survey", phases) != -1 {
		// not yet implemented
	}
	if indexOf("espionage", phases) != -1 {
		// not yet implemented
	}
	if indexOf("movement", phases) != -1 {
		// not yet implemented
	}
	if indexOf("draft", phases) != -1 {
		// not yet implemented
	}
	if indexOf("pay", phases) != -1 {
		// not yet implemented
	}
	if indexOf("ration", phases) != -1 {
		// not yet implemented
	}
	if indexOf("control", phases) != -1 {
		for _, err := range e.ExecuteControlPhase(pos) {
			log.Printf("execute: control: %v\n", err)
		}
	}
	return nil
}

func indexOf(s string, sl []string) int {
	for i, p := range sl {
		if s == p {
			return i
		}
	}
	return -1
}

// ExecuteAssemblyPhase runs all the orders in the assembly phase.
func (e *Engine) ExecuteAssemblyPhase(pos []*PhaseOrders) []error {
	var errs []error
	for _, po := range pos {
		if len(po.Assembly) == 0 {
			continue
		}
		log.Printf("execute: %s: assembly\n", po.Player.Handle)
	}
	return errs
}

// ExecuteControlPhase runs all the orders in the control phase.
func (e *Engine) ExecuteControlPhase(pos []*PhaseOrders) []error {
	var errs []error
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
	log.Printf("execute: %s: control: colony %q\n", p.Handle, o.Id)
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
	log.Printf("execute: %s: control: ship %q\n", p.Handle, o.Id)
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
	log.Printf("execute: %s: name: colony %q %s\n", p.Handle, o.Id, o.Name)
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
	c.Name = o.Name
	return nil
}

// Execute changes the name of a ship controlled by a player.
// Will fail if the ship is not controlled by the player.
func (o *NameShipOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	log.Printf("execute: %s: name: ship %q %s\n", p.Handle, o.Id, o.Name)
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
	s.Name = o.Name
	return nil
}

// ExecuteRetoolPhase runs all the orders in the retool phase.
func (e *Engine) ExecuteRetoolPhase(pos []*PhaseOrders) []error {
	var errs []error
	for _, po := range pos {
		if len(po.Retool) == 0 {
			continue
		}
		log.Printf("execute: %s: retool\n", po.Player.Handle)
	}
	return errs
}
