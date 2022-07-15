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
	Assembly    []*orders.Order
	Trade       []*orders.Order
	Survey      []*orders.Order
	Espionage   []*orders.Order
	Movement    []*orders.Order
	Draft       []*orders.Order
	Pay         []*orders.Order
	Ration      []*orders.Order
	Control     []*ControlPhaseOrder
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
		for _, po := range pos {
			if len(po.Retool) == 0 {
				continue
			}
			log.Printf("execute: %s: %s\n", "retool", po.Player.Handle)
		}
	}
	if indexOf("transfer", phases) != -1 {
		// not yet implemented
	}
	if indexOf("assembly", phases) != -1 {
		// not yet implemented
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
		for _, po := range pos {
			if len(po.Control) == 0 {
				continue
			}
			log.Printf("execute: %s: %s\n", "control", po.Player.Handle)
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
