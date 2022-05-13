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

package greeter

import (
	"context"
	"github.com/pkg/errors"
	"strings"
)

type Service struct{}

func (s Service) Greet(ctx context.Context, request GreetRequest) (*GreetResponse, error) {
	// consider rejecting unknown fields (https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body)
	if request.Name == "" {
		return nil, errors.New("missing 'name'")
	}
	return &GreetResponse{
		Greeting: strings.ToUpper(request.Name),
	}, nil
}
