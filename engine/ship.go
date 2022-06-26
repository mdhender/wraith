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

package engine

type XShip struct {
	Id         int // key for database
	No         int
	Name       string
	Kind       string
	Location   *Planet
	TechLevel  int
	Population struct {
		Professional      XPopulation
		Soldier           XPopulation
		Unskilled         XPopulation
		Unemployed        XPopulation
		ConstructionCrews int
		SpyTeams          int
		RebelPct          float64
	}
	Hull          []*Inventory // units used to build ship
	Inventory     []*Inventory // units stored in ship
	FactoryGroups []*FactoryGroup
	FarmGroups    []*XGroup
}
