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
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"
)

func (e *Engine) Report(spId int) error {
	// load the setup store
	var s Store
	b, err := os.ReadFile(filepath.Clean(filepath.Join(e.config.game.Store, fmt.Sprintf("%d", spId), "setup.json")))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// run the report for just the one nation
	for _, n := range e.config.nations.Index {
		if n.Id != spId {
			continue
		}
		log.Printf("reporting for %q %q (%q)\n", n.Id, n.Name, n.Store)

		reportFile := filepath.Clean(filepath.Join(n.Store, "report.txt"))
		log.Printf("reporting for %q %q (%q)\n", n.Id, n.Name, reportFile)

		w, err := os.OpenFile(reportFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}

		return s.Report(w, spId)
	}
	return nil
}

func (s *Store) Report(w io.Writer, spId int) error {
	p := message.NewPrinter(language.English)

	rptDate := time.Now().Format("2006/01/02")

	for _, player := range s.Players {
		_, _ = p.Fprintf(w, "Status Report\n")
		_, _ = p.Fprintf(w, "Game: %-6s    Turn: %5d    Player: %3d    Date: %s\n", s.Game.Id, s.Game.Turn, spId, rptDate)

		_, _ = p.Fprintf(w, "\nSpecies %3d ------------------------------------------------------------------\n", spId)
		_, _ = p.Fprintf(w, "  Bureaucracy:   %2d    Biology: %2d    Gravitics: %2d    LifeSupport: %2d\n",
			player.Skills.Bureaucracy, player.Skills.Biology, player.Skills.Gravitics, player.Skills.LifeSupport)
		_, _ = p.Fprintf(w, "  Manufacturing: %2d    Mining:  %2d    Military:  %2d    Shields:     %2d\n",
			player.Skills.Manufacturing, player.Skills.Mining, player.Skills.Military, player.Skills.Shields)

		for _, colony := range player.Colonies {
			_, _ = p.Fprintf(w, "\nColony Activity Report -------------------------------------------------------\n")
			kind := "unknown"
			switch colony.Kind {
			case "enclosed":
				kind = "ENCLOSED"
			case "open":
				kind = "OPEN"
			case "orbital":
				kind = "ORBITAL"
			default:
				panic(fmt.Sprintf("assert(colony.Kind != %q)", colony.Kind))
			}
			name := colony.Name
			if name == "" {
				name = "NOT NAMED"
			}
			colonyId, colonyKind := fmt.Sprintf("C%d", colony.Id), fmt.Sprintf("%s COLONY", kind)
			_, _ = p.Fprintf(w, "%-6s Tech:%2d  %14s: %-22s  System: %s #%d\n", colonyId, colony.TechLevel, colonyKind, name, colony.System, colony.Orbit)

			_, _ = p.Fprintf(w, "\nConstruction --------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Engines\n")
			_, _ = p.Fprintf(w, "    N/A\n")

			_, _ = p.Fprintf(w, "\nPeople --------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Group____________  Population_Units  Pay_____         CNGD/Turn  Ration__         FOOD/Turn\n")
			_, _ = p.Fprintf(w, "  Professional       %16d  %7.3f%%  %16d  %7.3f%%  %16d\n", colony.Population.PRO.Qty, colony.Population.PRO.Pay*100, colony.Population.PRO.TotalPay(), colony.Population.PRO.Ration*100, colony.Population.PRO.TotalRation())
			_, _ = p.Fprintf(w, "  Soldier            %16d  %7.3f%%  %16d  %7.3f%%  %16d\n", colony.Population.SLD.Qty, colony.Population.SLD.Pay*100, colony.Population.SLD.TotalPay(), colony.Population.SLD.Ration*100, colony.Population.SLD.TotalRation())
			_, _ = p.Fprintf(w, "  Unskilled          %16d  %7.3f%%  %16d  %7.3f%%  %16d\n", colony.Population.USK.Qty, colony.Population.USK.Pay*100, colony.Population.USK.TotalPay(), colony.Population.USK.Ration*100, colony.Population.USK.TotalRation())
			_, _ = p.Fprintf(w, "  Unemployed         %16d  %7.3f%%  %16d  %7.3f%%  %16d\n", colony.Population.UEM.Qty, colony.Population.UEM.Pay*100, colony.Population.UEM.TotalPay(), colony.Population.UEM.Ration*100, colony.Population.UEM.TotalRation())
			_, _ = p.Fprintf(w, "  ----------------   %16d  --------  %16d  --------  %16d\n", colony.Population.TotalPopulation(), colony.Population.TotalPay(), colony.Population.TotalRation())

			if colony.Population.Births == 0 && colony.Population.Deaths == 0 {
				colony.Population.Births = colony.Population.TotalPopulation() / 1600
				colony.Population.Deaths = colony.Population.Births
			}

			_, _ = p.Fprintf(w, "\n  Crew/Team________  Units___________\n")
			_, _ = p.Fprintf(w, "  Construction Crew  %16d\n", colony.Population.CNW)
			_, _ = p.Fprintf(w, "  Spy Team           %16d\n", colony.Population.SPY)

			_, _ = p.Fprintf(w, "\n  Changes__________  Population_Units\n")
			_, _ = p.Fprintf(w, "  Births             %16d\n", colony.Population.Births)
			_, _ = p.Fprintf(w, "  Non-Combat Deaths  %16d\n", colony.Population.Deaths)

			availSUs := 0
			operMUs, operEMUs, operSUs := 0, 0, 0
			_, _ = p.Fprintf(w, "\nOperational ------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Item-TL     Quantity          MUs  Hudnut         EMUs  SUs_Required\n")
			for _, unit := range colony.Operational {
				mu, emu := unit.MassUnits(), unit.EnclosedMassUnits()
				var sus int
				if unit.Name == "structural" || unit.Name == "light-structural" || unit.Name == "super-light-structural" {
					sus = 0
					availSUs += unit.Qty
				} else {
					switch colony.Kind {
					case "open":
						sus = emu * 1
					case "enclosed":
						sus = emu * 5
					case "orbital":
						sus = emu * 10
					}
				}
				_, _ = p.Fprintf(w, "  %-7s  %11d  %11d  %-6v  %11d  %12d\n", unit.Code(), unit.Qty, mu, unit.Hudnut(), emu, sus)
				operMUs += mu
				operEMUs += emu
				operSUs += sus
			}
			_, _ = p.Fprintf(w, "  Total                 %11d          %11d  %12d\n", operMUs, operEMUs, operSUs)

			storMUs, storEMUs, storSUs := 0, 0, 0
			_, _ = p.Fprintf(w, "\nStorage ----------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Item-TL     Quantity          MUs  Hudnut         EMUs  SUs_Required\n")
			for _, unit := range colony.Storage {
				mu, emu := unit.MassUnits(), unit.EnclosedMassUnits()
				var sus int
				switch colony.Kind {
				case "open":
					sus = emu * 1
				case "enclosed":
					sus = emu * 5
				case "orbital":
					sus = emu * 10
				}
				_, _ = p.Fprintf(w, "  %-7s  %11d  %11d  %-6v  %11d  %12d\n", unit.Code(), unit.Qty, mu, unit.Hudnut(), emu, sus)
				storMUs += mu
				storEMUs += emu
				storSUs += sus
			}
			_, _ = p.Fprintf(w, "  Total                 %11d          %11d  %12d\n", storMUs, storEMUs, storSUs)

			_, _ = p.Fprintf(w, "\nHull -------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Area_______           EMUs   SUs_Required  SUs_Available\n")
			_, _ = p.Fprintf(w, "  Operational  %13d  %13d\n", operEMUs, operSUs)
			_, _ = p.Fprintf(w, "  Storage      %13d  %13d\n", storEMUs, storSUs)
			_, _ = p.Fprintf(w, "  Total        %13d  %13d  %13d\n", operEMUs+storMUs, operSUs+storSUs, availSUs)

			//_, _ = p.Fprintf(w, "\nFarming --------------------------------------------------------------------------------------------------------\n")
			//_, _ = p.Fprintf(w, "  Group  Orders       Quantity  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")
			//if farm := colony.Units["FRM-1"]; farm != nil {
			//	group, techLevel := 1, 1
			//	qty := farm.Qty * 100 / 4
			//	_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
			//		group, "FOOD", farm.Qty, techLevel, int(math.Ceil(float64(farm.Qty)*0.5)), 1*farm.Qty, 3*farm.Qty, qty, qty, qty)
			//}

			_, _ = p.Fprintf(w, "\nFarming --------------------------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Group  Orders          Farms  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")
			for _, group := range colony.FarmGroups {
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.Qty) * 0.5))
					proLabor, uskLabor := 1*unit.Qty, 3*unit.Qty
					_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
						group.Id, group.Code(), unit.Qty, unit.TechLevel, fuelPerTurn, proLabor, uskLabor, unit.Stages[0], unit.Stages[1], unit.Stages[2])
				}
			}

			//_, _ = p.Fprintf(w, "\nMining ---------------------------------------------------------------------------------------------------------\n")
			//_, _ = p.Fprintf(w, "  Group  Orders        MIN_Qty  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")
			//if mine := colony.Units["MIN-1"]; mine != nil {
			//	techLevel := 1
			//	var no, qty int
			//	no = 50_000
			//	qty = no * 100 / 4
			//	_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
			//		1, "FUEL", no, techLevel, int(math.Ceil(float64(no)*0.5)), 1*no, 3*no, qty, qty, qty)
			//	no = mine.Qty
			//	qty = no * 100 / 4
			//	_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
			//		2, "METS", no, techLevel, int(math.Ceil(float64(no)*0.5)), 1*no, 3*no, qty, qty, qty)
			//	no = mine.Qty
			//	qty = no * 100 / 4
			//	_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
			//		3, "NMET", no, techLevel, int(math.Ceil(float64(no)*0.5)), 1*no, 3*no, qty, qty, qty)
			//}

			_, _ = p.Fprintf(w, "\nMining ---------------------------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  Group  Orders          Mines  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")
			for _, group := range colony.MiningGroups {
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.Qty) * 0.5))
					proLabor, uskLabor := 1*unit.Qty, 3*unit.Qty
					_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
						group.Id, group.Code(), unit.Qty, unit.TechLevel, fuelPerTurn, proLabor, uskLabor, unit.Stages[0], unit.Stages[1], unit.Stages[2])
				}
			}

			//_, _ = p.Fprintf(w, "\nProduction -----------------------------------------------------------------------------------------------------\n")
			//_, _ = p.Fprintf(w, "Input ----\n")
			//_, _ = p.Fprintf(w, "  Group  Orders       Quantity  TL  Ingest/Turn    METS/Unit    NMTS/Unit    METS/Turn    NMTS/Turn   Units/Turn\n")
			//if fact := colony.Units["FCT-1"]; fact != nil {
			//	group, techLevel := 1, 1
			//
			//	// FACT-1 can ingest a maximum of 5 MU of resources per turn
			//	ingestTurn := fact.Qty * 5
			//
			//	// CNGD costs 0.2 METS 0.4 NMET
			//	metsUnit, nmetUnit := 0.2, 0.4
			//
			//	// quantity is units per turn
			//	qty := int(math.Floor(float64(ingestTurn) / (metsUnit + nmetUnit)))
			//	metsTurn := int(metsUnit * float64(qty))
			//	nmetTurn := int(nmetUnit * float64(qty))
			//
			//	_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d    %9.3f    %9.3f  %11d  %11d  %11d\n",
			//		group, "CNGD", fact.Qty, techLevel, ingestTurn, metsUnit, nmetUnit, metsTurn, nmetTurn, qty)
			//}
			//
			//_, _ = p.Fprintf(w, "Output ---\n")
			//_, _ = p.Fprintf(w, "  Group  Orders       Quantity  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")
			//if fact := colony.Units["FCT-1"]; fact != nil {
			//	group, techLevel := 1, 1
			//	var pro, usk int
			//	if fact.Qty >= 50_000 {
			//		pro, usk = 1*fact.Qty, 3*fact.Qty
			//	} else if fact.Qty >= 5_000 {
			//		pro, usk = 2*fact.Qty, 6*fact.Qty
			//	} else if fact.Qty >= 500 {
			//		pro, usk = 3*fact.Qty, 9*fact.Qty
			//	} else if fact.Qty >= 50 {
			//		pro, usk = 4*fact.Qty, 12*fact.Qty
			//	} else if fact.Qty >= 5 {
			//		pro, usk = 5*fact.Qty, 15*fact.Qty
			//	} else {
			//		pro, usk = 6*fact.Qty, 18*fact.Qty
			//	}
			//
			//	// FACT-1 can ingest 20 MU of resources per YEAR
			//	ingest := fact.Qty * 20
			//	// CNGD costs 0.2 METS 0.4 NMET
			//	mets, nmet := 0.2, 0.4
			//	// quantity is units per turn
			//	qty := int(math.Floor(float64(ingest)/(mets+nmet))) / 4
			//
			//	_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
			//		group, "CNGD", fact.Qty, techLevel, int(math.Ceil(float64(fact.Qty)*0.5)), pro, usk, qty, qty, qty)
			//}

			_, _ = p.Fprintf(w, "\nProduction -----------------------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "Input ----\n")
			_, _ = p.Fprintf(w, "  Group  Orders      Factories  TL  Ingest/Turn    METS/Unit    NMTS/Unit    METS/Turn    NMTS/Turn   Units/Turn\n")
			for _, group := range colony.FactoryGroups {
				for _, unit := range group.Units {
					u := Unit{Name: "factory", TechLevel: unit.TechLevel, Qty: unit.Qty}

					metsUnit, nmetUnit := Unit{Name: group.Name}.RawMaterials()

					// quantity is units per turn
					qty := int(math.Floor(float64(u.IngestPerTurn()) / (metsUnit + nmetUnit)))
					metsTurn := int(metsUnit * float64(qty))
					nmetTurn := int(nmetUnit * float64(qty))

					_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d    %9.3f    %9.3f  %11d  %11d  %11d\n",
						group.Id, group.Code(), u.Qty, u.TechLevel, u.IngestPerTurn(), metsUnit, nmetUnit, metsTurn, nmetTurn, qty)
				}
			}

			_, _ = p.Fprintf(w, "Output ---\n")
			_, _ = p.Fprintf(w, "  Group  Orders      Factories  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")
			for _, group := range colony.FactoryGroups {
				for _, unit := range group.Units {
					u := Unit{Name: "factory", TechLevel: unit.TechLevel, Qty: unit.Qty}

					proLabor, uskLabor := u.LaborPerTurn()

					_, _ = p.Fprintf(w, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
						group.Id, group.Code(), unit.Qty, unit.TechLevel, u.Qty*u.FuelPerTurn(), proLabor, uskLabor, unit.Stages[0], unit.Stages[1], unit.Stages[2])
				}
			}

			_, _ = p.Fprintf(w, "\nEspionage --------------------------------------------------------------------\n")
			_, _ = p.Fprintf(w, "  No activity.\n")
		}

		orbits := []Orbit{
			{Id: 1, Kind: "terrestrial", HabitabilityNumber: 0},
			{Id: 2, Kind: "terrestrial", HabitabilityNumber: 0},
			{Id: 3, Kind: "terrestrial", HabitabilityNumber: 25},
			{Id: 4, Kind: "terrestrial", HabitabilityNumber: 0},
			{Id: 5, Kind: "asteroid belt", HabitabilityNumber: 0},
			{Id: 6, Kind: "gas giant", HabitabilityNumber: 0},
			{Id: 7, Kind: "gas giant", HabitabilityNumber: 0},
			{Id: 8, Kind: "terrestrial", HabitabilityNumber: 0},
			{Id: 9, Kind: "terrestrial", HabitabilityNumber: 0},
			{Id: 10, Kind: "asteroid belt", HabitabilityNumber: 0},
		}
		_, _ = p.Fprintf(w, "\nSurvey Report ----------------------------------------------------------------\n")
		for _, orbit := range orbits {
			name, controlledBy := "NOT NAMED", "N/A"
			if orbit.Id == 3 {
				name, controlledBy = "My Homeworld", "SP018"
			}
			_, _ = p.Fprintf(w, "  System %s %-3s   %-24s    Controlled By: %s\n", "0/0/0", fmt.Sprintf("#%d", orbit.Id), name, controlledBy)
			_, _ = p.Fprintf(w, "    Kind: %-13s    Habitability: %2d\n", orbit.Kind, orbit.HabitabilityNumber)
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

type Store struct {
	Game struct {
		Id   string `json:"id"`
		Turn int    `json:"turn"`
	} `json:"game"`
	Players map[string]*Player `json:"players,omitempty"`
}

type Player struct {
	Skills   *Skills   `json:"skills,omitempty"`
	Colonies []*Colony `json:"colonies,omitempty"`
	Ships    []*Ship   `json:"ships,omitempty"`
}

type PopUnit struct {
	Code   string  `json:"code"`
	Qty    int     `json:"qty,omitempty"`
	Pay    float64 `json:"pay,omitempty"`
	Ration float64 `json:"ration,omitempty"`
}

// TotalPay assumes that the base rates are per unit of population
//  PROFESSIONAL      0.375 CONSUMER GOODS
//  SOLDIER           0.250 CONSUMER GOODS
//  UNSKILLED WORKER  0.125 CONSUMER GOODS
//  UNEMPLOYABLE      0.000 CONSUMER GOODS
func (p *PopUnit) TotalPay() int {
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

// TotalRation assumes that base ration is 0.25 food units per unit of population
func (p *PopUnit) TotalRation() int {
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

type Unit struct {
	Name      string `json:"name"`
	TechLevel int    `json:"tech-level,omitempty"`
	Qty       int    `json:"qty,omitempty"`
	Stowed    bool   `json:"stowed,omitempty"`
	code      string
}

func (u Unit) Code() string {
	if u.code == "" {
		switch u.Name {
		case "anti-missile":
			u.code = fmt.Sprintf("ANM-%d", u.TechLevel)
		case "assault-craft":
			u.code = fmt.Sprintf("ASC-%d", u.TechLevel)
		case "assault-weapon":
			u.code = fmt.Sprintf("ASW-%d", u.TechLevel)
		case "automation":
			u.code = fmt.Sprintf("AUT-%d", u.TechLevel)
		case "consumer-goods":
			u.code = "CNGD"
		case "energy-shield":
			u.code = fmt.Sprintf("ESH-%d", u.TechLevel)
		case "energy-weapon":
			u.code = fmt.Sprintf("EWP-%d", u.TechLevel)
		case "factory":
			u.code = fmt.Sprintf("FCT-%d", u.TechLevel)
		case "farm":
			u.code = fmt.Sprintf("FRM-%d", u.TechLevel)
		case "food":
			u.code = "FOOD"
		case "fuel":
			u.code = "FUEL"
		case "gold":
			u.code = "GOLD"
		case "hyper-drive":
			u.code = fmt.Sprintf("HDR-%d", u.TechLevel)
		case "life-support":
			u.code = fmt.Sprintf("LSP-%d", u.TechLevel)
		case "light-structural":
			u.code = "LTSU"
		case "metallics":
			u.code = "MTLS"
		case "military-robots":
			u.code = fmt.Sprintf("MLR-%d", u.TechLevel)
		case "military-supplies":
			u.code = "MLSP"
		case "mine":
			u.code = fmt.Sprintf("MIN-%d", u.TechLevel)
		case "missile":
			u.code = fmt.Sprintf("MSS-%d", u.TechLevel)
		case "missile-launcher":
			u.code = fmt.Sprintf("MSL-%d", u.TechLevel)
		case "non-metallics":
			u.code = "NMTS"
		case "sensor":
			u.code = fmt.Sprintf("SNR-%d", u.TechLevel)
		case "space-drive":
			u.code = fmt.Sprintf("SDR-%d", u.TechLevel)
		case "structural":
			u.code = "STUN"
		case "super-light-structural":
			u.code = "SLSU"
		case "transport":
			u.code = fmt.Sprintf("TPT-%d", u.TechLevel)
		default:
			panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
		}
	}
	return u.code
}

func (u Unit) FuelPerTurn() int {
	switch u.Name {
	case "anti-missile":
		return 0
	case "assault-craft":
		return 0
	case "assault-weapon":
		return 2 * u.TechLevel
	case "automation":
		return 0
	case "consumer-goods":
		return 0
	case "energy-shield":
		return 0
	case "energy-weapon":
		return 0
	case "factory":
		return int(math.Ceil(float64(u.TechLevel) / 2))
	case "farm":
		if u.TechLevel < 6 {
			return int(math.Ceil(float64(u.TechLevel) / 2))
		}
		return u.TechLevel
	case "hyper-drive":
		return 0
	case "life-support":
		return 1
	case "light-structural":
		return 0
	case "military-robots":
		return 0
	case "military-supplies":
		return 0
	case "mine":
		return int(math.Ceil(float64(u.TechLevel) / 2))
	case "missile":
		return 0
	case "missile-launcher":
		return 0
	case "sensor":
		return int(math.Ceil(float64(u.TechLevel) / 20))
	case "space-drive":
		return 0
	case "structural":
		return 0
	case "super-light-structural":
		return 0
	case "transport":
		return int(math.Ceil(float64(u.TechLevel*u.TechLevel) / 10))
	default:
		panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
	}
}

func (u Unit) IngestPerTurn() int {
	if u.Name != "factory" {
		return 0
	}
	// can ingest 20 MU of resources per tech-level per YEAR
	return u.Qty * u.TechLevel * 20 / 4
}

func (u Unit) LaborPerTurn() (pro, usk int) {
	if u.Name != "factory" {
		return 1 * u.Qty, 3 * u.Qty
	} else if u.Qty >= 50_000 {
		return 1 * u.Qty, 3 * u.Qty
	} else if u.Qty >= 5_000 {
		return 2 * u.Qty, 6 * u.Qty
	} else if u.Qty >= 500 {
		return 3 * u.Qty, 9 * u.Qty
	} else if u.Qty >= 50 {
		return 4 * u.Qty, 12 * u.Qty
	} else if u.Qty >= 5 {
		return 5 * u.Qty, 15 * u.Qty
	}
	return 6 * u.Qty, 18 * u.Qty
}

func (u Unit) EnclosedMassUnits() int {
	mus := u.MassUnits()
	if !u.Hudnut() || !u.Stowed {
		return mus
	}
	return int(math.Ceil(float64(mus) / 2))
}

func (u Unit) MassUnits() int {
	switch u.Name {
	case "anti-missile":
		return 4 * u.TechLevel * u.Qty
	case "assault-craft":
		return 5 * u.TechLevel * u.Qty
	case "assault-weapon":
		return 2 * u.TechLevel * u.Qty
	case "automation":
		return 4 * u.TechLevel * u.Qty
	case "consumer-goods":
		return int(math.Ceil(0.6 * float64(u.Qty)))
	case "energy-shield":
		return 50 * u.TechLevel * u.Qty
	case "energy-weapon":
		return 10 * u.TechLevel * u.Qty
	case "factory":
		return (12 + 2*u.TechLevel) * u.Qty
	case "farm":
		return (6 + 2*u.TechLevel) * u.Qty
	case "food":
		return 6 * u.Qty
	case "fuel":
		return 1 * u.Qty
	case "gold":
		return 1 * u.Qty
	case "hyper-drive":
		return 45 * u.TechLevel * u.Qty
	case "life-support":
		return 8 * u.TechLevel * u.Qty
	case "light-structural":
		return int(math.Ceil(0.05 * float64(u.Qty)))
	case "metallics":
		return 1 * u.Qty
	case "military-robots":
		return (20 + 2*u.TechLevel) * u.Qty
	case "military-supplies":
		return int(math.Ceil(0.04 * float64(u.Qty)))
	case "mine":
		return (10 + 2*u.TechLevel) * u.Qty
	case "missile":
		return 4 * u.TechLevel * u.Qty
	case "missile-launcher":
		return 25 * u.TechLevel * u.Qty
	case "non-metallics":
		return 1 * u.Qty
	case "sensor":
		return 40 * u.TechLevel * u.Qty
	case "space-drive":
		return 25 * u.TechLevel * u.Qty
	case "structural":
		return int(math.Ceil(0.5 * float64(u.Qty)))
	case "super-light-structural":
		return int(math.Ceil(0.005 * float64(u.Qty)))
	case "transport":
		return int(math.Ceil(0.1 * float64(u.TechLevel*u.TechLevel*u.Qty)))
	default:
		panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
	}
}

func (u Unit) Assembled() string {
	if !u.Stowed {
		return "??"
	}
	return "???"
}

// Hudnut is borrowed from the Sniglet for the leftover
// bolts and such from a some-assembly-required project.
// Returns true if the unit can be disassembled for storage.
func (u Unit) Hudnut() bool {
	switch u.Name {
	case "anti-missile":
		return false
	case "assault-craft":
		return false
	case "assault-weapon":
		return false
	case "automation":
		return true
	case "consumer-goods":
		return false
	case "energy-shield":
		return true
	case "energy-weapon":
		return true
	case "factory":
		return true
	case "farm":
		return true
	case "food":
		return false
	case "fuel":
		return false
	case "gold":
		return false
	case "hyper-drive":
		return true
	case "life-support":
		return true
	case "light-structural":
		return true
	case "metallics":
		return false
	case "military-robots":
		return false
	case "military-supplies":
		return false
	case "mine":
		return true
	case "missile":
		return false
	case "missile-launcher":
		return true
	case "non-metallics":
		return false
	case "sensor":
		return true
	case "space-drive":
		return true
	case "structural":
		return true
	case "super-light-structural":
		return true
	case "transport":
		return false
	default:
		panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
	}
}

func (u Unit) Operational() string {
	if !u.Stowed {
		return "yes"
	}
	switch u.code {
	case "FCT":
		return "no"
	case "MIN":
		return "no"
	}
	return "yes"
}

func (u Unit) RawMaterials() (mets, nmts float64) {
	tl := float64(u.TechLevel)
	switch u.Name {
	case "anti-missile":
		return 2 * tl, 2 * tl
	case "assault-craft":
		return 3 * tl, 2 * tl
	case "assault-weapon":
		return 1 * tl, 1 * tl
	case "automation":
		return 2 * tl, 2 * tl
	case "consumer-goods":
		return 0.2, 0.4
	case "energy-shield":
		return 25 * tl, 25 * tl
	case "energy-weapon":
		return 5 * tl, 5 * tl
	case "factory":
		return 8 * tl, 4 * tl
	case "farm":
		return 4 + tl, 2 + tl
	case "hyper-drive":
		return 25 * tl, 20 * tl
	case "life-support":
		return 3 * tl, 5 * tl
	case "light-structural":
		return 0.01, 0.04
	case "military-robots":
		return 10 * tl, 10 * tl
	case "military-supplies":
		return 0.02, 0.02
	case "mine":
		return 5 + tl, 5 + tl
	case "missile":
		return 2 * tl, 2 * tl
	case "missile-launcher":
		return 15 * tl, 10 * tl
	case "sensor":
		return 10 * tl, 20 * tl
	case "space-drive":
		return 15 * tl, 10 * tl
	case "structural":
		return 0.1, 0.4
	case "super-light-structural":
		return 0.001, 0.004
	case "transport":
		return 3 * tl, tl
	}
	panic(fmt.Sprintf("assert(unit.name != %q)", u.Name))
}

func UnitAttributes(name string, techLevel int) (mets, nmts, totalMassUnits, fuelPerTurn, fuelPerCombatRound float64) {
	tl := float64(techLevel)
	switch name {
	case "anti-missile":
		return 2 * tl, 2 * tl, 4 * tl, 0, 0
	case "assault-craft":
		return 3 * tl, 2 * tl, 5 * tl, 0, 0.1
	case "assault-weapon":
		return 1 * tl, 1 * tl, 2 * tl, 2 * tl * tl, 0
	case "automation":
		return 2 * tl, 2 * tl, 4 * tl, 0, 0
	case "consumer-goods":
		return 0.2, 0.4, 0.6, 0, 0
	case "energy-shield":
		return 25 * tl, 25 * tl, 50 * tl, 0, 10 * tl
	case "energy-weapon":
		return 5 * tl, 5 * tl, 10 * tl, 0, 4 * tl
	case "factory":
		return 8 * tl, 4 * tl, 12 + 2*tl, 0.5 * tl, 4 * tl
	case "farm":
		if techLevel == 1 {
			return 4 + tl, 2 + tl, 6 + 2*tl, 0.5 * tl, 0
		} else if techLevel < 6 {
			return 4 + tl, 4 + tl, 6 + 2*tl, 0.5 * tl, 0
		}
		return 4 + tl, 2 + tl, 6 + 2*tl, tl, 0
	case "hyper-drive":
		return 25 * tl, 20 * tl, 45 * tl, 0, 0
	case "life-support":
		return 3 * tl, 5 * tl, 8 * tl, 1, 0
	case "light-structural":
		return 0.01, 0.04, 0.05, 0, 0
	case "military-robots":
		return 10 * tl, 10 * tl, 20 + 2*tl, 0, 0
	case "military-supplies":
		return 0.02, 0.02, 0.04, 0, 0
	case "mine":
		return 5 + tl, 5 + tl, 10 + (2 * tl), 0.5 * tl, 0
	case "missile":
		return 2 * tl, 2 * tl, 4 * tl, 0, 0
	case "missile-launcher":
		return 15 * tl, 10 * tl, 25 * tl, 0, 0
	case "sensor":
		return 10 * tl, 20 * tl, 40 * tl, tl / 20, 0
	case "space-drive":
		return 15 * tl, 10 * tl, 25 * tl, 0, tl * tl
	case "structural":
		return 0.1, 0.4, 0.5, 0, 0
	case "super-light-structural":
		return 0.001, 0.004, 0.005, 0, 0
	case "transport":
		return 3 * tl, tl, 4 * tl, 0.1 * tl * tl, 0.01 * tl * tl
	}
	panic(fmt.Sprintf("assert(unit.name != %q)", name))
}

type Group struct {
	Id        int          `json:"id"`
	Name      string       `json:"name"`
	TechLevel int          `json:"tech-level,omitempty"`
	Units     []*GroupUnit `json:"units,omitempty"`
	code      string
}

func (g Group) Code() string {
	if g.code == "" {
		g.code = Unit{Name: g.Name, TechLevel: g.TechLevel}.Code()
	}
	return g.code
}

type GroupUnit struct {
	TechLevel int   `json:"tech-level,omitempty"`
	Qty       int   `json:"qty,omitempty"`
	Stages    []int `json:"stages,omitempty"`
}

type Colony struct {
	Id            int         `json:"id"`
	Name          string      `json:"name,omitempty"`
	System        string      `json:"system"`
	Orbit         int         `json:"orbit"`
	Kind          string      `json:"kind"`
	TechLevel     int         `json:"tech-level"`
	Population    *Population `json:"population,omitempty"`
	Operational   []*Unit     `json:"operational,omitempty"`
	Storage       []*Unit     `json:"storage,omitempty"`
	FactoryGroups []*Group    `json:"factory-groups"`
	FarmGroups    []*Group    `json:"farm-groups,omitempty"`
	MiningGroups  []*Group    `json:"mining-groups,omitempty"`
}

type Ship struct {
	Id          int         `json:"id"`
	Name        string      `json:"name,omitempty"`
	TechLevel   int         `json:"tech-level"`
	Population  *Population `json:"population,omitempty"`
	Operational []*Unit     `json:"units,omitempty"`
	Storage     []*Unit     `json:"storage,omitempty"`
	FarmGroups  []*Group    `json:"farm-groups,omitempty"`
}

type Orbit struct {
	Id                 int
	Name               string
	Kind               string
	HabitabilityNumber int
}

type Population struct {
	PRO    PopUnit `json:"pro,omitempty"`
	SLD    PopUnit `json:"sld,omitempty"`
	USK    PopUnit `json:"usk,omitempty"`
	UEM    PopUnit `json:"uem,omitempty"`
	CNW    int     `json:"cnw,omitempty"`
	SPY    int     `json:"spy,omitempty"`
	Births int     `json:"births,omitempty"`
	Deaths int     `json:"deaths,omitempty"`
}

func (p *Population) TotalPay() int {
	if p == nil {
		return 0
	}
	return p.PRO.TotalPay() + p.SLD.TotalPay() + p.USK.TotalPay() + p.UEM.TotalPay()

}

func (p *Population) TotalPopulation() int {
	if p == nil {
		return 0
	}
	return p.PRO.Qty + p.SLD.Qty + p.USK.Qty + p.UEM.Qty
}

func (p *Population) TotalRation() int {
	if p == nil {
		return 0
	}
	return p.PRO.TotalRation() + p.SLD.TotalRation() + p.USK.TotalRation() + p.UEM.TotalRation()
}

type Skills struct {
	Biology       int `json:"biology,omitempty"`
	Bureaucracy   int `json:"bureaucracy,omitempty"`
	Gravitics     int `json:"gravitics,omitempty"`
	LifeSupport   int `json:"life-support,omitempty"`
	Manufacturing int `json:"manufacturing,omitempty"`
	Military      int `json:"military,omitempty"`
	Mining        int `json:"mining,omitempty"`
	Shields       int `json:"shields,omitempty"`
}

type Cluster []System

type System struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Z      int    `json:"z"`
	Id     string `json:"id"`
	SysHab int    `json:"sys_hab"`
	Orbits []struct {
		Orbit int    `json:"orbit"`
		PType string `json:"ptype"`
		Hab   int    `json:"hab"`
	} `json:"orbits"`
}
