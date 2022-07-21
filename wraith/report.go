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
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"log"
	"math"
	"sort"
	"time"
)

func (e *Engine) Report(w io.Writer, playerIds ...int) error {
	p := message.NewPrinter(language.English)

	asOfTurn := fmt.Sprintf("%04d/%d", e.Game.Turn.Year, e.Game.Turn.Quarter)
	rptDate := time.Now().Format("2006/01/02")
	showSUs := false

	var food *Unit
	for _, unit := range e.Units {
		if unit.Kind == "food" {
			food = unit
			break
		}
	}

	for _, playerId := range playerIds {
		player, ok := e.Players[playerId]
		if !ok {
			log.Printf("engine: reportWriter: playerId %d: not found\n", playerId)
			continue
		}
		nation := player.MemberOf
		if nation == nil {
			log.Printf("engine: reportWriter: playerId %d: nation: not found\n", playerId)
			continue
		}

		_, _ = p.Fprintf(w, "Status Report\n")
		_, _ = p.Fprintf(w, "Game: %-8s   Turn: %s   Nation: %3d   Player: %3d   Date: %s\n", e.Game.Code, asOfTurn, nation.No, player.Id, rptDate)

		_, _ = p.Fprintf(w, "\n------------------------------------------------------------------------------\n")
		_, _ = p.Fprintf(w, "Name: %-40s  Member Of: %s\n", player.Name, nation.Name)
		_, _ = p.Fprintf(w, "     Tech Level: %2d    Research Points: %11d\n", nation.TechLevel, nation.ResearchPointsPool)
		_, _ = p.Fprintf(w, "    Bureaucracy: %2d    Biology: %2d    Gravitics: %2d    LifeSupport: %2d\n", nation.Skills.Bureaucracy, nation.Skills.Biology, nation.Skills.Gravitics, nation.Skills.LifeSupport)
		_, _ = p.Fprintf(w, "  Manufacturing: %2d     Mining: %2d     Military: %2d        Shields: %2d\n", nation.Skills.Manufacturing, nation.Skills.Mining, nation.Skills.Military, nation.Skills.Shields)

		for _, cs := range player.Colonies {
			_, _ = p.Fprintf(w, "\n------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "%s Activity Report\n", cs.HullId)

			kind := "unknown"
			switch cs.Kind {
			case "enclosed":
				kind = "ENCLOSED"
			case "open", "surface":
				kind = "OPEN"
			case "orbital":
				kind = "ORBITAL"
			default:
				panic(fmt.Sprintf("assert(cs.Kind != %q)", cs.Kind))
			}
			name := cs.Name
			if name == "" {
				name = "NOT NAMED"
			}
			colonyKind := fmt.Sprintf("%s COLONY", kind)
			_, _ = p.Fprintf(w, "  Location: %s #%d    Tech: %2d  %14s: %-22s\n", cs.Planet.System.Coords.String(), cs.Planet.OrbitNo, cs.TechLevel, colonyKind, name)
			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Group____________  Population_Units  Pay_____  Rations_         CNGD/Turn         FOOD/Turn\n")
			_, _ = p.Fprintf(w, "  Professional       %16d  %7.3f%%  %7.3f%%  %16d  %16d\n", cs.Population.ProfessionalQty, cs.Pay.ProfessionalPct*100, cs.Rations.ProfessionalPct*100, cs.Pay.totalPay(cs.Population, "PRO"), cs.Rations.totalRations(cs.Population, "PRO"))
			_, _ = p.Fprintf(w, "  Soldier            %16d  %7.3f%%  %7.3f%%  %16d  %16d\n", cs.Population.SoldierQty, cs.Pay.SoldierPct*100, cs.Rations.SoldierPct*100, cs.Pay.totalPay(cs.Population, "SLD"), cs.Rations.totalRations(cs.Population, "SLD"))
			_, _ = p.Fprintf(w, "  Unskilled          %16d  %7.3f%%  %7.3f%%  %16d  %16d\n", cs.Population.UnskilledQty, cs.Pay.UnskilledPct*100, cs.Rations.UnskilledPct*100, cs.Pay.totalPay(cs.Population, "USK"), cs.Rations.totalRations(cs.Population, "USK"))
			_, _ = p.Fprintf(w, "  Unemployed         %16d  %7.3f%%  %7.3f%%  %16d  %16d\n", cs.Population.UnemployedQty, 0.0, cs.Rations.UnemployedPct*100, cs.Pay.totalPay(cs.Population, "UEM"), cs.Rations.totalRations(cs.Population, "UEM"))
			tPop := cs.Population.ProfessionalQty + cs.Population.SoldierQty + cs.Population.UnskilledQty + cs.Population.UnemployedQty
			tPay := cs.Pay.totalPay(cs.Population, "PRO") + cs.Pay.totalPay(cs.Population, "SLD") + cs.Pay.totalPay(cs.Population, "USK") + cs.Pay.totalPay(cs.Population, "UEM")
			tRations := cs.Rations.totalRations(cs.Population, "PRO") + cs.Rations.totalRations(cs.Population, "SLD") + cs.Rations.totalRations(cs.Population, "USK") + cs.Rations.totalRations(cs.Population, "UEM")
			_, _ = p.Fprintf(w, "  ----------------   %16d  --------  --------  %16d  %16d\n", tPop, tPay, tRations)

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Crew/Team________  Units___________\n")
			_, _ = p.Fprintf(w, "  Construction Crew  %16d\n", cs.Population.ConstructionCrewQty)
			_, _ = p.Fprintf(w, "  Spy Team           %16d\n", cs.Population.SpyTeamQty)
			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Changes__________  Population_Units\n")
			_, _ = p.Fprintf(w, "  Births             %16d\n", cs.Population.BirthsPriorTurn)
			_, _ = p.Fprintf(w, "  Non-Combat Deaths  %16d\n", cs.Population.NaturalDeathsPriorTurn)

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Hull and Systems ----------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Item-TL   Operational  Mass_______  Volume_____  Fuel_Cost__\n")
			operMass, operVolume, suOper, fuOper := 0, 0, 0, 0
			for _, item := range cs.Hull {
				mu := item.totalMass()
				var emu, su int
				if item.Unit.Code == "STUN" || item.Unit.Code == "LTSU" || item.Unit.Code == "SLSU" {
					// hull structures shouldn't require emu
				} else {
					emu = item.totalVolume()
					switch cs.Kind {
					case "open", "surface":
						su = emu * 1
					case "enclosed":
						su = emu * 5
					case "orbital":
						su = emu * 10
					default:
						panic(fmt.Sprintf("assert(cs.Kind != %q)", cs.Kind))
					}
				}

				fuelPerTurn := int(math.Ceil(float64(item.ActiveQty) * item.Unit.FuelPerUnitPerTurn))
				operMass, operVolume, suOper, fuOper = operMass+mu, operVolume+emu, suOper+su, fuOper+fuelPerTurn
				if showSUs {
					_, _ = p.Fprintf(w, "  %-7s  %12d  %11d  %11d  %11d  %11d\n", item.Unit.Code, item.ActiveQty, mu, emu, fuelPerTurn, su)
				} else {
					_, _ = p.Fprintf(w, "  %-7s  %12d  %11d  %11d  %11d\n", item.Unit.Code, item.ActiveQty, mu, emu, fuelPerTurn)
				}
			}
			_, _ = p.Fprintf(w, "  -------  ------------  -----------  -----------  -----------\n")
			if showSUs {
				_, _ = p.Fprintf(w, "   Totals  ------------  %11d  %11d  %11d  %11d\n", operMass, operVolume, fuOper, suOper)
			} else {
				_, _ = p.Fprintf(w, "   Totals  ------------  %11d  %11d  %11d\n", operMass, operVolume, fuOper)
			}

			availSUs := 0
			cargoMass, cargoVolume, suCargo := 0, 0, 0
			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Cargo and Supplies --------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Item-TL   Active_____  Stowed______  Total_Qty____  Mass________  Volume______\n")
			for _, item := range cs.Inventory {
				mu := item.totalMass()
				if item.Unit.Code == "STUN" || item.Unit.Code == "LTSU" || item.Unit.Code == "SLSU" {
					availSUs += item.ActiveQty
				}
				emu, su := item.totalVolume(), 0
				switch cs.Kind {
				case "open", "surface":
					su = emu * 1
				case "enclosed":
					su = emu * 5
				case "orbital":
					su = emu * 10
				default:
					panic(fmt.Sprintf("assert(cs.Kind != %q)", cs.Kind))
				}
				if showSUs {
					_, _ = p.Fprintf(w, "  %-7s  %12d  %12d  %13d  %12d  %12d  %13d\n", item.Unit.Code, item.ActiveQty, item.StowedQty, item.ActiveQty+item.StowedQty, mu, emu, su)
				} else {
					_, _ = p.Fprintf(w, "  %-7s  %12d  %12d  %13d  %12d  %12d\n", item.Unit.Code, item.ActiveQty, item.StowedQty, item.ActiveQty+item.StowedQty, mu, emu)
				}
				cargoMass, cargoVolume, suCargo = cargoMass+mu, cargoVolume+emu, suCargo+su
			}
			_, _ = p.Fprintf(w, "  -------  ------------  ------------  -------------  ------------  ------------\n")
			if showSUs {
				_, _ = p.Fprintf(w, "   Totals  ------------  ------------  -------------  %12d  %12d  %13d\n", cargoMass, cargoVolume, suCargo)
			} else {
				_, _ = p.Fprintf(w, "   Totals  ------------  ------------  -------------  %12d  %12d\n", cargoMass, cargoVolume)
			}

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Farming --------------------------------------------------------------------------------------------------------\n")
			for _, group := range cs.FarmGroups {
				_, _ = p.Fprintf(w, "  Group: %2d  Produces: %s\n", group.No, food.Code)
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.ActiveQty) * unit.Unit.FuelPerUnitPerTurn))
					proLabor, uskLabor := 1*unit.ActiveQty, 3*unit.ActiveQty
					_, _ = p.Fprintf(w, "     Input:  Farms__  Quantity_____  Professionals  Unskilled____  FUEL/Turn____\n")
					_, _ = p.Fprintf(w, "             %-7s  %13d  %13d  %13d  %13d\n", unit.Unit.Code, unit.ActiveQty, proLabor, uskLabor, fuelPerTurn)
				}
				_, _ = p.Fprintf(w, "    Output:  Unit___  Stage_1______  Stage_2______  Stage_3______\n")
				_, _ = p.Fprintf(w, "             %-7s  %13d  %13d  %13d\n", food.Code, group.StageQty[0], group.StageQty[1], group.StageQty[2])
			}

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Mining ---------------------------------------------------------------------------------------------------------\n")
			for _, group := range cs.MineGroups {
				_, _ = p.Fprintf(w, "  Group: %2d  Deposit: DP%-3d   Yield %6.3f%%    Remaining: %13d %-5s\n", group.No, group.Deposit.No, 100*group.Deposit.YieldPct, group.Deposit.RemainingQty, group.Deposit.Product.Code)
				fuelPerTurn := int(math.Ceil(float64(group.Unit.ActiveQty) * group.Unit.Unit.FuelPerUnitPerTurn))
				extractPerTurn := group.Unit.ActiveQty * 100 * group.Unit.Unit.TechLevel / 4
				proLabor, uskLabor := 1*group.Unit.ActiveQty, 3*group.Unit.ActiveQty
				_, _ = p.Fprintf(w, "     Input:  Mines__  Quantity_____  Professionals  Unskilled____  FUEL/Turn____  Extract/Turn_  Yield/Turn___\n")
				_, _ = p.Fprintf(w, "             %-7s  %13d  %13d  %13d  %13d  %13d  %13d\n", group.Unit.Unit.Code, group.Unit.ActiveQty, proLabor, uskLabor, fuelPerTurn, extractPerTurn, int(math.Ceil(float64(extractPerTurn)*group.Deposit.YieldPct)))
				_, _ = p.Fprintf(w, "    Output:  Unit___  Stage_1______  Stage_2______  Stage_3______\n")
				_, _ = p.Fprintf(w, "             %-7s  %13d  %13d  %13d\n", group.Deposit.Product.Code, group.StageQty[0], group.StageQty[1], group.StageQty[2])
			}

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Factory --------------------------------------------------------------------------------------------------------\n")
			for _, group := range cs.FactoryGroups {
				_, _ = p.Fprintf(w, "  Group: %2d  Produces: %s\n", group.No, group.Product.Code)
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.TotalQty) * unit.Unit.FuelPerUnitPerTurn))
					metsPerTurn := int(math.Ceil(float64(unit.TotalQty) * unit.Unit.MetsPerUnitPerTurn))
					nonMetsPerTurn := int(math.Ceil(float64(unit.TotalQty) * unit.Unit.NonMetsPerUnitPerTurn))
					proLabor, uskLabor := 1*unit.TotalQty, 3*unit.TotalQty
					_, _ = p.Fprintf(w, "     Input:  Facts__  Quantity_____  Professionals  Unskilled____  METS/Turn____  NMTS/Turn____  FUEL/Turn____\n")
					_, _ = p.Fprintf(w, "             %-7s  %13d  %13d  %13d  %13d  %13d  %13d\n", unit.Unit.Code, unit.TotalQty, proLabor, uskLabor, fuelPerTurn, metsPerTurn, nonMetsPerTurn)
				}
				_, _ = p.Fprintf(w, "    Output:  Unit___  Stage_1______  Stage_2______  Stage_3______\n")
				_, _ = p.Fprintf(w, "             %-7s  %13d  %13d  %13d\n", group.Product.Code, group.StageQty[0], group.StageQty[1], group.StageQty[2])
			}

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Espionage --------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  No activity.\n")
		}

		_, _ = p.Fprintf(w, "\n------------------------------------------------------------------------------\n")
		_, _ = p.Fprintf(w, "Surveys\n")
		starMap := make(map[int]*Star)
		for _, cs := range player.Colonies {
			starMap[cs.Planet.Star.Id] = cs.Planet.Star
		}
		for _, cs := range player.Ships {
			starMap[cs.Planet.Star.Id] = cs.Planet.Star
		}
		var stars Stars
		for _, star := range starMap {
			stars = append(stars, star)
		}
		sort.Sort(stars)
		for _, star := range stars {
			_, _ = p.Fprintf(w, "     Star: %s%s\n", star.System.Coords.String(), star.Sequence)
			for _, planet := range star.Planets {
				if planet == nil {
					continue
				}
				var habNo string
				if planet.HabitabilityNo > 0 {
					habNo = fmt.Sprintf("Habitability: %2d", planet.HabitabilityNo)
				}
				_, _ = p.Fprintf(w, "   Planet: %-13s    Kind: %-14s   %-17s\n", fmt.Sprintf("%s%s#%d", star.System.Coords.String(), star.Sequence, planet.OrbitNo), planet.Kind, habNo)
				for _, r := range planet.Deposits {
					_, _ = p.Fprintf(w, "           Deposit: %2d   %-6s   Yield: %7.3f%%   MUs remaining: %12d\n", r.No, r.Product.Code, r.YieldPct*100, r.RemainingQty)
				}
			}
		}

		_, _ = p.Fprintf(w, "\nMarket Report ---------------------------------------------------------------------\n")
		_, _ = p.Fprintf(w, "  News: *** 6 years ago, scientists discovered a signal broadcasting from the 10th\n")
		_, _ = p.Fprintf(w, "        orbit. 4 years ago, they decoded the message and recovered plans for an\n")
		_, _ = p.Fprintf(w, "        in-system engine (the \"space-drive\") and a faster-than-light engine\n")
		_, _ = p.Fprintf(w, "        (the \"hyper-drive\"). Work on both has recently completed.\n")
		_, _ = p.Fprintf(w, "  News: *** Moments after the hyper-drive was successfully tested in the orbital\n")
		_, _ = p.Fprintf(w, "        colony, the broadcast from the 10th orbit stopped.\n")

		_, _ = p.Fprintf(w, "\nCombat Report ---------------------------------------------------------------------\n")
		_, _ = p.Fprintf(w, "  No activity.\n")
	}

	return nil
}
