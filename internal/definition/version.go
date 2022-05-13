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

package definition

// VersionService makes nice greetings.
type VersionService interface {
	// Version returns the version of the engine.
	Version(VersionRequest) VersionResponse
}

// VersionRequest is the request object for VersionService.Version.
type VersionRequest struct {
	// Semver is a boolean.
	// example: true
	Semver bool
}

// VersionResponse is the response object containing a
// person's greeting.
type VersionResponse struct {
	// Version is the semantic version of the engine.
	// example: "0.1.0"
	Version string
}
