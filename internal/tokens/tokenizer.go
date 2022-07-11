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

package tokens

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Tokenizer struct {
	line, offset int
	buffer       []byte
	pb           []*Token
}

func FromBytes(b []byte) *Tokenizer {
	return &Tokenizer{line: 1, buffer: b}
}

func FromString(s string) *Tokenizer {
	return FromBytes([]byte(s))
}

// IsEof returns true if we are at end of input and the pushback buffer is empty
func (z *Tokenizer) IsEof() bool {
	return len(z.buffer) <= z.offset && len(z.pb) == 0
}

// Next returns the next available token.
// If tokens have been pushed back, they will be returned first.
// Returns EOF continuously at end of input.
func (z *Tokenizer) Next() *Token {
	if len(z.pb) != 0 {
		tok := z.pb[len(z.pb)-1]
		z.pb = z.pb[:len(z.pb)-1]
		return tok
	}

	var r rune
	var w int

	for !z.IsEof() {
		r, w = utf8.DecodeRune(z.buffer[z.offset:])
		z.offset += w

		if r == '\n' {
			z.line++
			return &Token{Line: z.line - 1, Kind: EOL}
		} else if r == ';' {
			for !z.IsEof() {
				if r, w = utf8.DecodeRune(z.buffer[z.offset:]); r == '\n' {
					break
				}
				z.offset += w
			}
		} else if unicode.IsControl(r) {
			continue
		} else if !unicode.IsSpace(r) {
			break
		}
	}

	if z.IsEof() {
		return &Token{Line: z.line, Kind: EOF}
	}

	if r == '{' {
		return &Token{Line: z.line, Kind: BlockOpen, Text: z.buffer[z.offset-w : z.offset]}
	} else if r == '}' {
		return &Token{Line: z.line, Kind: BlockClose, Text: z.buffer[z.offset-w : z.offset]}
	}

	var word []byte
	if r == '"' {
		word = append(word, '"')
		for !z.IsEof() {
			if r, w = utf8.DecodeRune(z.buffer[z.offset:]); r == '\n' {
				break
			}
			z.offset += w
			if r == '\t' {
				word = append(word, ' ')
				continue
			} else if r == '"' {
				word = append(word, '"')
				break
			} else if unicode.IsControl(r) {
				continue
			}
			word = append(word, z.buffer[z.offset-w:z.offset]...)
		}
		return &Token{Line: z.line, Kind: QuotedText, Text: word}
	}

	word = append(word, z.buffer[z.offset-w:z.offset]...)
	for !z.IsEof() {
		if r, w = utf8.DecodeRune(z.buffer[z.offset:]); r == '\n' {
			break
		} else if !(r == '-' || r == ',' || r == '.' || r == '%' || unicode.IsLetter(r) || unicode.IsDigit(r)) {
			break
		}
		z.offset += w
		word = append(word, z.buffer[z.offset-w:z.offset]...)
	}

	if isId(word) {
		word = bytes.ToUpper(word)
		switch word[0] {
		case 'C':
			return &Token{Line: z.line, Kind: ColonyId, Text: word}
		case 'D':
			return &Token{Line: z.line, Kind: DepositId, Text: word}
		case 'S':
			return &Token{Line: z.line, Kind: ShipId, Text: word}
		}
		panic(fmt.Sprintf("assert(word != %q)", string(word)))
	}

	if i, ok := toInteger(word); ok {
		return &Token{Line: z.line, Kind: Integer, Text: word, Integer: i}
	}

	if n, err := strconv.ParseFloat(string(word), 64); err == nil {
		return &Token{Line: z.line, Kind: Number, Text: word, Number: n}
	}

	if bytes.Equal(word, []byte("assemble")) {
		return &Token{Line: z.line, Kind: Assemble, Text: word}
	}
	if bytes.HasPrefix(word, []byte("automation-")) {
		return &Token{Line: z.line, Kind: AutomationUnit, Text: word}
	}
	if bytes.Equal(word, []byte("consumer-goods")) {
		return &Token{Line: z.line, Kind: ConsumerGoodsUnit, Text: word}
	}
	if bytes.HasPrefix(word, []byte("factory-")) {
		return &Token{Line: z.line, Kind: FactoryUnit, Text: word}
	}
	if bytes.HasPrefix(word, []byte("farm-")) {
		return &Token{Line: z.line, Kind: FarmUnit, Text: word}
	}
	if bytes.HasPrefix(word, []byte("hyper-drive-")) {
		return &Token{Line: z.line, Kind: HyperDriveUnit, Text: word}
	}
	if bytes.HasPrefix(word, []byte("mine-")) {
		return &Token{Line: z.line, Kind: MineUnit, Text: word}
	}
	if bytes.Equal(word, []byte("name")) {
		return &Token{Line: z.line, Kind: Name, Text: word}
	}
	if bytes.Equal(word, []byte("research")) {
		return &Token{Line: z.line, Kind: ResearchUnit, Text: word}
	}
	if bytes.HasPrefix(word, []byte("sensor-")) {
		return &Token{Line: z.line, Kind: SensorUnit, Text: word}
	}
	if bytes.HasPrefix(word, []byte("space-drive-")) {
		return &Token{Line: z.line, Kind: SpaceDriveUnit, Text: word}
	}
	if bytes.Equal(word, []byte("structural")) {
		return &Token{Line: z.line, Kind: StructuralUnit, Text: word}
	}
	if bytes.HasPrefix(word, []byte("transport-")) {
		return &Token{Line: z.line, Kind: TransportUnit, Text: word}
	}

	return &Token{Line: z.line, Kind: Text, Text: bytes.ToLower(word)}
}

// UnGet adds the token to the pushback buffer
func (z *Tokenizer) UnGet(t *Token) {
	if t.Kind == EOF {
		return
	}
	z.pb = append(z.pb, t)
}
