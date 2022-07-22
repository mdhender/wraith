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

func (o *Order) expectAssemble(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.ColonyId, tokens.ShipId); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: expected ship or colony id", o.Line))
		o.reject(z)
		return false
	}
	o.Args = append(o.Args, t)
	if t = accept(z, tokens.Integer); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: expected quantity", o.Line))
		o.reject(z)
		return false
	}
	o.Args = append(o.Args, t)
	if t = accept(z, tokens.FactoryUnit); t != nil {
		o.Verb.Kind = tokens.AssembleFactoryGroup
		o.Args = append(o.Args, t)
		return o.expectFactoryGroup(z)
	}
	if t = accept(z, tokens.FarmUnit); t != nil {
		o.Verb.Kind = tokens.AssembleFarmGroup
		o.Args = append(o.Args, t)
		return o.expectFarmGroup(z)
	}
	if t = accept(z, tokens.MineUnit); t != nil {
		o.Verb.Kind = tokens.AssembleMineGroup
		o.Args = append(o.Args, t)
		return o.expectMineGroup(z)
	}
	o.Errors = append(o.Errors, fmt.Errorf("%d: unexpected input on assemble group order", o.Line))
	return true
}

func (o *Order) expectCorSId(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.ColonyId, tokens.ShipId); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: expected ship or colony id", o.Line))
		o.reject(z)
		return false
	}
	o.Args = append(o.Args, t)
	if t = accept(z, tokens.EOL, tokens.EOF); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: unexpected input following ship or colony id", o.Line))
		o.reject(z)
		return false
	}
	return true
}

func (o *Order) expectFactoryGroup(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z,
		tokens.AutomationUnit, tokens.ConsumerGoodsUnit,
		tokens.FactoryUnit, tokens.FarmUnit,
		tokens.HyperDriveUnit, tokens.LifeSupportUnit,
		tokens.MineUnit, tokens.ResearchUnit,
		tokens.SensorUnit, tokens.SpaceDriveUnit,
		tokens.StructuralUnit, tokens.TransportUnit); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: expected unit to produce", o.Line))
		o.reject(z)
		return false
	}
	o.Args = append(o.Args, t)
	if t = accept(z, tokens.EOL, tokens.EOF); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: unexpected input on assemble factory group order", o.Line))
		o.reject(z)
		return false
	}
	return true
}

func (o *Order) expectFarmGroup(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.FoodUnit); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: expected unit to produce", o.Line))
		o.reject(z)
		return false
	}
	o.Args = append(o.Args, t)
	if t = accept(z, tokens.EOL, tokens.EOF); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: unexpected input on assemble farm group order", o.Line))
		o.reject(z)
		return false
	}
	return true
}

func (o *Order) expectMineGroup(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.DepositId); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: expected deposit to mine", o.Line))
		o.reject(z)
		return false
	}
	o.Args = append(o.Args, t)
	if t = accept(z, tokens.EOL, tokens.EOF); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: unexpected input on assemble mine group order", o.Line))
		o.reject(z)
		return false
	}
	return true
}

func (o *Order) expectName(z *tokens.Tokenizer) bool {
	var t *tokens.Token
	if t = accept(z, tokens.ColonyId, tokens.ShipId); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: expected ship or colony id", o.Line))
		o.reject(z)
		return false
	}
	o.Args = append(o.Args, t)
	if t = accept(z, tokens.QuotedText); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: expected name", o.Line))
		o.reject(z)
		return false
	}
	o.Args = append(o.Args, t)
	if t = accept(z, tokens.EOL, tokens.EOF); t == nil {
		o.Errors = append(o.Errors, fmt.Errorf("%d: unexpected input on name order", o.Line))
		o.reject(z)
		return false
	}
	return true
}

// consume until we find EOL or EOF token.
// the slice of tokens returned will not include EOL or EOF
func (o *Order) reject(z *tokens.Tokenizer) {
	for t := z.Next(); t.Kind != tokens.EOF; t = z.Next() {
		if t.Kind == tokens.EOL {
			z.UnGet(t)
			break
		}
		o.Reject = append(o.Reject, t)
	}
}
