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
	"sort"
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
	CorS       string // id of ship or colony to assemble in
	Quantity   int
	Unit       string
	Product    string
	cons, fuel requisition
}
type AssembleMiningGroupOrder struct {
	CorS     string // id of ship or colony to assemble in
	Quantity int
	Unit     string
	Group    string
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
		for _, err := range e.ExecuteFuelAllocationPhase(pos) {
			log.Printf("execute: fuel-allocation: %v\n", err)
		}
	}
	if indexOf("labor-allocation", phases) != -1 {
		log.Printf("execute: labor-allocation\n")
		for _, err := range e.ExecuteLaborAllocationPhase(pos) {
			log.Printf("execute: labor-allocation: %v\n", err)
		}
	}
	if indexOf("life-support", phases) != -1 {
		log.Printf("execute: life-support phase\n")
		for _, err := range e.ExecuteLifeSupportPhase(pos) {
			log.Printf("execute: life-support: %v\n", err)
		}
	}
	if indexOf("farm-production", phases) != -1 {
		log.Printf("execute: farm-production\n")
		for _, err := range e.ExecuteFarmProductionPhase(pos) {
			log.Printf("execute: farm-production: %v\n", err)
		}
	}
	if indexOf("mine-production", phases) != -1 {
		log.Printf("execute: mine-production phase\n")
		for _, err := range e.ExecuteMineProductionPhase(pos) {
			log.Printf("execute: mine-production: %v\n", err)
		}
	}
	if indexOf("factory-production", phases) != -1 {
		log.Printf("execute: factory-production phase\n")
		for _, err := range e.ExecuteFactoryProductionPhase(pos) {
			log.Printf("execute: factory-production: %v\n", err)
		}
	}
	if indexOf("combat", phases) != -1 {
		log.Printf("execute: combat phase\n")
		for _, err := range e.ExecuteCombatPhase(pos) {
			log.Printf("execute: combat: %v\n", err)
		}
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
		for _, err := range e.ExecuteAssemblyPhase(pos) {
			log.Printf("execute: assembly: %v\n", err)
		}
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
	for _, cs := range e.CorSById {
		// population changes
		if cs.Kind == "ship" {
			cs.Population.BirthsPriorTurn = 0
		} else {
			// TODO: create a standard of living metric and change rate to 0.25% to 2.5%
			birthRate := 0.0025 // 0.25% per year baseline
			cs.Population.BirthsPriorTurn = int(float64(totalPop(cs)) * birthRate / 4)
		}
		cs.Population.ProfessionalQty = cs.pro.available
		cs.Population.SoldierQty = cs.sol.available
		cs.Population.UnskilledQty = cs.uns.available
		cs.Population.UnemployedQty = cs.uem.available + cs.Population.BirthsPriorTurn
		cs.Population.NaturalDeathsPriorTurn = cs.nonCombatDeaths

		// update fuel depot
		foundFuel := false
		for _, u := range cs.Inventory {
			if u.Unit.Kind == "fuel" {
				if foundFuel {
					u.ActiveQty, u.StowedQty = 0, 0
				} else {
					foundFuel, u.ActiveQty, u.StowedQty = true, 0, cs.fuel.available
				}
			}
		}

		// inventory changes
		for _, group := range cs.FarmGroups {
			var unit *InventoryUnit
			for _, u := range cs.Inventory {
				if u.Unit.Id == group.Product.Id {
					unit = u
					break
				}
			}
			if unit == nil {
				unit = &InventoryUnit{Unit: group.Product}
				cs.Inventory = append(cs.Inventory, unit)
				sort.Sort(cs.Inventory)
			}
			unit.StowedQty += group.StageQty[3]
			group.StageQty[3] = 0
		}
		for _, group := range cs.MineGroups {
			var unit *InventoryUnit
			for _, u := range cs.Inventory {
				if u.Unit.Id == group.Deposit.Product.Id {
					unit = u
					break
				}
			}
			if unit == nil {
				unit = &InventoryUnit{Unit: group.Deposit.Product}
				cs.Inventory = append(cs.Inventory, unit)
				sort.Sort(cs.Inventory)
			}
			unit.StowedQty += group.StageQty[3]
			group.StageQty[3] = 0
		}
		for _, group := range cs.FactoryGroups {
			var unit *InventoryUnit
			for _, u := range cs.Inventory {
				if u.Unit.Id == group.Product.Id {
					unit = u
					break
				}
			}
			if unit == nil {
				unit = &InventoryUnit{Unit: group.Product}
				cs.Inventory = append(cs.Inventory, unit)
				sort.Sort(cs.Inventory)
			}
			unit.StowedQty += group.StageQty[3]
			group.StageQty[3] = 0

			//for _, u := range group.Units {
			//	log.Printf("%-6s: fg: %d %s: unit %s %d %d\n",
			//		cs.HullId, group.No, group.Product.Code,
			//		u.Unit.Code, u.ActiveQty, u.activeQty)
			//	u.ActiveQty = 995876
			//}
		}
	}

	return nil
}

// ExecuteFuelAllocationPhase runs all the orders in the fuel allocation phase.
func (e *Engine) ExecuteFuelAllocationPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nFuel Allocation -------------------------------------------------\n")
	}
	for _, cs := range e.CorSById {
		fuelInitialization(cs, pos)
	}
	return errs
}

// ExecuteLifeSupportPhase runs all the orders in the life support phase.
func (e *Engine) ExecuteLifeSupportPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nLife Support ----------------------------------------------------\n")
	}
	for _, cors := range e.CorSById {
		cors.lifeSupportInitialization(pos)
	}
	for _, cors := range e.CorSById {
		cors.lifeSupportCheck()
	}
	return errs
}

// ExecuteLaborAllocationPhase runs all the orders in the fuel allocation phase.
func (e *Engine) ExecuteLaborAllocationPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nLabor Allocation ------------------------------------------------\n")
	}
	for _, cs := range e.CorSById {
		laborInitialization(cs, pos)
	}
	return errs
}

// ExecuteFarmProductionPhase runs all the orders in the farm production phase.
func (e *Engine) ExecuteFarmProductionPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nFarm Production -------------------------------------------------\n")
	}
	for _, cs := range e.CorSById {
		if len(cs.FarmGroups) != 0 {
			farmProduction(cs, pos)
		}
	}
	return errs
}

// ExecuteMineProductionPhase runs all the orders in the mine production phase.
func (e *Engine) ExecuteMineProductionPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nMine Production -------------------------------------------------\n")
	}
	for _, cs := range e.CorSById {
		if len(cs.MineGroups) != 0 {
			mineProduction(cs, pos)
		}
	}
	return errs
}

// ExecuteFactoryProductionPhase runs all the orders in the factory production phase.
func (e *Engine) ExecuteFactoryProductionPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nFactory Production ----------------------------------------------\n")
	}
	for _, cs := range e.CorSById {
		if len(cs.FactoryGroups) != 0 {
			factoryProduction(cs, pos)
		}
	}
	return errs
}

// ExecuteCombatPhase runs all the orders in the combat phase.
func (e *Engine) ExecuteCombatPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nCombat ----------------------------------------------------------\n")
		o.Player.Log("  Not Implemented!\n")
	}
	return errs
}

// ExecuteAssemblyPhase runs all the orders in the assembly phase.
func (e *Engine) ExecuteAssemblyPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nAssembly --------------------------------------------------------\n")
		for _, order := range o.Assembly {
			if err := order.FactoryGroup.Execute(e, o.Player); err != nil {
				errs = append(errs, err)
			}
			if err := order.MiningGroup.Execute(e, o.Player); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// Execute assembles a factory group on the colony or ship.
// Will fail if the colony or ship is not controlled by the player.
func (o *AssembleFactoryGroupOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}

	// find colony or ship
	cs, ok := e.findColony(o.CorS)
	if !ok {
		if cs, ok = e.findShip(o.CorS); !ok {
			p.Log("  assemble %s: no such colony or ship\n", o.CorS)
			return fmt.Errorf("no such colony or ship %q", o.CorS)
		}
	}
	// fail if controlled by another player
	if cs.ControlledBy != nil && cs.ControlledBy != p {
		p.Log("  assemble %s: no such colony or ship\n", o.CorS)
		return fmt.Errorf("no such colony or ship %q", o.CorS)
	}

	// find the factory and product in the order
	factory, ok := unitFromString(e, o.Unit)
	if !ok {
		p.Log("  assemble %s: no such unit %q\n", o.CorS, o.Unit)
		return fmt.Errorf("no such unit %q", o.Unit)
	} else if factory.TechLevel > cs.TechLevel {
		p.Log("  assemble %s: unit %q: invalid tech level\n", o.CorS, o.Unit)
		return fmt.Errorf("invalid tech level %q", o.Product)
	}
	product, ok := unitFromString(e, o.Product)
	if !ok {
		p.Log("  assemble %s: no such unit %q product %q\n", o.CorS, o.Unit, o.Product)
		return fmt.Errorf("no such unit %q", o.Product)
	} else if product.TechLevel > cs.TechLevel {
		p.Log("  assemble %s: unit %q product %q: invalid tech level\n", o.CorS, o.Unit, o.Product)
		return fmt.Errorf("invalid tech level %q", o.Product)
	}

	// is there a group already producing this product or must we create a new group?
	var fg *FactoryGroup
	for _, group := range cs.FactoryGroups {
		if group.Product.Code == product.Code {
			// an existing group
			fg = group
			break
		}
	}
	if fg == nil {
		fg = &FactoryGroup{
			CorS:    cs,
			Id:      e.NextSeq(),
			No:      0,
			Product: product,
		}
		var idx [30]bool
		for _, group := range cs.FactoryGroups {
			idx[group.No] = true
		}
		for no := 1; fg.No == 0 && no < 30; no++ {
			if !idx[no] {
				fg.No = no
			}
		}
		if fg.No == 0 {
			p.Log("  assemble %s: unit %q product %q: no factory groups available", o.CorS, o.Unit, o.Product)
			return fmt.Errorf("no factory groups available")
		}
		cs.FactoryGroups = append(cs.FactoryGroups, fg)
		sort.Sort(cs.FactoryGroups)
	}
	p.Log("  assemble %s: group %2d: %-11s product %s\n", o.CorS, fg.No, factory.Name, product.Name)

	// no fuel required to assemble factory units!
	o.fuel.needed = 0

	// allocate labor. 1 CON per 500 tonnes.
	o.cons.needed = o.Quantity
	if availableCon(cs) < o.cons.needed {
		o.cons.allocated = availableCon(cs)
	} else {
		o.cons.allocated = o.cons.needed
	}
	cs.cons.allocated += o.cons.allocated

	// are there already units in the group or must we add them?
	var unit *InventoryUnit
	for _, u := range fg.Units {
		if u.Unit.Code == factory.Code {
			unit = u
			break
		}
	}
	if unit == nil {
		unit = &InventoryUnit{Unit: factory}
		fg.Units = append(fg.Units, unit)
		sort.Sort(fg.Units)
	}
	unit.ActiveQty += o.Quantity

	return nil
}

// Execute assembles a mine group on the colony or ship.
// Will fail if the colony or ship is not controlled by the player.
func (o *AssembleMiningGroupOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	mine, ok := unitFromString(e, o.Unit)
	if !ok {
		p.Log("  assemble %s: no such unit %q\n", o.CorS, o.Unit)
		return fmt.Errorf("no such unit %q", o.Unit)
	}
	//product, ok := unitFromString(e, o.Product)
	//if !ok {
	//	p.Log("  assemble %s: no such unit %q %q:\n%v\n", o.CorS, o.Unit, o.Product, o)
	//	return fmt.Errorf("no such unit %q", o.Product)
	//}
	p.Log("  assemble %s: mine %q group %q\n", o.CorS, mine.Code, o.Group)

	// find colony or ship
	var cs *CorS
	if cs, ok = e.findColony(o.CorS); !ok {
		if cs, ok = e.findShip(o.CorS); !ok {
			p.Log("  assemble %s: no such colony or ship\n", o.CorS)
			return fmt.Errorf("no such colony or ship %q", o.CorS)
		}
	}
	// fail if controlled by another player
	if cs.ControlledBy != nil && cs.ControlledBy != p {
		p.Log("  assemble %s: no such colony or ship\n", o.CorS)
		return fmt.Errorf("no such colony or ship %q", o.CorS)
	}
	return nil
}

// ExecuteControlPhase runs all the orders in the control phase.
func (e *Engine) ExecuteControlPhase(pos []*PhaseOrders) (errs []error) {
	for _, o := range pos {
		o.Player.Log("\n\nControl ---------------------------------------------------------\n")
		for _, order := range o.Control {
			if err := order.ControlColony.Execute(e, o.Player); err != nil {
				errs = append(errs, err)
			}
			if err := order.ControlShip.Execute(e, o.Player); err != nil {
				errs = append(errs, err)
			}
			if err := order.NameColony.Execute(e, o.Player); err != nil {
				errs = append(errs, err)
			}
			if err := order.NameShip.Execute(e, o.Player); err != nil {
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
	// find colony
	c, ok := e.findColony(o.Id)
	if !ok {
		p.Log("  control %s: no such colony\n", o.Id)
		return fmt.Errorf("no such colony %q", o.Id)
	}
	// fail if controlled by another player
	if c.ControlledBy != nil && c.ControlledBy != p {
		p.Log("  control %s: no such colony\n", o.Id)
		return fmt.Errorf("no such colony %q", o.Id)
	}
	// update the controller to the player
	c.ControlledBy = p
	p.Log("  control %s: now controlled by %d\n", o.Id, p.Id)
	return nil
}

// Execute changes the controller of a ship to the player.
// Will fail if the ship is controlled by another player.
func (o *ControlShipOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	p.Log("execute: %s: control: ship %q\n", p.Name, o.Id)
	// find ship
	s, ok := e.findShip(o.Id)
	if !ok {
		p.Log("  control %s: no such ship\n", o.Id)
		return fmt.Errorf("no such ship %q", o.Id)
	}
	// fail if controlled by another player
	if s.ControlledBy != nil && s.ControlledBy != p {
		p.Log("  control %s: no such ship\n", o.Id)
		return fmt.Errorf("no such ship %q", o.Id)
	}
	// update the controller to the player
	s.ControlledBy = p
	p.Log("  control %s: now controlled by %d\n", o.Id, p.Id)
	return nil
}

// Execute changes the name of a colony controlled by a player.
// Will fail if the colony is not controlled by the player.
func (o *NameColonyOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	// find colony
	c, ok := e.findColony(o.Id)
	if !ok {
		p.Log("  name %s: no such colony\n", o.Id)
		return fmt.Errorf("no such colony %q", o.Id)
	}
	// fail if controlled by another player
	if c.ControlledBy != p {
		p.Log("  name %s: no such colony\n", o.Id)
		return fmt.Errorf("no such colony %q", o.Id)
	}
	// update the name
	c.Name = strings.Trim(o.Name, `"`)
	p.Log("  name %s: now named %q\n", o.Id, c.Name)
	return nil
}

// Execute changes the name of a ship controlled by a player.
// Will fail if the ship is not controlled by the player.
func (o *NameShipOrder) Execute(e *Engine, p *Player) error {
	if o == nil {
		return nil
	}
	// find ship
	s, ok := e.findShip(o.Id)
	if !ok {
		p.Log("  name %s: no such ship\n", o.Id)
		return fmt.Errorf("no such ship %q", o.Id)
	}
	// fail if controlled by another player
	if s.ControlledBy != nil && s.ControlledBy != p {
		p.Log("  name %s: no such ship\n", o.Id)
		return fmt.Errorf("no such ship %q", o.Id)
	}
	// update the name
	s.Name = strings.Trim(o.Name, `"`)
	p.Log("  name %s: now named %q\n", o.Id, s.Name)
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
