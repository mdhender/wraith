////////////////////////////////////////////////////////////////////////////////
// wraith - the wraith game engine and Server
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

package cheese

type ClusterList []*ClusterListItem

type ClusterListItem struct {
	X, Y, Z  int
	QtyStars int
	Distance float64
}

// Len implements the sort.Sort interface
func (u ClusterList) Len() int {
	return len(u)
}

// Less implements the sort.Sort interface
func (u ClusterList) Less(i, j int) bool {
	if u[i].Distance < u[j].Distance {
		return true
	} else if u[i].Distance > u[j].Distance {
		return false
	}

	if u[i].X < u[j].X {
		return true
	} else if u[i].X > u[j].X {
		return false
	}

	if u[i].Y < u[j].Y {
		return true
	} else if u[i].Y > u[j].Y {
		return false
	}

	return u[i].Z < u[j].Z
}

// Swap implements the sort.Sort interface
func (u ClusterList) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
