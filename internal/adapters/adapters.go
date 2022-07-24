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

// OrdersToPhaseOrders converts orders into the Engine's expected format while splitting them into buckets for each phase.
// Because it appends to the bucket, it does not change the relative order of commands in a phase.
// Invalid or unknown orders are dropped.
func OrdersToPhaseOrders(epo *wraith.PhaseOrders, o ...*orders.Order) *wraith.PhaseOrders {
	for _, order := range o {
		if order == nil || order.Verb == nil || order.Errors != nil || order.Reject != nil {
			continue
		}
		switch order.Verb.Kind {
		case tokens.AssembleConstructionCrew:
			o := &wraith.AssembleConstructionCrewOrder{
				CorS:     string(order.Args[0].Text),
				Quantity: order.Args[1].Integer,
			}
			epo.Assembly = append(epo.Assembly, &wraith.AssemblyPhaseOrder{ConstructionCrew: o})
		case tokens.AssembleFactoryGroup:
			o := &wraith.AssembleFactoryGroupOrder{
				CorS:     string(order.Args[0].Text),
				Quantity: order.Args[1].Integer,
				Unit:     order.Args[2].String(),
				Product:  order.Args[3].String(),
			}
			epo.Assembly = append(epo.Assembly, &wraith.AssemblyPhaseOrder{FactoryGroup: o})
		case tokens.AssembleFarmGroup:
			o := &wraith.AssembleFarmGroupOrder{
				CorS:     string(order.Args[0].Text),
				Quantity: order.Args[1].Integer,
				Unit:     order.Args[2].String(),
				Product:  order.Args[3].String(),
			}
			epo.Assembly = append(epo.Assembly, &wraith.AssemblyPhaseOrder{FarmGroup: o})
		case tokens.AssembleMineGroup:
			epo.Assembly = append(epo.Assembly, &wraith.AssemblyPhaseOrder{MiningGroup: &wraith.AssembleMineGroupOrder{
				CorS:     string(order.Args[0].Text),
				Quantity: order.Args[1].Integer,
				Unit:     order.Args[2].String(),
				Deposit:  order.Args[3].String(),
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
