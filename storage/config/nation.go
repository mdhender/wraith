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

package config

// Nation configuration
type Nation struct {
	Id          string
	Name        string `json:"name"`
	Description string `json:"description"`
	Speciality  string `json:"speciality"`
	Government  struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
	} `json:"government"`
	Homeworld struct {
		Name     string `json:"name"`
		Location struct {
			X     int `json:"x"`
			Y     int `json:"y"`
			Z     int `json:"z"`
			Orbit int `json:"orbit"`
		} `json:"location"`
	} `json:"homeworld"`
}
