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

package orders

import (
	"fmt"
	"log"
	"strings"
)

type Arg struct {
	Integer    int
	Number     float64
	Percentage float64
	QuotedText string
	Text       string
}

type Order struct {
	Line    int
	Command interface{}
}

type AssembleFactoryGroup struct {
	Id      string // id of colony or ship to assemble the group in
	Factory string // type of factory unit
	Qty     int    // number of factory units
	Product string // type of unit to produce
}

type AssembleMineGroup struct {
	Id        string // id of colony or ship to assemble the group in
	Mine      string // type of mine unit
	Qty       int    // number of mine units
	DepositId string // id of deposit to mine
}

type Error struct {
	Words []string
	Error error
}

type NameOrder struct {
	Id   string // id of object to name
	Name string // new name of object
}

func Parse(b []byte) ([]*Order, error) {
	var orders []*Order
	for tz := newTokenizer(b); !tz.eof(); {
		if accept(tz, "eol") != nil {
			continue
		}
		if verb := accept(tz, "assemble"); verb != nil {
			order, err := expectAssemble(tz, verb)
			if err != nil {
				log.Println(err)
			} else {
				orders = append(orders, order)
			}
		} else if verb = accept(tz, "name"); verb != nil {
			order, err := expectName(tz, verb)
			if err != nil {
				log.Println(err)
			} else {
				orders = append(orders, order)
			}
		} else {
			// unknown order
			log.Printf("unknown order\n")
			_ = reject(tz, nil)
		}
	}
	return orders, nil
}

func accept(tz *tokenizer, kinds ...string) *Token {
	tok := tz.next()
	for _, kind := range kinds {
		if tok.Kind == kind {
			return tok
		}
	}
	tz.unget(tok)
	return nil
}

func expectAssemble(tz *tokenizer, verb *Token) (*Order, error) {
	cors := accept(tz, "colony-id", "ship-id")
	if cors == nil {
		return nil, reject(tz, fmt.Errorf("%d: expected ship or colony id", verb.Line))
	}
	qty := accept(tz, "integer")
	if qty == nil {
		return nil, fmt.Errorf("%d: expected quantity", verb.Line)
	}
	if mine := accept(tz, "mine"); mine != nil {
		deposit := accept(tz, "deposit-id")
		if deposit == nil {
			return nil, reject(tz, fmt.Errorf("%d: expected deposit to mine", verb.Line))
		}
		nl := accept(tz, "eol", "eof")
		if nl == nil {
			return nil, reject(tz, fmt.Errorf("%d: unexpected input on order", verb.Line))
		}
		return &Order{Line: verb.Line, Command: &AssembleMineGroup{Id: cors.Id, Mine: mine.Text, Qty: qty.Integer, DepositId: deposit.Id}}, nil
	}
	if factory := accept(tz, "factory"); factory != nil {
		unit := accept(tz, "automation", "consumer-goods", "factory", "farm", "hyper-drive", "mine", "research", "structural", "sensor", "space-drive", "structural", "transport")
		if unit == nil {
			return nil, reject(tz, fmt.Errorf("%d: expected unit to produce", verb.Line))
		}
		nl := accept(tz, "eol", "eof")
		if nl == nil {
			return nil, reject(tz, fmt.Errorf("%d: unexpected input on order", verb.Line))
		}
		return &Order{Line: verb.Line, Command: &AssembleFactoryGroup{Id: cors.Id, Factory: factory.Text, Qty: qty.Integer, Product: unit.Unit}}, nil
	}
	order := &Order{Line: verb.Line, Command: &Error{Error: fmt.Errorf("%d: expected mine or factory unit", verb.Line)}}
	return order, reject(tz, fmt.Errorf("%d: expected mine or factory unit", verb.Line))
}

func expectName(tz *tokenizer, verb *Token) (*Order, error) {
	cors := accept(tz, "colony-id", "ship-id")
	if cors == nil {
		return nil, reject(tz, fmt.Errorf("%d: expected ship or colony id", verb.Line))
	}
	name := accept(tz, "text")
	if name == nil {
		return nil, reject(tz, fmt.Errorf("%d: expected name", verb.Line))
	} else if name.Kind != "text" || !strings.HasPrefix(name.Text, "\"") || !strings.HasSuffix(name.Text, "\"") {
		return nil, reject(tz, fmt.Errorf("%d: expected name to be quoted text", verb.Line))
	}
	nl := accept(tz, "eol", "eof")
	if nl == nil {
		return nil, reject(tz, fmt.Errorf("%d: unexpected input on order", verb.Line))
	}
	return &Order{Line: verb.Line, Command: &NameOrder{Id: cors.Id, Name: name.Text}}, nil
}

func reject(tz *tokenizer, err error) error {
	// consume until we find eol or eof token
	for tok := tz.next(); !(tok.Kind == "eof" || tok.Kind == "eol"); {
		tok = tz.next()
	}
	return err
}
