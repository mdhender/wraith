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

package identity

import (
	"context"
	"github.com/pkg/errors"
)

type Service struct{}

func (s Service) AuthenticateUser(ctx context.Context, request AuthenticateUserRequest) (*AuthenticatedUserResponse, error) {
	if request.Email == "fred.flintrock@example.com" && request.Secret == "$2a$12$7hRCqIDEQ05MD7nmzFgfxuTl/cwZkDF8.Hb3s1bpC78ey.dbUdXCW" {
		return &AuthenticatedUserResponse{
			Id:    "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d",
			Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		}, nil
	}
	return nil, errors.New("not authorized")
}

func (s Service) CreateUser(ctx context.Context, request CreateUserRequest) (*UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) DeleteUser(ctx context.Context, request DeleteUserRequest) (*UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) FetchUser(ctx context.Context, request FetchUserRequest) (*UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) UpdateUser(ctx context.Context, request UpdateUserRequest) (*UserResponse, error) {
	//TODO implement me
	panic("implement me")
}
