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

package models

import "errors"

var ErrDuplicateKey = errors.New("duplicate key")
var ErrInvalidField = errors.New("invalid field")
var ErrMissingField = errors.New("missing field")
var ErrNoConnection = errors.New("no connection")
var ErrNoDataFound = errors.New("no data found")
var ErrNotImplemented = errors.New("not implemented")
var ErrUnauthorized = errors.New("unauthorized")
