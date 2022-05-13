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

// AuthorizationService makes nice authorizations.
type AuthorizationService interface {
	// Authorize validates the id and secret and returns a signed JWT if successful.
	Authorize(AuthorizeRequest) AuthorizeResponse
	CreateSigningKeyRequest(CreateSigningKeyRequest) SigningKeyResponse
}

// AuthorizeRequest is the request object for AuthorizationService.Authorize.
type AuthorizeRequest struct {
	// Id is the account identifier to authenticate against.
	// example: "fred.flintrock@example.com"
	Id string
	// Secret is the hex-encoded passphrase used to authenticate the request.
	// example: "6261644d6f6f7365"
	Secret string
}

// AuthorizeResponse is the response object containing the signed JWT if authorized.
type AuthorizeResponse struct {
	// Token is the signed JWT if the request is authorized.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string
}

// CreateSigningKeyRequest is the request object for AuthorizationService.CreateSigningKey
type CreateSigningKeyRequest struct {
	// Id is the identifier to assign to the new key.
	// Example: "archer"
	Id string
	// PrivateKey is the hex-encoded key used to sign requests.
	// example: "feab7c54dec2"
	PrivateKey string
	// PublicKey is the hex-encoded key used to verify signatures.
	// example: "909d7a4202ff"
	PublicKey string
}

// SigningKeyResponse is the response object containing the public key to verify signatures
type SigningKeyResponse struct {
	// PublicKey is the hex-encoded key used to verify signatures.
	// example: "909d7a4202ff"
	PublicKey string
}
