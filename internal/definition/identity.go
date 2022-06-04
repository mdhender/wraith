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
//   USER - A USER is the person logging in to the system to use it.
//          The User has three unique identifiers:
//            Id
//            Handle
//            Email Address
//          The Id never changes. Handle is probably not needed, but was
//          added so that a User could change their e-mail address.
////////////////////////////////////////////////////////////////////////////////

// IdentityService provides methods to create, update, and authenticate users.
type IdentityService interface {
	// AuthenticateCredentials verifies the credentials.
	// It returns the identity data with an authorization token if successful.
	AuthenticateCredentials(AuthenticateCredentialsRequest) IdentityResponse
	// AuthorizeToken blah
	AuthorizeToken(AuthorizeTokenRequest) RolesResponse
	// CreateIdentity creates a new user and returns the user data if successful.
	CreateIdentity(CreateIdentityRequest) IdentityResponse
	// DeleteIdentity deletes an existing user and returns an OK response if successful.
	DeleteIdentity(DeleteIdentityRequest) OKResponse
	// FetchIdentity retrieves an existing user and returns the user data if successful.
	FetchIdentity(FetchIdentityRequest) IdentityResponse
	// UpdateIdentity updates an existing identity and returns the updated data if successful.
	UpdateIdentity(UpdateIdentityRequest) IdentityResponse
}

// AuthenticateCredentialsRequest is the request object for IdentityService.AuthenticateCredentials.
// The caller must supply an email or handle, never both.
type AuthenticateCredentialsRequest struct {
	// Email is the e-mail address associated with the identity.
	// Optional - omit or set to null if not using.
	// example: "fred.flintrock@example.com"
	Email *string
	// Handle is the nickname associated with the identity.
	// Optional - omit or set to null if not using.
	// example: "flintrock"
	Handle *string
	// Secret is the hex-encoded passphrase used to authenticate the request.
	// Required.
	// example: "6261644d6f6f7365"
	Secret string
}

// AuthorizeTokenRequest is the request object for IdentityService.AuthorizeToken.
type AuthorizeTokenRequest struct {
	// Token is the signed JWT to validate.
	// Required.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string
}

// CreateIdentityRequest is the request object for IdentityService.CreateIdentity.
type CreateIdentityRequest struct {
	// Email is the e-mail address associated with the identity.
	// Required.
	// example: "fred.flintrock@example.com"
	Email string
	// Handle is the nickname associated with the identity.
	// Required.
	// example: "flintrock"
	Handle string
	// Secret is the hex-encoded passphrase to use when authenticating the identity.
	// Required.
	// example: "6261644d6f6f7365"
	Secret string
}

// DeleteIdentityRequest is the request object for IdentityService.DeleteIdentity.
type DeleteIdentityRequest struct {
	// Id is the unique identifier for the identity.
	// Required.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
}

// FetchIdentityRequest is the request object for IdentityService.FetchIdentity.
type FetchIdentityRequest struct {
	// Id is the unique identifier for the identity.
	// Optional - omit or set to null if not querying on.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address associated with the identity.
	// Optional - omit or set to null if not querying on.
	// example: "fred.flintrock@example.com"
	Email *string
	// Handle is the nickname associated with the identity.
	// Optional - omit or set to null if not querying on.
	// example: "flintrock"
	Handle *string
}

// UpdateIdentityRequest is the request object for IdentityService.UpdateIdentity.
type UpdateIdentityRequest struct {
	// Id is the unique identifier for the identity.
	// Required.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address associated with the identity.
	// Optional - omit or set to null if not updating.
	// example: "fred.flintrock@example.com"
	Email *string
	// Handle is the nickname associated with the identity.
	// Optional - omit or set to null if not updating.
	// example: "flintrock"
	Handle *string
	// Secret is the hex-encoded passphrase used to authenticate the request.
	// Optional - omit or set to null if not updating.
	// example: "6261644d6f6f7365"
	Secret *string
}

// IdentityResponse is the response object containing the identity and an authorization token if authenticated.
type IdentityResponse struct {
	// Id is the unique identifier for the identity.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address associated with the identity.
	// example: "fred.flintrock@example.com"
	Email string
	// Handle is the nickname associated with the identity.
	// example: "flintrock"
	Handle string
	// Roles is the list of roles granted to the identity.
	// example: []{"authenticated"}
	Roles []string
	// Token is the signed JWT if the identity is authenticated.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string
}

// OKResponse is the response object.
type OKResponse struct {
	// OK is a boolean.
	// example: true
	OK bool
}

// RolesResponse is the response object.
type RolesResponse struct {
	// Roles is the list of roles granted to the identity.
	// example: []string{"authenticated"}
	Roles []string
	// Token is the signed JWT if the identity is authenticated.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string
}

// AuthenticateUserRequest is the request object for IdentityService.AuthenticateUser.
// The caller must supply an email or handle, never both.
type AuthenticateUserRequest struct {
	// Email is the e-mail address the user registered with.
	// Optional - omit or set to null if not using.
	// example: "fred.flintrock@example.com"
	Email *string
	// Handle is the nickname for the user.
	// Optional - omit or set to null if not using.
	// example: "flintrock"
	Handle *string
	// Secret is the hex-encoded passphrase used to authenticate the request.
	// Required.
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
	// Handle is the nickname for the user.
	// example: "flintrock"
	Handle string
	// Roles is the list of roles granted to the user.
	// example: []{"authenticated"}
	Roles []string
	// Token is the signed JWT if the user is authenticated.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string
}

// AuthenticatedUserResponse is the response object containing the user's bearer token if authenticated.
type AuthenticatedUserResponse struct {
	// Id is the unique identifier for the user.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address the user registered with.
	// example: "fred.flintrock@example.com"
	Email string
	// Handle is the nickname for the user.
	// example: "flintrock"
	Handle string
	// Roles is the list of roles granted to the user.
	// example: []{"authenticated"}
	Roles []string
	// Token is the signed JWT if the user is authenticated.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string
}

// CreateUserRequest is the request object for IdentityService.CreateUser.
type CreateUserRequest struct {
	// Email is the e-mail address the user registered with.
	// Required.
	// example: "fred.flintrock@example.com"
	Email string
	// Handle is the nickname for the user.
	// Required.
	// example: "flintrock"
	Handle string
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
	// Optional - omit or set to null if not using.
	// example: "fred.flintrock@example.com"
	Email *string
}

// FetchUserRequest is the request object for IdentityService.FetchUser.
type FetchUserRequest struct {
	// Id is the identifier for the user to retrieve.
	// Optional - omit or set to null if not using.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id *string
	// Email is the e-mail address for the user to retrieve.
	// Optional - omit or set to null if not using.
	// example: "fred.flintrock@example.com"
	Email *string
	// Handle is the nickname for the user.
	// Optional - omit or set to null if not using.
	// example: "flintrock"
	Handle *string
}

// UpdateUserRequest is the request object for IdentityService.UpdateUser.
type UpdateUserRequest struct {
	// Id is the identifier for the user to update.
	// Required.
	// example: "fb6c1b87-41ef-4e92-91cc-1a5c59e5cd2d"
	Id string
	// Email is the e-mail address the user registered with.
	// Optional - omit or set to null if not updating.
	// example: "fred.flintrock@example.com"
	Email *string
	// Handle is the nickname for the user.
	// Optional - omit or set to null if not updating.
	// example: "flintrock"
	Handle *string
	// Secret is the hex-encoded passphrase used to authenticate the user on future requests.
	// Optional - omit or set to null if not updating.
	// example: "6261644d6f6f7365"
	Secret *string
}
