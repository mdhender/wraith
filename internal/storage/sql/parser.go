/*
 * Wraith Game Engine
 * Copyright (c) 2022 Michael D. Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package sql

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"unicode"
	"unicode/utf8"
)

func Parse(b []byte) error {
	for len(b) > 0 {
		token, buffer := next(b)
		if len(token) != 0 {
			fmt.Printf("token(%q)\n", string(token))
		}
		b = buffer
	}
	return errors.New("not implemented")
}

func next(b []byte) (token, buffer []byte) {
	delimiters := []byte{'.', ',', '(', ')', ';', '=', '+', '-', '*', '/', '\n'}

	if len(b) == 0 || b[0] == ';' {
		return nil, nil
	}

	// always treat some characters specially
	if bytes.IndexByte([]byte(delimiters), b[0]) != -1 {
		return b[:1], b[1:]
	}

	// token will either be a run of spaces/invalid runes or a run of runes to a space or invalid rune.
	r, w := utf8.DecodeRune(b)
	length := w
	if r == utf8.RuneError || unicode.IsSpace(r) {
		// consume runes up to the next non-space rune
		for len(b) < length && b[length] != '\n' {
			r, w = utf8.DecodeRune(b[length:])
			if !(r == utf8.RuneError || unicode.IsSpace(r)) {
				break
			}
			length += w
		}
	} else {
		// return a run of runes up to the next space, special character, or invalid rune
		for len(b) < length {
			if bytes.IndexByte(delimiters, b[length]) != -1 {
				return b[:1], b[1:]
			}
			r, w = utf8.DecodeRune(b)
			if r == utf8.RuneError || unicode.IsSpace(r) {
				break
			}
			length += w
		}
	}
	return b[:length], b[length:]
}
