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
	"github.com/mdhender/wraith/models"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"math"
	"time"
)

func (e *Engine) Report(spId int) error {
	panic("!")
	//// run the report for just the one nation
	//for _, n := range e.game.Nations {
	//	if n.No != spId {
	//		continue
	//	}
	//	log.Printf("reporting for %d %q (%d)\n", n.No, n.Name, n.Id)
	//
	//	//reportFile := filepath.Clean(filepath.Join(n.Store, "report.txt"))
	//	//log.Printf("reporting for %d %q (%q)\n", n.Id, n.Name, reportFile)
	//
	//	//w, err := os.OpenFile(reportFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	//	//if err != nil {
	//	//	return err
	//	//}
	//
	//	return e.ReportWriter(os.Stdout, spId)
	//}
	//
	//return nil
}

func (e *Engine) ReportWriter(game *models.Game, w io.Writer) error {
	p := message.NewPrinter(language.English)

	asOfTurn := game.CurrentTurn
	rptDate := time.Now().Format("2006/01/02")

	for _, nation := range game.Nations {
		if nation.Details[0].Name != "Yinshei" {
			continue
		}

		_, _ = p.Fprintf(w, "Status Report\n")
		_, _ = p.Fprintf(w, "Game: %-8s   Turn: %s   Nation: %3d    Date: %s\n", game.ShortName, asOfTurn.String(), nation.No, rptDate)

		_, _ = p.Fprintf(w, "\nNation %3d -------------------------------------------------------------------\n", nation.No)
		for _, detail := range nation.Details {
			if detail.ControlledBy == nil {
				_, _ = p.Fprintf(w, "  Name: %-40s  Controlled By: %s\n", detail.Name, "*nobody*")
			} else {
				_, _ = p.Fprintf(w, "  Name: %-40s  Controlled By: %s\n", detail.Name, detail.ControlledBy.Details[0].Handle)
			}
		}
		for _, research := range nation.Research {
			_, _ = p.Fprintf(w, "     Tech Level: %2d    Research Points: %11d\n", research.TechLevel, research.ResearchPointsPool)
		}
		for _, skills := range nation.Skills {
			_, _ = p.Fprintf(w, "    Bureaucracy: %2d    Biology: %2d    Gravitics: %2d    LifeSupport: %2d\n", skills.Bureaucracy, skills.Biology, skills.Gravitics, skills.LifeSupport)
			_, _ = p.Fprintf(w, "  Manufacturing: %2d     Mining: %2d     Military: %2d        Shields: %2d\n", skills.Manufacturing, skills.Mining, skills.Military, skills.Shields)
		}

		for _, cs := range nation.Colonies {
			msnId := fmt.Sprintf("C%d", cs.MSN)
			_, _ = p.Fprintf(w, "\n%s Activity Report -------------------------------------------------------\n", msnId)

			kind := "unknown"
			switch cs.Kind {
			case "enclosed":
				kind = "ENCLOSED"
			case "open":
				kind = "OPEN"
			case "orbital":
				kind = "ORBITAL"
			default:
				panic(fmt.Sprintf("assert(cs.Kind != %q)", cs.Kind))
			}
			name := cs.Details[0].Name
			if name == "" {
				name = "NOT NAMED"
			}
			techLevel := cs.Details[0].TechLevel
			location := cs.Locations[0]
			planet := location.Location
			orbitNo := planet.OrbitNo
			star := planet.Star
			system := star.System
			colonyKind := fmt.Sprintf("%s COLONY", kind)
			_, _ = p.Fprintf(w, "  Location: %s #%d    Tech: %2d  %14s: %-22s\n", system.Coords.String(), orbitNo, techLevel, colonyKind, name)

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Group____________  Population_Units  Pay_____  Rations_         CNGD/Turn         FOOD/Turn\n")
			_, _ = p.Fprintf(w, "  Professional       %16d  %7.3f%%  %7.3f%%  %16d  %16d\n", cs.Population[0].QtyProfessional, cs.Pay[0].ProfessionalPct*100, cs.Rations[0].ProfessionalPct*100, totalPay("PRO", cs.Pay[0].ProfessionalPct, cs.Population[0].QtyProfessional), totalRations("PRO", cs.Rations[0].ProfessionalPct, cs.Population[0].QtyProfessional))
			_, _ = p.Fprintf(w, "  Soldier            %16d  %7.3f%%  %7.3f%%  %16d  %16d\n", cs.Population[0].QtySoldier, cs.Pay[0].SoldierPct*100, cs.Rations[0].SoldierPct*100, totalPay("SLD", cs.Pay[0].SoldierPct, cs.Population[0].QtySoldier), totalRations("SLD", cs.Rations[0].SoldierPct, cs.Population[0].QtySoldier))
			_, _ = p.Fprintf(w, "  Unskilled          %16d  %7.3f%%  %7.3f%%  %16d  %16d\n", cs.Population[0].QtyUnskilled, cs.Pay[0].UnskilledPct*100, cs.Rations[0].UnskilledPct*100, totalPay("USK", cs.Pay[0].UnskilledPct, cs.Population[0].QtyUnskilled), totalRations("USK", cs.Rations[0].UnskilledPct, cs.Population[0].QtyUnskilled))
			_, _ = p.Fprintf(w, "  Unemployed         %16d  %7.3f%%  %7.3f%%  %16d  %16d\n", cs.Population[0].QtyUnemployed, cs.Pay[0].UnemployedPct*100, cs.Rations[0].UnemployedPct*100, totalPay("UEM", cs.Pay[0].UnemployedPct, cs.Population[0].QtyUnemployed), totalRations("UEM", cs.Rations[0].UnemployedPct, cs.Population[0].QtyUnemployed))
			tPop := cs.Population[0].QtyProfessional + cs.Population[0].QtySoldier + cs.Population[0].QtyUnskilled + cs.Population[0].QtyUnemployed
			tPay := totalPay("PRO", cs.Pay[0].ProfessionalPct, cs.Population[0].QtyProfessional) + totalPay("SLD", cs.Pay[0].SoldierPct, cs.Population[0].QtySoldier) + totalPay("USK", cs.Pay[0].UnskilledPct, cs.Population[0].QtyUnskilled) + totalPay("UEM", cs.Pay[0].UnemployedPct, cs.Population[0].QtyUnemployed)
			tRations := totalRations("PRO", cs.Rations[0].ProfessionalPct, cs.Population[0].QtyProfessional) + totalRations("SLD", cs.Rations[0].SoldierPct, cs.Population[0].QtySoldier) + totalRations("USK", cs.Rations[0].UnskilledPct, cs.Population[0].QtyUnskilled) + totalRations("UEM", cs.Rations[0].UnemployedPct, cs.Population[0].QtyUnemployed)
			_, _ = p.Fprintf(w, "  ----------------   %16d  --------  --------  %16d  %16d\n", tPop, tPay, tRations)

			//		if cs.Population[0].Births == 0 && cs.Population[0].Deaths == 0 {
			//			cs.Population[0].Births = cs.Population[0].TotalPopulation() / 1600
			//			cs.Population[0].Deaths = cs.Population[0].Births
			//		}

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Crew/Team________  Units___________\n")
			_, _ = p.Fprintf(w, "  Construction Crew  %16d\n", cs.Population[0].QtyConstructionCrew)
			_, _ = p.Fprintf(w, "  Spy Team           %16d\n", cs.Population[0].QtySoldier)
			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Changes__________  Population_Units\n")
			_, _ = p.Fprintf(w, "  Births             %16d\n", cs.Population[0].Births)
			_, _ = p.Fprintf(w, "  Non-Combat Deaths  %16d\n", cs.Population[0].Deaths)

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Components ----------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Item-TL   Operational          MUs         EMUs  SUs_Required\n")
			muOper, emuOper, suOper := 0, 0, 0
			for _, item := range cs.Hull {
				mu := totalMass(item.Unit, item.QtyOperational, 0)
				var emu, su int
				if item.Unit.Code == "STUN" || item.Unit.Code == "LTSU" || item.Unit.Code == "SLSU" {
					// hull structures shouldn't require emu
				} else {
					emu = totalVolume(item.Unit, item.QtyOperational, 0)
					switch cs.Kind {
					case "open":
						su = emu * 1
					case "enclosed":
						su = emu * 5
					case "orbital":
						su = emu * 10
					}
				}
				muOper, emuOper, suOper = muOper+mu, emuOper+emu, suOper+su
				_, _ = p.Fprintf(w, "  %7s  %12d  %11d  %11d   %11d\n", item.Unit.Code, item.QtyOperational, mu, emu, su)
			}
			_, _ = p.Fprintf(w, "   Totals  ------------  %11d  %11d   %11d\n", muOper, emuOper, suOper)

			availSUs := 0
			muCargo, emuCargo, suCargo := 0, 0, 0
			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Cargo ---------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Item-TL   Operational        Stowed      Total Qty           MUs          EMUs   SUs_Required\n")
			for _, item := range cs.Inventory {
				mu := totalMass(item.Unit, item.QtyOperational, item.QtyStowed)
				if item.Unit.Code == "STUN" || item.Unit.Code == "LTSU" || item.Unit.Code == "SLSU" {
					availSUs += item.QtyOperational
				}
				emu, su := totalVolume(item.Unit, item.QtyOperational, item.QtyStowed), 0
				switch cs.Kind {
				case "open":
					su = emu * 1
				case "enclosed":
					su = emu * 5
				case "orbital":
					su = emu * 10
				}
				_, _ = p.Fprintf(w, "  %7s  %12d  %12d  %13d  %12d  %12d  %13d\n",
					item.Unit.Code, item.QtyOperational, item.QtyStowed, item.QtyOperational+item.QtyStowed, mu, emu, su)
				muCargo, emuCargo, suCargo = muCargo+mu, emuCargo+emu, suCargo+su
			}
			_, _ = p.Fprintf(w, "   Totals  ------------  ------------  -------------  %12d  %12d  %13d\n", muCargo, emuCargo, suCargo)

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Farming --------------------------------------------------------------------------------------------------------\n")
			for _, group := range cs.Farms {
				_, _ = p.Fprintf(w, "  Group: %2d  Produces: %s\n", group.No, group.Unit.Code)
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.QtyOperational) * 0.5))
					proLabor, uskLabor := 1*unit.QtyOperational, 3*unit.QtyOperational
					_, _ = p.Fprintf(w, "     Input:  Farms_     Quantity  Professionals    Unskilled    FUEL/Turn\n")
					_, _ = p.Fprintf(w, "             %6s  %11d  %11d    %11d  %11d\n",
						unit.Unit.Code, unit.QtyOperational, proLabor, uskLabor, fuelPerTurn)
				}
				for _, stage := range group.Stages {
					_, _ = p.Fprintf(w, "    Output:  Unit__      Stage_1      Stage_2        Stage_3\n")
					_, _ = p.Fprintf(w, "             %6s  %11d  %11d    %11d\n",
						group.Unit.Code, stage.QtyStage1, stage.QtyStage2, stage.QtyStage3)
				}
			}

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Mining ---------------------------------------------------------------------------------------------------------\n")
			for _, group := range cs.Mines {
				_, _ = p.Fprintf(w, "  Group: %2d  Produces: %s\n", group.No, group.Deposit.Unit.Code)
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.QtyOperational) * 0.5))
					proLabor, uskLabor := 1*unit.QtyOperational, 3*unit.QtyOperational
					_, _ = p.Fprintf(w, "     Input:  Mines_     Quantity  Professionals    Unskilled    FUEL/Turn\n")
					_, _ = p.Fprintf(w, "             %6s  %11d  %11d    %11d  %11d\n",
						unit.Unit.Code, unit.QtyOperational, proLabor, uskLabor, fuelPerTurn)
				}
				for _, stage := range group.Stages {
					_, _ = p.Fprintf(w, "    Output:  Unit__      Stage_1      Stage_2        Stage_3\n")
					_, _ = p.Fprintf(w, "             %6s  %11d  %11d    %11d\n",
						group.Deposit.Unit.Code, stage.QtyStage1, stage.QtyStage2, stage.QtyStage3)
				}
			}

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Production -----------------------------------------------------------------------------------------------------\n")
			for _, group := range cs.Factories {
				_, _ = p.Fprintf(w, "  Group: %2d  Produces: %s\n", group.No, group.Unit.Code)
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.QtyOperational) * 0.5))
					proLabor, uskLabor := 1*unit.QtyOperational, 3*unit.QtyOperational
					_, _ = p.Fprintf(w, "     Input:  Facts_     Quantity  Professionals    Unskilled    METS/Turn    NMTS/Turn    FUEL/Turn\n")
					_, _ = p.Fprintf(w, "             %6s  %11d  %11d    %11d            ?            ?  %11d\n",
						unit.Unit.Code, unit.QtyOperational, proLabor, uskLabor, fuelPerTurn)
				}
				for _, stage := range group.Stages {
					_, _ = p.Fprintf(w, "    Output:  Unit__      Stage_1      Stage_2        Stage_3\n")
					_, _ = p.Fprintf(w, "             %6s  %11d  %11d    %11d\n",
						group.Unit.Code, stage.QtyStage1, stage.QtyStage2, stage.QtyStage3)
				}
			}

			_, _ = p.Fprintf(w, "\n")
			_, _ = p.Fprintf(w, "  Espionage --------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  No activity.\n")
		}

		//	orbits := []ReportOrbit{
		//		{Id: 1, Kind: "terrestrial", HabitabilityNumber: 0},
		//		{Id: 2, Kind: "terrestrial", HabitabilityNumber: 0},
		//		{Id: 3, Kind: "terrestrial", HabitabilityNumber: 25},
		//		{Id: 4, Kind: "terrestrial", HabitabilityNumber: 0},
		//		{Id: 5, Kind: "asteroid belt", HabitabilityNumber: 0},
		//		{Id: 6, Kind: "gas giant", HabitabilityNumber: 0},
		//		{Id: 7, Kind: "gas giant", HabitabilityNumber: 0},
		//		{Id: 8, Kind: "terrestrial", HabitabilityNumber: 0},
		//		{Id: 9, Kind: "terrestrial", HabitabilityNumber: 0},
		//		{Id: 10, Kind: "asteroid belt", HabitabilityNumber: 0},
		//	}
		_, _ = p.Fprintf(w, "\nSurvey Report ----------------------------------------------------------------\n")
		//	for _, orbit := range orbits {
		//		name, controlledBy := "NOT NAMED", "N/A"
		//		if orbit.Id == 3 {
		//			name, controlledBy = "My Homeworld", "SP018"
		//		}
		//		_, _ = p.Fprintf(w, "  ReportSystem %s %-3s   %-24s    Controlled By: %s\n", "0/0/0", fmt.Sprintf("#%d", orbit.Id), name, controlledBy)
		//		_, _ = p.Fprintf(w, "    Kind: %-13s    Habitability: %2d\n", orbit.Kind, orbit.HabitabilityNumber)
		//	}
		//
		//_, _ = p.Fprintf(w, "\nMarket Report ---------------------------------------------------------------------\n")
		//_, _ = p.Fprintf(w, "  News: *** 6 years ago, scientists discovered a signal broadcasting from the 10th\n")
		//_, _ = p.Fprintf(w, "        orbit. 4 years ago, they decoded the message and recovered plans for an\n")
		//_, _ = p.Fprintf(w, "        in-system engine (the \"space-drive\") and a faster-than-light engine\n")
		//_, _ = p.Fprintf(w, "        (the \"hyper-drive\"). Work on both has recently completed.\n")
		//_, _ = p.Fprintf(w, "  News: *** Moments after the hyper-drive was successfully tested in the orbital\n")
		//_, _ = p.Fprintf(w, "        colony, the broadcast from the 10th orbit stopped.\n")
		//
		//_, _ = p.Fprintf(w, "\nCombat Report ---------------------------------------------------------------------\n")
		//_, _ = p.Fprintf(w, "  No activity.\n")

		break
	}
	return nil
}

func totalMass(unit *models.Unit, oper, stowed int) int {
	return int(math.Ceil(float64(oper+stowed) * unit.MassPerUnit))
}

func totalVolume(unit *models.Unit, oper, stowed int) int {
	return int(math.Ceil(float64(oper)*unit.VolumePerUnit)) + int(math.Ceil(float64(stowed)*unit.StowedVolumePerUnit))
}

// totalPay assumes that the base rates are per unit of population
//  PROFESSIONAL      0.375 CONSUMER GOODS
//  SOLDIER           0.250 CONSUMER GOODS
//  UNSKILLED WORKER  0.125 CONSUMER GOODS
//  UNEMPLOYABLE      0.000 CONSUMER GOODS
func totalPay(code string, pct float64, qty int) int {
	switch code {
	case "PRO":
		return int(math.Ceil((0.375 * pct) * float64(qty)))
	case "SLD":
		return int(math.Ceil((0.250 * pct) * float64(qty)))
	case "USK":
		return int(math.Ceil((0.125 * pct) * float64(qty)))
	case "UEM":
		return 0
	default:
		panic(fmt.Sprintf("assert(ReportPopUnit.Code != %q)", code))
	}
}

// totalRations assumes that base rates are per unit of population
//  PROFESSIONAL      0.250 FOOD
//  SOLDIER           0.250 FOOD
//  UNSKILLED WORKER  0.250 FOOD
//  UNEMPLOYABLE      0.250 FOOD
func totalRations(code string, pct float64, qty int) int {
	switch code {
	case "PRO":
		return int(math.Ceil((0.25 * pct) * (float64(qty))))
	case "SLD":
		return int(math.Ceil((0.25 * pct) * (float64(qty))))
	case "USK":
		return int(math.Ceil((0.25 * pct) * (float64(qty))))
	case "UEM":
		return int(math.Ceil((0.25 * pct) * (float64(qty))))
	default:
		panic(fmt.Sprintf("assert(ReportPopUnit.Code != %q)", code))
	}
}

//type ReportStore struct {
//	Game struct {
//		Id   string `json:"id"`
//		Turn int    `json:"turn"`
//	} `json:"game"`
//	Players map[string]*ReportPlayer `json:"players,omitempty"`
//}
//
//type ReportPlayer struct {
//	Skills   *Skills         `json:"skills,omitempty"`
//	Colonies []*ReportColony `json:"colonies,omitempty"`
//	Ships    []*ReportShip   `json:"ships,omitempty"`
//}
//
//type ReportPopUnit struct {
//	Code   string  `json:"code"`
//	Qty    int     `json:"qty,omitempty"`
//	Pay    float64 `json:"pay,omitempty"`
//	Ration float64 `json:"ration,omitempty"`
//}
//
//type ReportUnit struct {
//	Name      string `json:"name"`
//	TechLevel int    `json:"tech-level,omitempty"`
//	Qty       int    `json:"qty,omitempty"`
//	Stowed    bool   `json:"stowed,omitempty"`
//	code      string
//}
//
//func (u ReportUnit) Code() string {
//	if u.code == "" {
//		switch u.Name {
//		case "anti-missile":
//			u.code = fmt.Sprintf("ANM-%d", u.TechLevel)
//		case "assault-craft":
//			u.code = fmt.Sprintf("ASC-%d", u.TechLevel)
//		case "assault-weapon":
//			u.code = fmt.Sprintf("ASW-%d", u.TechLevel)
//		case "automation":
//			u.code = fmt.Sprintf("AUT-%d", u.TechLevel)
//		case "consumer-goods":
//			u.code = "CNGD"
//		case "energy-shield":
//			u.code = fmt.Sprintf("ESH-%d", u.TechLevel)
//		case "energy-weapon":
//			u.code = fmt.Sprintf("EWP-%d", u.TechLevel)
//		case "factory":
//			u.code = fmt.Sprintf("FCT-%d", u.TechLevel)
//		case "farm":
//			u.code = fmt.Sprintf("FRM-%d", u.TechLevel)
//		case "food":
//			u.code = "FOOD"
//		case "fuel":
//			u.code = "FUEL"
//		case "gold":
//			u.code = "GOLD"
//		case "hyper-drive":
//			u.code = fmt.Sprintf("HDR-%d", u.TechLevel)
//		case "life-support":
//			u.code = fmt.Sprintf("LSP-%d", u.TechLevel)
//		case "light-structural":
//			u.code = "LTSU"
//		case "metallics":
//			u.code = "MTLS"
//		case "military-robots":
//			u.code = fmt.Sprintf("MLR-%d", u.TechLevel)
//		case "military-supplies":
//			u.code = "MLSP"
//		case "mine":
//			u.code = fmt.Sprintf("MIN-%d", u.TechLevel)
//		case "missile":
//			u.code = fmt.Sprintf("MSS-%d", u.TechLevel)
//		case "missile-launcher":
//			u.code = fmt.Sprintf("MSL-%d", u.TechLevel)
//		case "non-metallics":
//			u.code = "NMTS"
//		case "sensor":
//			u.code = fmt.Sprintf("SNR-%d", u.TechLevel)
//		case "space-drive":
//			u.code = fmt.Sprintf("SDR-%d", u.TechLevel)
//		case "structural":
//			u.code = "STUN"
//		case "super-light-structural":
//			u.code = "SLSU"
//		case "transport":
//			u.code = fmt.Sprintf("TPT-%d", u.TechLevel)
//		default:
//			panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
//		}
//	}
//	return u.code
//}
//
//func (u ReportUnit) FuelPerTurn() int {
//	switch u.Name {
//	case "anti-missile":
//		return 0
//	case "assault-craft":
//		return 0
//	case "assault-weapon":
//		return 2 * u.TechLevel
//	case "automation":
//		return 0
//	case "consumer-goods":
//		return 0
//	case "energy-shield":
//		return 0
//	case "energy-weapon":
//		return 0
//	case "factory":
//		return int(math.Ceil(float64(u.TechLevel) / 2))
//	case "farm":
//		if u.TechLevel < 6 {
//			return int(math.Ceil(float64(u.TechLevel) / 2))
//		}
//		return u.TechLevel
//	case "hyper-drive":
//		return 0
//	case "life-support":
//		return 1
//	case "light-structural":
//		return 0
//	case "military-robots":
//		return 0
//	case "military-supplies":
//		return 0
//	case "mine":
//		return int(math.Ceil(float64(u.TechLevel) / 2))
//	case "missile":
//		return 0
//	case "missile-launcher":
//		return 0
//	case "sensor":
//		return int(math.Ceil(float64(u.TechLevel) / 20))
//	case "space-drive":
//		return 0
//	case "structural":
//		return 0
//	case "super-light-structural":
//		return 0
//	case "transport":
//		return int(math.Ceil(float64(u.TechLevel*u.TechLevel) / 10))
//	default:
//		panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
//	}
//}
//
//func (u ReportUnit) IngestPerTurn() int {
//	if u.Name != "factory" {
//		return 0
//	}
//	// can ingest 20 MU of resources per tech-level per YEAR
//	return u.Qty * u.TechLevel * 20 / 4
//}
//
//func (u ReportUnit) LaborPerTurn() (pro, usk int) {
//	if u.Name != "factory" {
//		return 1 * u.Qty, 3 * u.Qty
//	} else if u.Qty >= 50_000 {
//		return 1 * u.Qty, 3 * u.Qty
//	} else if u.Qty >= 5_000 {
//		return 2 * u.Qty, 6 * u.Qty
//	} else if u.Qty >= 500 {
//		return 3 * u.Qty, 9 * u.Qty
//	} else if u.Qty >= 50 {
//		return 4 * u.Qty, 12 * u.Qty
//	} else if u.Qty >= 5 {
//		return 5 * u.Qty, 15 * u.Qty
//	}
//	return 6 * u.Qty, 18 * u.Qty
//}
//
//func (u ReportUnit) EnclosedMassUnits() int {
//	mus := u.MassUnits()
//	if !u.Hudnut() || !u.Stowed {
//		return mus
//	}
//	return int(math.Ceil(float64(mus) / 2))
//}
//
//func (u ReportUnit) MassUnits() int {
//	switch u.Name {
//	case "anti-missile":
//		return 4 * u.TechLevel * u.Qty
//	case "assault-craft":
//		return 5 * u.TechLevel * u.Qty
//	case "assault-weapon":
//		return 2 * u.TechLevel * u.Qty
//	case "automation":
//		return 4 * u.TechLevel * u.Qty
//	case "consumer-goods":
//		return int(math.Ceil(0.6 * float64(u.Qty)))
//	case "energy-shield":
//		return 50 * u.TechLevel * u.Qty
//	case "energy-weapon":
//		return 10 * u.TechLevel * u.Qty
//	case "factory":
//		return (12 + 2*u.TechLevel) * u.Qty
//	case "farm":
//		return (6 + 2*u.TechLevel) * u.Qty
//	case "food":
//		return 6 * u.Qty
//	case "fuel":
//		return 1 * u.Qty
//	case "gold":
//		return 1 * u.Qty
//	case "hyper-drive":
//		return 45 * u.TechLevel * u.Qty
//	case "life-support":
//		return 8 * u.TechLevel * u.Qty
//	case "light-structural":
//		return int(math.Ceil(0.05 * float64(u.Qty)))
//	case "metallics":
//		return 1 * u.Qty
//	case "military-robots":
//		return (20 + 2*u.TechLevel) * u.Qty
//	case "military-supplies":
//		return int(math.Ceil(0.04 * float64(u.Qty)))
//	case "mine":
//		return (10 + 2*u.TechLevel) * u.Qty
//	case "missile":
//		return 4 * u.TechLevel * u.Qty
//	case "missile-launcher":
//		return 25 * u.TechLevel * u.Qty
//	case "non-metallics":
//		return 1 * u.Qty
//	case "sensor":
//		return 40 * u.TechLevel * u.Qty
//	case "space-drive":
//		return 25 * u.TechLevel * u.Qty
//	case "structural":
//		return int(math.Ceil(0.5 * float64(u.Qty)))
//	case "super-light-structural":
//		return int(math.Ceil(0.005 * float64(u.Qty)))
//	case "transport":
//		return int(math.Ceil(0.1 * float64(u.TechLevel*u.TechLevel*u.Qty)))
//	default:
//		panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
//	}
//}
//
//func (u ReportUnit) Assembled() string {
//	if !u.Stowed {
//		return "??"
//	}
//	return "???"
//}
//
//// Hudnut is borrowed from the Sniglet for the leftover
//// bolts and such from a some-assembly-required project.
//// Returns true if the unit can be disassembled for storage.
//func (u ReportUnit) Hudnut() bool {
//	switch u.Name {
//	case "anti-missile":
//		return false
//	case "assault-craft":
//		return false
//	case "assault-weapon":
//		return false
//	case "automation":
//		return true
//	case "consumer-goods":
//		return false
//	case "energy-shield":
//		return true
//	case "energy-weapon":
//		return true
//	case "factory":
//		return true
//	case "farm":
//		return true
//	case "food":
//		return false
//	case "fuel":
//		return false
//	case "gold":
//		return false
//	case "hyper-drive":
//		return true
//	case "life-support":
//		return true
//	case "light-structural":
//		return true
//	case "metallics":
//		return false
//	case "military-robots":
//		return false
//	case "military-supplies":
//		return false
//	case "mine":
//		return true
//	case "missile":
//		return false
//	case "missile-launcher":
//		return true
//	case "non-metallics":
//		return false
//	case "sensor":
//		return true
//	case "space-drive":
//		return true
//	case "structural":
//		return true
//	case "super-light-structural":
//		return true
//	case "transport":
//		return false
//	default:
//		panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
//	}
//}
//
//func (u ReportUnit) Operational() string {
//	if !u.Stowed {
//		return "yes"
//	}
//	switch u.code {
//	case "FCT":
//		return "no"
//	case "MIN":
//		return "no"
//	}
//	return "yes"
//}
//
//func (u ReportUnit) RawMaterials() (mets, nmts float64) {
//	tl := float64(u.TechLevel)
//	switch u.Name {
//	case "anti-missile":
//		return 2 * tl, 2 * tl
//	case "assault-craft":
//		return 3 * tl, 2 * tl
//	case "assault-weapon":
//		return 1 * tl, 1 * tl
//	case "automation":
//		return 2 * tl, 2 * tl
//	case "consumer-goods":
//		return 0.2, 0.4
//	case "energy-shield":
//		return 25 * tl, 25 * tl
//	case "energy-weapon":
//		return 5 * tl, 5 * tl
//	case "factory":
//		return 8 * tl, 4 * tl
//	case "farm":
//		return 4 + tl, 2 + tl
//	case "hyper-drive":
//		return 25 * tl, 20 * tl
//	case "life-support":
//		return 3 * tl, 5 * tl
//	case "light-structural":
//		return 0.01, 0.04
//	case "military-robots":
//		return 10 * tl, 10 * tl
//	case "military-supplies":
//		return 0.02, 0.02
//	case "mine":
//		return 5 + tl, 5 + tl
//	case "missile":
//		return 2 * tl, 2 * tl
//	case "missile-launcher":
//		return 15 * tl, 10 * tl
//	case "sensor":
//		return 10 * tl, 20 * tl
//	case "space-drive":
//		return 15 * tl, 10 * tl
//	case "structural":
//		return 0.1, 0.4
//	case "super-light-structural":
//		return 0.001, 0.004
//	case "transport":
//		return 3 * tl, tl
//	}
//	panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
//}
//
//func UnitAttributes(name string, techLevel int) (mets, nmts, totalMassUnits, fuelPerTurn, fuelPerCombatRound float64) {
//	tl := float64(techLevel)
//	switch name {
//	case "anti-missile":
//		return 2 * tl, 2 * tl, 4 * tl, 0, 0
//	case "assault-craft":
//		return 3 * tl, 2 * tl, 5 * tl, 0, 0.1
//	case "assault-weapon":
//		return 1 * tl, 1 * tl, 2 * tl, 2 * tl * tl, 0
//	case "automation":
//		return 2 * tl, 2 * tl, 4 * tl, 0, 0
//	case "consumer-goods":
//		return 0.2, 0.4, 0.6, 0, 0
//	case "energy-shield":
//		return 25 * tl, 25 * tl, 50 * tl, 0, 10 * tl
//	case "energy-weapon":
//		return 5 * tl, 5 * tl, 10 * tl, 0, 4 * tl
//	case "factory":
//		return 8 * tl, 4 * tl, 12 + 2*tl, 0.5 * tl, 4 * tl
//	case "farm":
//		if techLevel == 1 {
//			return 4 + tl, 2 + tl, 6 + 2*tl, 0.5 * tl, 0
//		} else if techLevel < 6 {
//			return 4 + tl, 4 + tl, 6 + 2*tl, 0.5 * tl, 0
//		}
//		return 4 + tl, 2 + tl, 6 + 2*tl, tl, 0
//	case "hyper-drive":
//		return 25 * tl, 20 * tl, 45 * tl, 0, 0
//	case "life-support":
//		return 3 * tl, 5 * tl, 8 * tl, 1, 0
//	case "light-structural":
//		return 0.01, 0.04, 0.05, 0, 0
//	case "military-robots":
//		return 10 * tl, 10 * tl, 20 + 2*tl, 0, 0
//	case "military-supplies":
//		return 0.02, 0.02, 0.04, 0, 0
//	case "mine":
//		return 5 + tl, 5 + tl, 10 + (2 * tl), 0.5 * tl, 0
//	case "missile":
//		return 2 * tl, 2 * tl, 4 * tl, 0, 0
//	case "missile-launcher":
//		return 15 * tl, 10 * tl, 25 * tl, 0, 0
//	case "sensor":
//		return 10 * tl, 20 * tl, 40 * tl, tl / 20, 0
//	case "space-drive":
//		return 15 * tl, 10 * tl, 25 * tl, 0, tl * tl
//	case "structural":
//		return 0.1, 0.4, 0.5, 0, 0
//	case "super-light-structural":
//		return 0.001, 0.004, 0.005, 0, 0
//	case "transport":
//		return 3 * tl, tl, 4 * tl, 0.1 * tl * tl, 0.01 * tl * tl
//	}
//	panic(fmt.Sprintf("assert(unit.name != %q)", name))
//}
//
//type ReportGroup struct {
//	Id        int                `json:"id"`
//	Name      string             `json:"name"`
//	TechLevel int                `json:"tech-level,omitempty"`
//	Units     []*ReportGroupUnit `json:"units,omitempty"`
//	code      string
//}
//
//func (g ReportGroup) Code() string {
//	if g.code == "" {
//		g.code = ReportUnit{Name: g.Name, TechLevel: g.TechLevel}.Code()
//	}
//	return g.code
//}
//
//type ReportGroupUnit struct {
//	TechLevel int   `json:"tech-level,omitempty"`
//	Qty       int   `json:"qty,omitempty"`
//	Stages    []int `json:"stages,omitempty"`
//}
//
//type ReportColony struct {
//	Id            int               `json:"id"`
//	Name          string            `json:"name,omitempty"`
//	System        string            `json:"system"`
//	Orbit         int               `json:"orbit"`
//	Kind          string            `json:"kind"`
//	TechLevel     int               `json:"tech-level"`
//	Population    *ReportPopulation `json:"population,omitempty"`
//	Operational   []*ReportUnit     `json:"operational,omitempty"`
//	Storage       []*ReportUnit     `json:"storage,omitempty"`
//	FactoryGroups []*ReportGroup    `json:"factory-groups"`
//	FarmGroups    []*ReportGroup    `json:"farm-groups,omitempty"`
//	MiningGroups  []*ReportGroup    `json:"mining-groups,omitempty"`
//}
//
//type ReportShip struct {
//	Id          int               `json:"id"`
//	Name        string            `json:"name,omitempty"`
//	TechLevel   int               `json:"tech-level"`
//	Population  *ReportPopulation `json:"population,omitempty"`
//	Operational []*ReportUnit     `json:"units,omitempty"`
//	Storage     []*ReportUnit     `json:"storage,omitempty"`
//	FarmGroups  []*ReportGroup    `json:"farm-groups,omitempty"`
//}
//
//type ReportOrbit struct {
//	Id                 int
//	Name               string
//	Kind               string
//	HabitabilityNumber int
//}
//
//type ReportPopulation struct {
//	PRO    ReportPopUnit `json:"pro,omitempty"`
//	SLD    ReportPopUnit `json:"sld,omitempty"`
//	USK    ReportPopUnit `json:"usk,omitempty"`
//	UEM    ReportPopUnit `json:"uem,omitempty"`
//	CNW    int           `json:"cnw,omitempty"`
//	SPY    int           `json:"spy,omitempty"`
//	Births int           `json:"births,omitempty"`
//	Deaths int           `json:"deaths,omitempty"`
//}
//
//func (p *ReportPopulation) TotalPay() int {
//	if p == nil {
//		return 0
//	}
//	return p.PRO.TotalPay() + p.SLD.TotalPay() + p.USK.TotalPay() + p.UEM.TotalPay()
//
//}
//
//func (p *ReportPopulation) TotalPopulation() int {
//	if p == nil {
//		return 0
//	}
//	return p.PRO.Qty + p.SLD.Qty + p.USK.Qty + p.UEM.Qty
//}
//
//func (p *ReportPopulation) TotalRation() int {
//	if p == nil {
//		return 0
//	}
//	return p.PRO.TotalRation() + p.SLD.TotalRation() + p.USK.TotalRation() + p.UEM.TotalRation()
//}
//
//type ReportCluster []ReportSystem
//
//type ReportSystem struct {
//	X      int    `json:"x"`
//	Y      int    `json:"y"`
//	Z      int    `json:"z"`
//	Id     string `json:"id"`
//	SysHab int    `json:"sys_hab"`
//	Orbits []struct {
//		Orbit int    `json:"orbit"`
//		PType string `json:"ptype"`
//		Hab   int    `json:"hab"`
//	} `json:"orbits"`
//}
//
//// ReportGame configuration
//type ReportGame struct {
//	Id           string         `json:"id"`
//	Description  string         `json:"description"`
//	NationsIndex []NationsIndex `json:"nations-index"`
//	TurnsIndex   []TurnsIndex   `json:"turns-index"`
//	Turn         int            // current turn in the game
//	Nations      []*ReportNation
//}
//
//// ReportNation configuration
//type ReportNation struct {
//	Store       string `json:"store"` // path to store data
//	Id          int    `json:"id"`
//	Name        string `json:"name"`
//	Description string `json:"description"`
//	Speciality  string `json:"speciality"`
//	Government  struct {
//		Kind string `json:"kind"`
//		Name string `json:"name"`
//	} `json:"government"`
//	HomePlanet struct {
//		Name     string `json:"name"`
//		Location struct {
//			X     int `json:"x"`
//			Y     int `json:"y"`
//			Z     int `json:"z"`
//			Orbit int `json:"orbit"`
//		} `json:"location"`
//	} `json:"homeworld"`
//	Skills   Skills
//	Colonies []*XColony
//	Ships    []*XShip
//}
//
//// Read loads a store from a JSON file.
//// It returns any errors.
//func (s *ReportNation) Read() error {
//	b, err := ioutil.ReadFile(filepath.Join(s.Store, "store.json"))
//	if err != nil {
//		return err
//	}
//	return json.Unmarshal(b, s)
//}
//
//// Write writes a store to a JSON file.
//// It returns any errors.
//func (s *ReportNation) Write() error {
//	if s.Store == "" {
//		return errors.New("missing nation store path")
//	}
//	b, err := json.MarshalIndent(s, "", "  ")
//	if err != nil {
//		return err
//	}
//	return ioutil.WriteFile(filepath.Join(s.Store, "store.json"), b, 0600)
//}
