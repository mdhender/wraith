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
	Line    int
	Command interface{}
}

type Error struct {
	Error  error
	Tokens []*tokens.Token
}

type AssembleGroup struct {
	Verb  *tokens.Token
	Id    *tokens.Token // id of colony or ship to assemble the group in
	Qty   *tokens.Token // number of units to add to group
	Error *Error        // nil unless there was an error parsing
}

type AssembleFactoryGroup struct {
	Verb    *tokens.Token
	Id      *tokens.Token // id of colony or ship to assemble the group in
	Qty     *tokens.Token // number of factory units
	Factory *tokens.Token // type of factory unit
	Product *tokens.Token // type of unit to produce
	Error   *Error        // nil unless there was an error parsing
}

type AssembleMineGroup struct {
	Verb      *tokens.Token
	Id        *tokens.Token // id of colony or ship to assemble the group in
	Qty       *tokens.Token // number of mine units
	Mine      *tokens.Token // type of mine unit
	DepositId *tokens.Token // id of deposit to mine
	Error     *Error        // nil unless there was an error parsing
}

type Name struct {
	Verb  *tokens.Token
	Id    *tokens.Token // id of object to name
	Name  *tokens.Token // new name of object
	Error *Error        // nil unless there was an error parsing
}

type Unknown struct {
	Verb  *tokens.Token
	Error *Error // nil unless there was an error parsing
}

func Parse(b []byte) ([]interface{}, error) {
	var orders []interface{}

	for z := tokens.FromBytes(b); !z.IsEof(); {
		if accept(z, tokens.EOL) != nil {
			continue
		}

		if verb := accept(z, tokens.Assemble); verb != nil {
			orders = append(orders, expectAssemble(z, &AssembleGroup{Verb: verb}))
		} else if verb = accept(z, tokens.Name); verb != nil {
			orders = append(orders, expectName(z, &Name{Verb: verb}))
		} else {
			orders = append(orders, &Unknown{Verb: z.Next(), Error: &Error{Error: fmt.Errorf("unknown order"), Tokens: reject(z)}})
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

func expectAssemble(z *tokens.Tokenizer, cmd *AssembleGroup) interface{} {
	if cmd.Id = accept(z, tokens.ColonyId, tokens.ShipId); cmd.Id == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: expected ship or colony id", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	if cmd.Qty = accept(z, tokens.Integer); cmd.Qty == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: expected quantity", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	if mine := accept(z, tokens.MineUnit); mine != nil {
		return expectMineGroup(z, &AssembleMineGroup{
			Verb: cmd.Verb,
			Id:   cmd.Id,
			Qty:  cmd.Qty,
			Mine: mine,
		})
	}
	if factory := accept(z, tokens.FactoryUnit); factory != nil {
		return expectFactoryGroup(z, &AssembleFactoryGroup{
			Verb:    cmd.Verb,
			Id:      cmd.Id,
			Qty:     cmd.Qty,
			Factory: factory,
		})
	}
	cmd.Error = &Error{Error: fmt.Errorf("%d: unexpected input on assemble group order", cmd.Verb.Line), Tokens: reject(z)}
	return cmd
}

func expectFactoryGroup(z *tokens.Tokenizer, cmd *AssembleFactoryGroup) *AssembleFactoryGroup {
	if cmd.Product = accept(z, tokens.AutomationUnit, tokens.ConsumerGoodsUnit, tokens.FactoryUnit, tokens.FarmUnit, tokens.HyperDriveUnit, tokens.MineUnit, tokens.ResearchUnit, tokens.SensorUnit, tokens.SpaceDriveUnit, tokens.StructuralUnit, tokens.TransportUnit); cmd.Product == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: expected unit to produce", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	if nl := accept(z, tokens.EOL, tokens.EOF); nl == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: unexpected input on assemble factory group order", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	return cmd
}

func expectMineGroup(z *tokens.Tokenizer, cmd *AssembleMineGroup) *AssembleMineGroup {
	if cmd.DepositId = accept(z, tokens.DepositId); cmd.DepositId == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: expected deposit to mine", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	if nl := accept(z, tokens.EOL, tokens.EOF); nl == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: unexpected input on assemble mine group order", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	return cmd
}

func expectName(z *tokens.Tokenizer, cmd *Name) *Name {
	if cmd.Id = accept(z, tokens.ColonyId, tokens.ShipId); cmd.Id == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: expected ship or colony id", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	if cmd.Name = accept(z, tokens.QuotedText); cmd.Name == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: expected name", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	if nl := accept(z, tokens.EOL, tokens.EOF); nl == nil {
		cmd.Error = &Error{Error: fmt.Errorf("%d: unexpected input on name order", cmd.Verb.Line), Tokens: reject(z)}
		return cmd
	}
	return cmd
}

// consume until we find EOL or EOF token.
// the slice of tokens returned will not include EOL or EOF
func reject(z *tokens.Tokenizer) []*tokens.Token {
	var toks []*tokens.Token
	for t := z.Next(); t.Kind != tokens.EOF; t = z.Next() {
		if t.Kind == tokens.EOL {
			z.UnGet(t)
			break
		}
		toks = append(toks, t)
	}
	return toks
}
