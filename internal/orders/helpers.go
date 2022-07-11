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

// isId returns true if the word could be a valid id
func isId(word []byte) bool {
	// all ids are at least two characters
	if len(word) < 2 {
		return false
	}

	// colony id is C##
	if (word[0] == 'c' || word[0] == 'C') && ('0' < word[1] && word[1] <= '9') {
		return true
	}
	// ship id is S##
	if (word[0] == 's' || word[0] == 'S') && ('0' < word[1] && word[1] <= '9') {
		return true
	}

	// remaining ids are all 2 or more characters
	if len(word) < 3 {
		return false
	}

	// deposit id is DP##
	if (word[0] == 'd' || word[0] == 'D') && (word[1] == 'p' || word[1] == 'P') && ('0' < word[2] && word[2] <= '9') {
		return true
	}

	// not a known id
	return false
}
