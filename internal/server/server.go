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

package server

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/jsonwt"
	"github.com/mdhender/jsonwt/signers"
	"github.com/mdhender/wraith/internal/way"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	http.Server

	router *way.Router
	www    struct {
		host string
		port string
	}
	game struct {
		source string
		data   *Game
	}
	authn struct {
		source  string
		Version string `json:"version"`
		Bcrypt  struct {
			MinCost int `json:"min_cost"`
		} `json:"bcrypt"`
		Users map[string]*user `json:"users"` // key is user handle
	}
	jwt struct {
		source   string
		Version  string `json:"version"`
		TTLHours int    `json:"ttl-hours"`
		Key      struct {
			Name   string `json:"name"`
			Secret string `json:"secret"`
		} `json:"key"`
		factory *jsonwt.Factory
		ttl     time.Duration
	}
}

type claims struct {
	Species string   `json:"species"`
	Roles   []string `json:"roles"`
}

type user struct {
	Id           string   `json:"id"`
	Handle       string   `json:"handle"`
	Roles        []string `json:"roles"`
	Secret       string   `json:"secret"`
	HashedSecret string   `json:"hashed-secret"`
}

func New(opts ...Option) (*Server, error) {
	// create a server with default values
	s := &Server{}
	s.authn.Users = make(map[string]*user)
	s.router = way.NewRouter()
	s.www.host, s.www.port = "", "8080"

	s.Addr = net.JoinHostPort(s.www.host, s.www.port)
	s.MaxHeaderBytes = 1 << 20 // 1mb?
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second

	// apply the list of options to the server
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	if b, err := os.ReadFile(s.authn.source); err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &s.authn); err != nil {
		return nil, err
	} else {
		for id, u := range s.authn.Users {
			if u.Id != id {
				u.Id = id
			}
			if u.Handle == "" {
				u.Handle = u.Id
			}
		}
	}

	if b, err := os.ReadFile(s.jwt.source); err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &s.jwt); err != nil {
		return nil, err
	} else if hsSigner, err := signers.NewHS256([]byte(s.jwt.Key.Secret)); err != nil {
		return nil, err
	} else {
		s.jwt.factory = jsonwt.NewFactory(s.jwt.Key.Name, hsSigner)
		s.jwt.ttl = time.Hour * time.Duration(s.jwt.TTLHours)
	}

	if b, err := os.ReadFile(s.game.source); err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &s.game.data); err != nil {
		return nil, err
	}

	//s.router.HandleFunc("GET", "/ui/add-user", s.handleGetAddUser)
	//s.router.HandleFunc("POST", "/ui/add-user", s.handlePostAddUser)
	s.router.HandleFunc("GET", "/ui", s.authenticatedOnly(s.handleGetIndex))
	s.router.HandleFunc("GET", "/ui/games/PT-1/players/SP1/report", s.authenticatedOnly(s.handleGetReport))
	s.router.HandleFunc("GET", "/ui/login", s.handleGetLogin)
	s.router.HandleFunc("POST", "/ui/login", s.handlePostLogin)
	s.router.HandleFunc("GET", "/ui/logout", s.handleLogout)
	s.router.HandleFunc("POST", "/ui/logout", s.handleLogout)

	return s, nil
}

func (s *Server) Mux() http.Handler {
	return s.router
}

func (s *Server) handleGetIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	page := fmt.Sprintf(`<body>
				<h1>Wraith UI</h1>
			</body>`)
	_, _ = w.Write([]byte(page))
}

func (s *Server) handleGetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	page := &strings.Builder{}
	page.WriteString(`<body><h1>Wraith SP1 Report</h1><code><pre>`)

	// run the report for just the one nation
	for _, n := range s.game.data.Nations {
		if n.Name != "SP1" {
			continue
		}
		log.Printf("reporting for %d %q (%q)\n", n.Id, n.Name, s.game.source)

		p := message.NewPrinter(language.English)

		rptDate := time.Now().Format("2006/01/02")

		_, _ = p.Fprintf(page, "Status Report\n")
		_, _ = p.Fprintf(page, "Game: %-6s    Turn: %5d    Player: %3d    Date: %s\n", s.game.data.Id, s.game.data.Turn, n.Id, rptDate)

		_, _ = p.Fprintf(page, "\nNation %s ------------------------------------------------------------------\n", n.Name)
		_, _ = p.Fprintf(page, "  Bureaucracy:   %2d    Biology: %2d    Gravitics: %2d    LifeSupport: %2d\n",
			n.Skills.Bureaucracy, n.Skills.Biology, n.Skills.Gravitics, n.Skills.LifeSupport)
		_, _ = p.Fprintf(page, "  Manufacturing: %2d    Mining:  %2d    Military:  %2d    Shields:     %2d\n",
			n.Skills.Manufacturing, n.Skills.Mining, n.Skills.Military, n.Skills.Shields)

		for _, colony := range n.Colonies {
			_, _ = p.Fprintf(page, "\nColony Activity Report -------------------------------------------------------\n")
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
			name := "UNNAMED" //colony.Name
			if name == "" {
				name = "NOT NAMED"
			}
			colonyId, colonyKind := fmt.Sprintf("C%d", colony.Id), fmt.Sprintf("%s COLONY", kind)
			_, _ = p.Fprintf(page, "%-6s Tech:%2d  %14s: %-22s  System: %d/%d/%d #%d\n", colonyId, colony.TechLevel, colonyKind, name, colony.Location.X, colony.Location.Y, colony.Location.Z, colony.Location.Orbit)

			_, _ = p.Fprintf(page, "\nConstruction --------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "  Engines\n")
			_, _ = p.Fprintf(page, "    N/A\n")

			_, _ = p.Fprintf(page, "\nPeople --------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "  Group____________  Population_Units  Pay_____         CNGD/Turn  Ration__         FOOD/Turn\n")
			_, _ = p.Fprintf(page, "  Professional       %16d  %7.3f%%  %16d  %7.3f%%  %16d\n", colony.Population.Professional.Qty, colony.Population.Professional.Pay*100, colony.Population.Professional.TotalPay(), colony.Population.Professional.Ration*100, colony.Population.Professional.TotalRation())
			_, _ = p.Fprintf(page, "  Soldier            %16d  %7.3f%%  %16d  %7.3f%%  %16d\n", colony.Population.Soldier.Qty, colony.Population.Soldier.Pay*100, colony.Population.Soldier.TotalPay(), colony.Population.Soldier.Ration*100, colony.Population.Soldier.TotalRation())
			_, _ = p.Fprintf(page, "  Unskilled          %16d  %7.3f%%  %16d  %7.3f%%  %16d\n", colony.Population.Unskilled.Qty, colony.Population.Unskilled.Pay*100, colony.Population.Unskilled.TotalPay(), colony.Population.Unskilled.Ration*100, colony.Population.Unskilled.TotalRation())
			_, _ = p.Fprintf(page, "  Unemployed         %16d  %7.3f%%  %16d  %7.3f%%  %16d\n", colony.Population.Unemployed.Qty, colony.Population.Unemployed.Pay*100, colony.Population.Unemployed.TotalPay(), colony.Population.Unemployed.Ration*100, colony.Population.Unemployed.TotalRation())
			_, _ = p.Fprintf(page, "  ----------------   %16d  --------  %16d  --------  %16d\n", colony.TotalPopulation(), colony.TotalPay(), colony.TotalRation())

			if colony.Population.Births == 0 && colony.Population.Deaths == 0 {
				colony.Population.Births = colony.TotalPopulation() / 1600
				colony.Population.Deaths = colony.Population.Births
			}

			_, _ = p.Fprintf(page, "\n  Crew/Team________  Units___________\n")
			_, _ = p.Fprintf(page, "  Construction Crew  %16d\n", colony.Population.ConstructionCrews)
			_, _ = p.Fprintf(page, "  Spy Team           %16d\n", colony.Population.SpyTeams)

			_, _ = p.Fprintf(page, "\n  Changes__________  Population_Units\n")
			_, _ = p.Fprintf(page, "  Births             %16d\n", colony.Population.Births)
			_, _ = p.Fprintf(page, "  Non-Combat Deaths  %16d\n", colony.Population.Deaths)

			availSUs := 0
			operMUs, operEMUs, operSUs := 0, 0, 0
			_, _ = p.Fprintf(page, "\nOperational ------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "  Item-TL     Quantity          MUs  Hudnut         EMUs  SUs_Required\n")
			for _, unit := range colony.Inventory {
				if unit.OperationalQty == 0 {
					continue
				}
				mu, emu := unit.MassUnits(), unit.EnclosedMassUnits()
				var sus int
				if unit.Name == "structural" || unit.Name == "light-structural" || unit.Name == "super-light-structural" {
					sus = 0
					availSUs += unit.OperationalQty
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
				_, _ = p.Fprintf(page, "  %-7s  %11d  %11d  %-6v  %11d  %12d\n", unit.Code, unit.OperationalQty, mu, false, emu, sus)
				operMUs += mu
				operEMUs += emu
				operSUs += sus
			}
			_, _ = p.Fprintf(page, "  Total                 %11d          %11d  %12d\n", operMUs, operEMUs, operSUs)

			storMUs, storEMUs, storSUs := 0, 0, 0
			_, _ = p.Fprintf(page, "\nStorage ----------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "  Item-TL     Quantity          MUs  Hudnut         EMUs  SUs_Required\n")
			for _, unit := range colony.Inventory {
				if unit.StowedQty == 0 {
					continue
				}
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
				_, _ = p.Fprintf(page, "  %-7s  %11d  %11d  %-6v  %11d  %12d\n", unit.Code, unit.StowedQty, mu, false, emu, sus)
				storMUs += mu
				storEMUs += emu
				storSUs += sus
			}
			_, _ = p.Fprintf(page, "  Total                 %11d          %11d  %12d\n", storMUs, storEMUs, storSUs)

			_, _ = p.Fprintf(page, "\nHull -------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "  Area_______           EMUs   SUs_Required  SUs_Available\n")
			_, _ = p.Fprintf(page, "  Operational  %13d  %13d\n", operEMUs, operSUs)
			_, _ = p.Fprintf(page, "  Storage      %13d  %13d\n", storEMUs, storSUs)
			_, _ = p.Fprintf(page, "  Total        %13d  %13d  %13d\n", operEMUs+storMUs, operSUs+storSUs, availSUs)

			_, _ = p.Fprintf(page, "\nFarming --------------------------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "  Group  Orders          Farms  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")
			for _, group := range colony.FarmGroups {
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.Qty) * 0.5))
					proLabor, uskLabor := 1*unit.Qty, 3*unit.Qty
					_, _ = p.Fprintf(page, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
						group.Id, "???", unit.Qty, unit.TechLevel, fuelPerTurn, proLabor, uskLabor, unit.Stages[0], unit.Stages[1], unit.Stages[2])
				}
			}

			_, _ = p.Fprintf(page, "\nMining ---------------------------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "  Group  Orders          Mines  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")
			for _, group := range colony.MiningGroups {
				for _, unit := range group.Units {
					fuelPerTurn := int(math.Ceil(float64(unit.Qty) * 0.5))
					proLabor, uskLabor := 1*unit.Qty, 3*unit.Qty
					_, _ = p.Fprintf(page, "  #%4d  %6s  %13d  %2d  %11d  %11d  %11d  %11d  %11d  %11d\n",
						group.Id, "???", unit.Qty, unit.TechLevel, fuelPerTurn, proLabor, uskLabor, unit.Stages[0], unit.Stages[1], unit.Stages[2])
				}
			}

			_, _ = p.Fprintf(page, "\nProduction -----------------------------------------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "Input ----\n")
			_, _ = p.Fprintf(page, "  Group  Orders      Factories  TL  Ingest/Turn    METS/Unit    NMTS/Unit    METS/Turn    NMTS/Turn   Units/Turn\n")

			_, _ = p.Fprintf(page, "Output ---\n")
			_, _ = p.Fprintf(page, "  Group  Orders      Factories  TL    FUEL/Turn    PRO_Labor    USK_Labor      Stage_1      Stage_2      Stage_3\n")

			_, _ = p.Fprintf(page, "\nEspionage --------------------------------------------------------------------\n")
			_, _ = p.Fprintf(page, "  No activity.\n")
		}

		_, _ = p.Fprintf(page, "\nMarket Report ---------------------------------------------------------------------\n")
		_, _ = p.Fprintf(page, "  News: *** 6 years ago, scientists discovered a signal broadcasting from the 10th\n")
		_, _ = p.Fprintf(page, "        orbit. 4 years ago, they decoded the message and recovered plans for an\n")
		_, _ = p.Fprintf(page, "        in-system engine (the \"space-drive\") and a faster-than-light engine\n")
		_, _ = p.Fprintf(page, "        (the \"hyper-drive\"). Work on both has recently completed.\n")
		_, _ = p.Fprintf(page, "  News: *** Moments after the hyper-drive was successfully tested in the orbital\n")
		_, _ = p.Fprintf(page, "        colony, the broadcast from the 10th orbit stopped.\n")

		_, _ = p.Fprintf(page, "\nCombat Report ---------------------------------------------------------------------\n")
		_, _ = p.Fprintf(page, "  No activity.\n")
	}

	page.WriteString(`</pre></code></body>`)

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(page.String()))
}
