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
	"github.com/mdhender/wraith/engine"
	"github.com/mdhender/wraith/internal/tokens"
)

type Arg struct {
	Integer    int
	Number     float64
	Percentage float64
	QuotedText string
	Text       string
}

type Order struct {
	Line   int
	Verb   *tokens.Token
	Args   []*tokens.Token
	Reject []*tokens.Token // nil unless there was an error parsing
	Errors []error         // nil unless there was an error parsing
}

func (o *Order) String() string {
	var s string
	if o == nil {
		return s
	}
	if o.Verb != nil {
		s += string(o.Verb.Text)
	}
	for _, t := range o.Args {
		if t != nil {
			s += " " + string(t.Text)
		}
	}
	for _, t := range o.Reject {
		if t != nil {
			s += " " + string(t.Text)
		}
	}
	for i, e := range o.Errors {
		if i != 0 {
			s += "\n"
		}
		s += fmt.Sprintf("  ;; %v", e)
	}
	return s
}

type IOrder interface {
	Line() int
	String() string
	DryRun() bool
	Execute(e *engine.Engine) error
}

type Error struct {
	Error  error
	Tokens []*tokens.Token
}

func (o *Error) String() string {
	var s string
	if o == nil {
		return s
	}
	for _, t := range o.Tokens {
		s += " " + string(t.Text)
	}
	if o.Error != nil {
		s += fmt.Sprintf(" ;; %v", o.Error)
	}
	return s
}

type AssembleGroup struct {
	Verb  *tokens.Token
	Id    *tokens.Token // id of colony or ship to assemble the group in
	Qty   *tokens.Token // number of units to add to group
	Error *Error        // nil unless there was an error parsing
}

func (o *AssembleGroup) String() string {
	var s string
	if o == nil {
		return s
	}
	for i, t := range []*tokens.Token{o.Verb, o.Id, o.Qty} {
		if t == nil {
			break
		} else if i != 0 {
			s += " "
		}
		s += string(t.Text)
	}
	return s + o.Error.String()
}

type AssembleFactoryGroup struct {
	Verb    *tokens.Token
	Id      *tokens.Token // id of colony or ship to assemble the group in
	Qty     *tokens.Token // number of factory units
	Factory *tokens.Token // type of factory unit
	Product *tokens.Token // type of unit to produce
	Error   *Error        // nil unless there was an error parsing
}

func (o *AssembleFactoryGroup) String() string {
	var s string
	if o == nil {
		return s
	}
	for i, t := range []*tokens.Token{o.Verb, o.Id, o.Qty, o.Factory, o.Product} {
		if t == nil {
			break
		} else if i != 0 {
			s += " "
		}
		s += string(t.Text)
	}
	return s + o.Error.String()
}

type AssembleMineGroup struct {
	Verb      *tokens.Token
	Id        *tokens.Token // id of colony or ship to assemble the group in
	Qty       *tokens.Token // number of mine units
	Mine      *tokens.Token // type of mine unit
	DepositId *tokens.Token // id of deposit to mine
	Error     *Error        // nil unless there was an error parsing
}

func (o *AssembleMineGroup) String() string {
	var s string
	if o == nil {
		return s
	}
	for i, t := range []*tokens.Token{o.Verb, o.Id, o.Qty, o.Mine, o.DepositId} {
		if t == nil {
			break
		} else if i != 0 {
			s += " "
		}
		s += string(t.Text)
	}
	return s + o.Error.String()
}

type Name struct {
	Verb  *tokens.Token
	Id    *tokens.Token // id of object to name
	Name  *tokens.Token // new name of object
	Error *Error        // nil unless there was an error parsing
}

func (o *Name) String() string {
	var s string
	if o == nil {
		return s
	}
	for i, t := range []*tokens.Token{o.Verb, o.Id, o.Name} {
		if t == nil {
			break
		} else if i != 0 {
			s += " "
		}
		s += string(t.Text)
	}
	return s + o.Error.String()
}

type Unknown struct {
	Verb  *tokens.Token
	Error *Error // nil unless there was an error parsing
}

func (o *Unknown) String() string {
	var s string
	if o == nil {
		return s
	}
	for i, t := range []*tokens.Token{o.Verb} {
		if t == nil {
			break
		} else if i != 0 {
			s += " "
		}
		s += string(t.Text)
	}
	return s + o.Error.String()
}

func Parse(b []byte) ([]*Order, error) {
	var orders []*Order

	for z := tokens.FromBytes(b); !z.IsEof(); {
		if accept(z, tokens.EOL) != nil {
			continue
		}

		if verb := accept(z, tokens.Assemble); verb != nil {
			cmd := &Order{Line: verb.Line, Verb: verb}
			cmd.expectAssemble(z)
			orders = append(orders, cmd)
		} else if verb = accept(z, tokens.Name); verb != nil {
			cmd := &Order{Line: verb.Line, Verb: verb}
			cmd.expectName(z)
			orders = append(orders, cmd)
		} else {
			cmd := &Order{Line: verb.Line, Verb: verb}
			cmd.Errors = append(cmd.Errors, fmt.Errorf("unknown order"))
			cmd.reject(z)
			orders = append(orders, cmd)
		}
	}

	return orders, nil
}

func accept(z *tokens.Tokenizer, kinds ...tokens.Kind) *tokens.Token {
	tok := z.Next()
	for _, kind := range kinds {
		if tok.Kind == kind {
			return tok
		}
	}
	z.UnGet(tok)
	return nil
}

func (cmd *Order) expectAssemble(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.ColonyId, tokens.ShipId); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: expected ship or colony id", cmd.Line))
		cmd.reject(z)
		return false
	}
	cmd.Args = append(cmd.Args, t)
	if t = accept(z, tokens.Integer); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: expected quantity", cmd.Line))
		cmd.reject(z)
		return false
	}
	cmd.Args = append(cmd.Args, t)
	if t = accept(z, tokens.FactoryUnit); t != nil {
		cmd.Args = append(cmd.Args, t)
		return cmd.expectFactoryGroup(z)
	}
	if t = accept(z, tokens.MineUnit); t != nil {
		cmd.Args = append(cmd.Args, t)
		return cmd.expectMineGroup(z)
	}
	cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: unexpected input on assemble group order", cmd.Line))
	return true
}

func (cmd *Order) expectFactoryGroup(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.AutomationUnit, tokens.ConsumerGoodsUnit, tokens.FactoryUnit, tokens.FarmUnit, tokens.HyperDriveUnit, tokens.MineUnit, tokens.ResearchUnit, tokens.SensorUnit, tokens.SpaceDriveUnit, tokens.StructuralUnit, tokens.TransportUnit); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: expected unit to produce", cmd.Line))
		cmd.reject(z)
		return false
	}
	cmd.Args = append(cmd.Args, t)
	if t = accept(z, tokens.EOL, tokens.EOF); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: unexpected input on assemble factory group order", cmd.Line))
		cmd.reject(z)
		return false
	}
	return true
}

func (cmd *Order) expectMineGroup(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.DepositId); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: expected deposit to mine", cmd.Line))
		cmd.reject(z)
		return false
	}
	cmd.Args = append(cmd.Args, t)
	if t = accept(z, tokens.EOL, tokens.EOF); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: unexpected input on assemble mine group order", cmd.Line))
		cmd.reject(z)
		return false
	}
	return true
}

func (cmd *Order) expectName(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.ColonyId, tokens.ShipId); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: expected ship or colony id", cmd.Line))
		cmd.reject(z)
		return false
	}
	cmd.Args = append(cmd.Args, t)
	if t = accept(z, tokens.QuotedText); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: expected name", cmd.Line))
		cmd.reject(z)
		return false
	}
	cmd.Args = append(cmd.Args, t)
	if t = accept(z, tokens.EOL, tokens.EOF); t == nil {
		cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: unexpected input on name order", cmd.Line))
		cmd.reject(z)
		return false
	}
	return true
}

// consume until we find EOL or EOF token.
// the slice of tokens returned will not include EOL or EOF
func (cmd *Order) reject(z *tokens.Tokenizer) {
	for t := z.Next(); t.Kind != tokens.EOF; t = z.Next() {
		if t.Kind == tokens.EOL {
			z.UnGet(t)
			break
		}
		cmd.Reject = append(cmd.Reject, t)
	}
}
