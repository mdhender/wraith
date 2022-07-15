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

func Parse(b []byte) ([]*Order, error) {
	var orders []*Order

	for z := tokens.FromBytes(b); !z.IsEof(); {
		if verb := accept(z, tokens.EOL); verb != nil {
			continue
		} else if verb = accept(z, tokens.Assemble); verb != nil {
			cmd := &Order{Line: verb.Line, Verb: verb}
			cmd.expectAssemble(z)
			orders = append(orders, cmd)
			continue
		} else if verb = accept(z, tokens.Control); verb != nil {
			cmd := &Order{Line: verb.Line, Verb: verb}
			cmd.expectCorSId(z)
			orders = append(orders, cmd)
			continue
		} else if verb = accept(z, tokens.Name); verb != nil {
			cmd := &Order{Line: verb.Line, Verb: verb}
			cmd.expectName(z)
			orders = append(orders, cmd)
			continue
		}

		// unknown order. reject the entire line.
		verb := z.Next()
		cmd := &Order{Line: verb.Line, Verb: verb}
		cmd.Errors = append(cmd.Errors, fmt.Errorf("unknown order %q", string(verb.Text)))
		cmd.reject(z)
		orders = append(orders, cmd)
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
