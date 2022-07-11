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

// isId returns true if the word could be a valid id
func isId(b []byte) bool {
	// all ids are at least two characters
	if len(b) < 2 {
		return false
	}

	// colony id is C##
	if (b[0] == 'c' || b[0] == 'C') && ('0' < b[1] && b[1] <= '9') {
		return true
	}
	// ship id is S##
	if (b[0] == 's' || b[0] == 'S') && ('0' < b[1] && b[1] <= '9') {
		return true
	}

	// remaining ids are all 2 or more characters
	if len(b) < 3 {
		return false
	}

	// deposit id is DP##
	if (b[0] == 'd' || b[0] == 'D') && (b[1] == 'p' || b[1] == 'P') && ('0' < b[2] && b[2] <= '9') {
		return true
	}

	// not a known id
	return false
}

func toInteger(b []byte) (i int, ok bool) {
	if len(b) == 0 {
		return 0, false
	}
	for len(b) != 0 {
		switch b[0] {
		case ',', '_':
			// ignore commas and underscores
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
		default:
			return 0, false
		}
		b = b[1:]
	}
	return i, true
}
