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

// Token is a token from the input buffer
type Token struct {
	Kind    Kind
	Line    int     // line number in the input
	Integer int     // populated only for Integers
	Number  float64 // populated for both number and percentage
	Text    []byte  // always populated
}

func (t *Token) String() string {
	if t == nil {
		return ""
	}
	return string(t.Text)
}

type Kind int

const (
	Unknown Kind = iota
	EOF
	EOL
	BlockOpen
	BlockClose

	// atoms, so to speak

	Integer
	Number
	Percentage
	QuotedText
	Text

	// identifiers

	ColonyId
	DepositId
	ShipId

	// order verbs

	Assemble
	Name

	// units

	AutomationUnit
	ConsumerGoodsUnit
	FactoryUnit
	FarmUnit
	HyperDriveUnit
	MineUnit
	ResearchUnit
	SensorUnit
	SpaceDriveUnit
	StructuralUnit
	TransportUnit
)
