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

// Package adapters implements functions to convert between data types
package adapters

import (
	"fmt"
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/internal/orders"
	"github.com/mdhender/wraith/internal/tokens"
	"github.com/mdhender/wraith/models"
	"github.com/mdhender/wraith/wraith"
)

func ModelsColoniesToEngineColonies(mc []*models.ColonyOrShip) []*engine.Colony {
	ec := make([]*engine.Colony, 0, len(mc))
	for _, c := range mc {
		ec = append(ec, ModelsColonyToEngineColony(c))
	}
	return ec
}

func ModelsColonyToEngineColony(mc *models.ColonyOrShip) *engine.Colony {
	ec := &engine.Colony{
		Id:           fmt.Sprintf("C%d", mc.MSN),
		ControlledBy: nil,
		Name:         "",
	}
	return ec
}

func ModelsPlayerToEnginePlayer(mp *models.Player) *engine.Player {
	var ep engine.Player
	ep.Id = mp.Id
	ep.Handle = mp.Details[0].Handle
	ep.Nation.Name = mp.MemberOf.Details[0].Name
	ep.Nation.Speciality = mp.MemberOf.Speciality
	if mp.MemberOf.HomePlanet != nil {
		ep.Nation.HomeWorld = mp.MemberOf.HomePlanet.String()
	}
	ep.Nation.GovtName = mp.MemberOf.Details[0].GovtName
	ep.Nation.GovtKind = mp.MemberOf.Details[0].GovtKind
	return &ep
}

// OrdersToEnginePhaseOrders converts orders into the Engine's expected format while splitting them into buckets for each phase.
// Because it appends to the bucket, it does not change the relative order of commands in a phase.
// Invalid or unknown orders are dropped.
func OrdersToEnginePhaseOrders(o ...*orders.Order) *engine.PhaseOrders {
	var epo engine.PhaseOrders
	for _, order := range o {
		if order == nil || order.Verb == nil || order.Errors != nil || order.Reject != nil {
			continue
		}
		switch order.Verb.Kind {
		case tokens.AssembleFactoryGroup:
			epo.Assembly = append(epo.Assembly, &engine.AssemblyPhaseOrder{FactoryGroup: &engine.AssembleFactoryGroupOrder{
				CorS:     string(order.Args[0].Text),
				Quantity: order.Args[1].Integer,
				Unit:     "",
				Product:  "",
			}})
		case tokens.AssembleMiningGroup:
			epo.Assembly = append(epo.Assembly, &engine.AssemblyPhaseOrder{MiningGroup: &engine.AssembleMiningGroupOrder{
				CorS:     string(order.Args[0].Text),
				Quantity: order.Args[1].Integer,
				Unit:     "",
				Product:  "",
			}})
		case tokens.Control:
			id := string(order.Args[0].Text)
			if order.Args[0].Kind == tokens.ColonyId {
				epo.Control = append(epo.Control, &engine.ControlPhaseOrder{ControlColony: &engine.ControlColonyOrder{Id: id}})
			} else if order.Args[0].Kind == tokens.ShipId {
				epo.Control = append(epo.Control, &engine.ControlPhaseOrder{ControlShip: &engine.ControlShipOrder{Id: id}})
			}
		case tokens.Name:
			id := string(order.Args[0].Text)
			name := string(order.Args[1].Text)
			if order.Args[0].Kind == tokens.ColonyId {
				epo.Control = append(epo.Control, &engine.ControlPhaseOrder{NameColony: &engine.NameColonyOrder{Id: id, Name: name}})
			} else if order.Args[0].Kind == tokens.ShipId {
				epo.Control = append(epo.Control, &engine.ControlPhaseOrder{NameShip: &engine.NameShipOrder{Id: id, Name: name}})
			}
		}
	}
	return &epo
}

// OrdersToPhaseOrders converts orders into the Engine's expected format while splitting them into buckets for each phase.
// Because it appends to the bucket, it does not change the relative order of commands in a phase.
// Invalid or unknown orders are dropped.
func OrdersToPhaseOrders(p *wraith.Player, o ...*orders.Order) *wraith.PhaseOrders {
	epo := &wraith.PhaseOrders{Player: p}
	for _, order := range o {
		if order == nil || order.Verb == nil || order.Errors != nil || order.Reject != nil {
			continue
		}
		switch order.Verb.Kind {
		case tokens.AssembleFactoryGroup:
			epo.Assembly = append(epo.Assembly, &wraith.AssemblyPhaseOrder{FactoryGroup: &wraith.AssembleFactoryGroupOrder{
				CorS:     string(order.Args[0].Text),
				Quantity: order.Args[1].Integer,
				Unit:     "",
				Product:  "",
			}})
		case tokens.AssembleMiningGroup:
			epo.Assembly = append(epo.Assembly, &wraith.AssemblyPhaseOrder{MiningGroup: &wraith.AssembleMiningGroupOrder{
				CorS:     string(order.Args[0].Text),
				Quantity: order.Args[1].Integer,
				Unit:     "",
				Product:  "",
			}})
		case tokens.Control:
			id := string(order.Args[0].Text)
			if order.Args[0].Kind == tokens.ColonyId {
				epo.Control = append(epo.Control, &wraith.ControlPhaseOrder{ControlColony: &wraith.ControlColonyOrder{Id: id}})
			} else if order.Args[0].Kind == tokens.ShipId {
				epo.Control = append(epo.Control, &wraith.ControlPhaseOrder{ControlShip: &wraith.ControlShipOrder{Id: id}})
			}
		case tokens.Name:
			id := string(order.Args[0].Text)
			name := string(order.Args[1].Text)
			if order.Args[0].Kind == tokens.ColonyId {
				epo.Control = append(epo.Control, &wraith.ControlPhaseOrder{NameColony: &wraith.NameColonyOrder{Id: id, Name: name}})
			} else if order.Args[0].Kind == tokens.ShipId {
				epo.Control = append(epo.Control, &wraith.ControlPhaseOrder{NameShip: &wraith.NameShipOrder{Id: id, Name: name}})
			}
		}
	}
	return epo
}
