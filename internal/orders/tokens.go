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
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	Line       int
	Word       string
	Kind       string
	Integer    int
	number     float64
	percentage float64
	Id         string
	Unit       string
	Text       string
}

type tokenizer struct {
	line, offset int
	buffer       []byte
	pb           []*Token
}

func newTokenizer(b []byte) *tokenizer {
	return &tokenizer{line: 1, buffer: b}
}

const eof = rune(-1)

// eof returns true if we are at end of input and the pushback buffer is empty
func (t *tokenizer) eof() bool {
	return len(t.buffer) <= t.offset && len(t.pb) == 0
}

// unget adds the token to the pushback buffer
func (t *tokenizer) unget(tok *Token) {
	if tok.Kind == "eof" {
		return
	}
	t.pb = append(t.pb, tok)
}

func (t *tokenizer) next() *Token {
	if len(t.pb) != 0 {
		tok := t.pb[len(t.pb)-1]
		t.pb = t.pb[:len(t.pb)-1]
		return tok
	}

	var r rune
	var w int

	for !t.eof() {
		r, w = utf8.DecodeRune(t.buffer[t.offset:])
		t.offset += w

		if r == '\n' {
			t.line++
			return &Token{Line: t.line - 1, Kind: "eol"}
		} else if r == ';' {
			for !t.eof() {
				if r, w = utf8.DecodeRune(t.buffer[t.offset:]); r == '\n' {
					break
				}
				t.offset += w
			}
		} else if unicode.IsControl(r) {
			continue
		} else if !unicode.IsSpace(r) {
			break
		}
	}

	if t.eof() {
		return &Token{Line: t.line, Kind: "eof"}
	}

	if r == '{' {
		return &Token{Line: t.line, Kind: "block-open", Word: "{"}
	} else if r == '}' {
		return &Token{Line: t.line, Kind: "block-close", Word: "{"}
	}

	var word []byte
	if r == '"' {
		word = append(word, '"')
		for !t.eof() {
			if r, w = utf8.DecodeRune(t.buffer[t.offset:]); r == '\n' {
				break
			}
			t.offset += w
			if r == '\t' {
				word = append(word, ' ')
				continue
			} else if r == '"' {
				word = append(word, '"')
				break
			} else if unicode.IsControl(r) {
				continue
			}
			word = append(word, t.buffer[t.offset-w:t.offset]...)
		}
		return &Token{Line: t.line, Kind: "text", Word: string(word), Text: string(word)}
	}

	word = append(word, t.buffer[t.offset-w:t.offset]...)
	for !t.eof() {
		if r, w = utf8.DecodeRune(t.buffer[t.offset:]); r == '\n' {
			break
		} else if !(r == '-' || r == ',' || r == '.' || r == '%' || unicode.IsLetter(r) || unicode.IsDigit(r)) {
			break
		}
		t.offset += w
		word = append(word, t.buffer[t.offset-w:t.offset]...)
	}

	if isId(word) {
		if word[0] == 'c' || word[0] == 'C' {
			return &Token{Line: t.line, Kind: "colony-id", Word: string(word), Id: strings.ToUpper(string(word))}
		} else if word[0] == 'd' || word[0] == 'D' {
			return &Token{Line: t.line, Kind: "deposit-id", Word: string(word), Id: strings.ToUpper(string(word))}
		} else if word[0] == 's' || word[0] == 'S' {
			return &Token{Line: t.line, Kind: "ship-id", Word: string(word), Id: strings.ToUpper(string(word))}
		}
		return &Token{Line: t.line, Kind: "id", Word: string(word), Id: strings.ToUpper(string(word))}
	}

	if i, ok := tonum(word); ok {
		return &Token{Line: t.line, Kind: "integer", Word: string(word), Integer: i}
	}

	s := strings.ToLower(string(word))
	switch s {
	case "assemble":
		return &Token{Line: t.line, Word: string(word), Kind: "assemble", Text: s}
	case "control":
		return &Token{Line: t.line, Word: string(word), Kind: "control", Text: s}
	case "name":
		return &Token{Line: t.line, Word: string(word), Kind: "name", Text: s}
	case "automation-1":
		return &Token{Line: t.line, Word: string(word), Kind: "automation", Unit: "automation-1"}
	case "consumer-goods":
		return &Token{Line: t.line, Word: string(word), Kind: "consumer-goods", Unit: "consumer-goods"}
	case "factory-1":
		return &Token{Line: t.line, Word: string(word), Kind: "factory", Unit: "factory-1"}
	case "farm-1":
		return &Token{Line: t.line, Word: string(word), Kind: "farm", Unit: "farm-1"}
	case "hyper-drive-1":
		return &Token{Line: t.line, Word: string(word), Kind: "hyper-drive", Unit: "hyper-drive-1"}
	case "life-support-1":
		return &Token{Line: t.line, Word: string(word), Kind: "life-support", Unit: "life-support-1"}
	case "mine-1":
		return &Token{Line: t.line, Word: string(word), Kind: "mine", Unit: "mine-1"}
	case "research":
		return &Token{Line: t.line, Word: string(word), Kind: "research", Unit: "research"}
	case "sensor-1":
		return &Token{Line: t.line, Word: string(word), Kind: "sensor", Unit: "sensor-1"}
	case "space-drive-1":
		return &Token{Line: t.line, Word: string(word), Kind: "space-drive", Unit: "space-drive-1"}
	case "structural":
		return &Token{Line: t.line, Word: string(word), Kind: "structural", Unit: "structural"}
	case "transport-1":
		return &Token{Line: t.line, Word: string(word), Kind: "transport", Unit: "transport-1"}
	default:
		return &Token{Line: t.line, Word: string(word), Kind: "unknown"}
	}
}

func tonum(b []byte) (i int, ok bool) {
	if len(b) == 0 {
		return 0, false
	}
	for len(b) != 0 {
		switch b[0] {
		case '0':
			i = i * 10
		case '1':
			i = i*10 + 1
		case '2':
			i = i*10 + 2
		case '3':
			i = i*10 + 3
		case '4':
			i = i*10 + 4
		case '5':
			i = i*10 + 5
		case '6':
			i = i*10 + 6
		case '7':
			i = i*10 + 7
		case '8':
			i = i*10 + 8
		case '9':
			i = i*10 + 9
		case ',':
		default:
			return 0, false
		}
		b = b[1:]
	}
	return i, true
}
