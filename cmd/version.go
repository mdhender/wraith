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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var globalVersion struct {
	Version    string
	major      int
	minor      int
	patch      int
	preRelease string
	build      string
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "show application version",
	Long:  `Show application version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", globalVersion.Version)
	},
}

func init() {
	globalVersion.major, globalVersion.minor, globalVersion.patch = 0, 1, 0
	globalVersion.preRelease = ""
	globalVersion.build = ""

	// format the version per https://semver.org/ rules
	globalVersion.Version = fmt.Sprintf("%d.%d.%d", globalVersion.major, globalVersion.minor, globalVersion.patch)
	if globalVersion.preRelease != "" {
		globalVersion.Version = globalVersion.Version + "-" + globalVersion.preRelease
	}
	if globalVersion.build != "" {
		globalVersion.Version = globalVersion.Version + "+" + globalVersion.build
	}

	cmdBase.AddCommand(cmdVersion)
}
