package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

type ClientClaims struct {
	Username string `json:"username"`
	ClientID string `json:"user_id"`
	jwt.RegisteredClaims
}

// OAuthAuthenticationClaims represents the claims of an OAuth authentication token
type OAuthAuthenticationClaims struct {
	jwt.RegisteredClaims
}

type OAuthAccessTokenClaims struct {
	UserID uint   `json:"user_id"`
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

// OAuthIDTokenClaims creates an OpenID Connect ID Token.
//
// Normative (MUST include claims):
//
//	iss   - issuer identifier
//	sub   - stable subject identifier
//	aud   - client_id
//	exp   - expiration time
//	iat   - issued-at time
//
// Conditionally Required:
//
//	nonce      - if provided in auth request
//	auth_time  - if max_age requested
//	at_hash    - if access_token returned in same response
//
// Normative Requirements:
//   - MUST be a JWT.
//   - MUST be signed.
//   - MUST use agreed algorithm.
//   - MUST validate audience and issuer consistency.
//
// Non-Normative:
//   - Claim ordering.
//   - Internal subject storage model.
//   - Token lifetime duration.
type OAuthIDTokenClaims struct {
	jwt.RegisteredClaims
	AuthTime time.Time `json:"auth_time"`
	Nonce    string    `json:"nonce"`
}
