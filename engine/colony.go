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

type XColony struct {
	Id       int
	Kind     string
	Location struct {
		X     int
		Y     int
		Z     int
		Star  int
		Orbit int
	}
	TechLevel  int
	Population struct {
		Professional      ReportPopulation
		Soldier           ReportPopulation
		Unskilled         ReportPopulation
		Unemployed        ReportPopulation
		ConstructionCrews int
		SpyTeams          int
		Births            int
		Deaths            int
	}
	Inventory     []*Inventory
	FactoryGroups []*XGroup
	FarmGroups    []*XGroup
	MiningGroups  []*XGroup
}

type XPopulation struct {
	Code   string
	Qty    int
	Pay    float64
	Ration float64
}

type XGroup struct {
	Id    int
	Name  string
	Units []*XGroupUnits
}

type XGroupUnits struct {
	TechLevel int
	Qty       int
	Stages    []int
}
