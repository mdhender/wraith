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

////////////////////////////////////////////////////////////////////////////////
// Context
//   USER - A USER is the person logging in to the system to use it. Currently,
//          the user is uniquely identified by their e-mail address. The user
//          also has a unique identifier assigned so that they can change their
//          e-mail address if needed.
////////////////////////////////////////////////////////////////////////////////

// IdentityService provides methods to create, update, and authenticate users.
type IdentityService interface {
	// AuthenticateUser verifies the email and secret against the stored values.
	// It returns the user data if successful.
	AuthenticateUser(AuthenticateUserRequest) AuthenticatedUserResponse
	// CreateUser creates a new user and returns the user data if successful.
	CreateUser(CreateUserRequest) UserResponse
	// DeleteUser deletes an existing user and returns the old user data if successful.
	DeleteUser(DeleteUserRequest) UserResponse
	// FetchUser retrieves an existing user and returns the user data if successful.
	FetchUser(FetchUserRequest) UserResponse
	// UpdateUser updates an existing user and returns the updated user data if successful.
	UpdateUser(UpdateUserRequest) UserResponse
}

// AuthenticateUserRequest is the request object for IdentityService.Authenticate.
type AuthenticateUserRequest struct {
	// Email is the e-mail address the user registered with.
	// Required.
	// example: "fred.flintrock@example.com"
	Email string
	// Secret is the hex-encoded passphrase used to authenticate the request.
	// example: "6261644d6f6f7365"
	Secret string
}

// AuthenticatedUserResponse is the response object containing the user's bearer token if authenticated.
type AuthenticatedUserResponse struct {
	// Id is the unique identifier for the user.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Token is the signed JWT if the request is authorized.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string
}

// CreateUserRequest is the request object for IdentityService.CreateUser.
type CreateUserRequest struct {
	// Email is the e-mail address the user registered with.
	// Required.
	// example: "fred.flintrock@example.com"
	Email string
	// Secret is the hex-encoded passphrase used to authenticate the user on future requests.
	// Required.
	// example: "6261644d6f6f7365"
	Secret string
}

// DeleteUserRequest is the request object for IdentityService.DeleteUser.
type DeleteUserRequest struct {
	// Id is the identifier for the user to delete.
	// Required.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address the user registered with.
	// Optional - include as a second consistency check.
	// example: "fred.flintrock@example.com"
	Email string
}

// FetchUserRequest is the request object for IdentityService.FetchUser.
// The caller may supply id or e-mail (or both).
type FetchUserRequest struct {
	// Id is the identifier for the user to retrieve.
	// Optional.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address for the user to retrieve.
	// Optional.
	// example: "fred.flintrock@example.com"
	Email string
}

// UpdateUserRequest is the request object for IdentityService.UpdateUser.
type UpdateUserRequest struct {
	// Id is the identifier for the user to update.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address the user registered with.
	// Optional - omit if not updating.
	// example: "fred.flintrock@example.com"
	Email string
	// Secret is the hex-encoded passphrase used to authenticate the user on future requests.
	// Optional - omit if not updating.
	// example: "6261644d6f6f7365"
	Secret string
}

// UserResponse is the response object containing the user's data if authenticated.
type UserResponse struct {
	// Id is the unique identifier for the user.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address the user registered with.
	// example: "fred.flintrock@example.com"
	Email string
}
