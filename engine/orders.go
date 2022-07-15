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
	"github.com/mdhender/wraith/internal/tokens"
	"log"
)

type PlayerOrders struct {
	Player *Player
	// orders sorted by phase
	Combat      []*orders.Order
	SetUp       []*orders.Order
	Disassembly []*orders.Order
	Retool      []*orders.Order
	Transfer    []*orders.Order
	Assembly    []*orders.Order
	Trade       []*orders.Order
	Survey      []*orders.Order
	Espionage   []*orders.Order
	Movement    []*orders.Order
	Draft       []*orders.Order
	Pay         []*orders.Order
	Ration      []*orders.Order
	Control     []*orders.Order
}

// PlayerOrders creates and initializes the struct.
// It sorts the orders into buckets for each phase.
// Because it appends to the bucket, it does not change the relative order of commands in a phase.
// Invalid or unknown orders are dropped.
func (e *Engine) PlayerOrders(p *Player, o []*orders.Order) *PlayerOrders {
	po := &PlayerOrders{Player: p}
	for _, order := range o {
		if order == nil || order.Verb == nil {
			continue
		}
		switch order.Verb.Kind {
		case tokens.AssembleFactoryGroup:
			po.Retool = append(po.Retool, order)
		case tokens.AssembleMiningGroup:
			po.Retool = append(po.Retool, order)
		case tokens.Control:
			po.Control = append(po.Control, order)
		case tokens.Name:
			po.Control = append(po.Control, order)
		}
	}
	return po
}

// Execute runs all the orders in the list of phases.
// If the list is empty, no phases will run.
func (e *Engine) Execute(pos []*PlayerOrders, phases ...string) error {
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
