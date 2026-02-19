package oauth

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type AuthorizationCodeClaims struct {
	ClientID    string            `json:"cid"`
	UserID      uint              `json:"uid"`
	RedirectURI string            `json:"ruri"`
	Scopes      string            `json:"scp"`
	ExpiresAt   int64             `json:"exp"`
	HeaderMap   map[string]string `json:"hmap"`
	jwt.RegisteredClaims
}

type AuthURL struct {
	URL          string
	State        string
	Code         string
	CodeVerifier string
}

// String returns the URL to redirect the user to for authentication
func (u *AuthURL) String() string {
	urlc, _ := url.Parse(u.URL)
	q := urlc.Query()
	redirectURI := q.Get("redirect_uri")
	q.Del("client_id")
	q.Del("redirect_uri")
	urlc.RawQuery = q.Encode()
	return redirectURI + urlc.String()
}

func GenerateAuthURL(
	clientID string,
	code string,
	clientRedirectURI string,
	scopes []string,
) (string, error) {
	state, stateErr := utils.GenerateRandomBase64String(24)
	if stateErr != nil {
		return "", fmt.Errorf("could not generate random state: %v", stateErr)
	}

	config := oauth2.Config{
		ClientID:    clientID,
		RedirectURL: clientRedirectURI,
		Scopes:      scopes,
	}

	authURL := AuthURL{
		URL: config.AuthCodeURL(
			state,
			oauth2.SetAuthURLParam("code", code),
		),
		Code:  code,
		State: state,
	}

	return authURL.String(), nil
}

// GenerateAuthorizationCode Authorization Code Flow — Code Creation & Validation
//
// Specification References:
//   - OAuth 2.0 (RFC 6749)
//   - OpenID Connect Core 1.0
//
// Overview:
//
//	The authorization code is a short-lived, single-use credential
//	representing successful resource owner authentication and client
//	authorization. It is issued by the Authorization Endpoint and
//	redeemed at the Token Endpoint.
//
// ---------------------------------------------------------------------
// Authorization Code — Normative Requirements
// ---------------------------------------------------------------------
//
// The server MUST:
//
//   - Generate a cryptographically strong, high-entropy value.
//   - Ensure the code is unpredictable and URL-safe.
//   - Bind the code to:
//   - client_id
//   - redirect_uri
//   - authenticated user (resource owner)
//   - approved scope
//   - PKCE code_challenge (if present)
//   - Make the code single-use.
//   - Expire the code after a short duration (typically 30–120 seconds).
//   - Reject the code if:
//   - Expired
//   - Already used
//   - client_id does not match
//   - redirect_uri does not match
//   - PKCE verification fails
//
// The server MUST NOT:
//
//   - Embed sensitive data directly inside the code unless protected.
//   - Allow reuse of the code.
//   - Accept codes issued to a different client.
//   - Accept codes with mismatched redirect_uri.
//
// ---------------------------------------------------------------------
// Authorization Code — Security Properties
// ---------------------------------------------------------------------
//
//   - Opaque (recommended)
//   - High entropy (>=128 bits, 256 bits preferred)
//   - Short lifetime
//   - One-time redeemable
//
// Recommended entropy example:
//
//	32 random bytes (256 bits)
//	base64.RawURLEncoding encoding
//
// ---------------------------------------------------------------------
// Lifecycle
// ---------------------------------------------------------------------
//
//  1. Authorization Endpoint:
//     - User authenticates
//     - Client + redirect validated
//     - Consent handled
//     - Code generated and stored
//     - Redirect to client with:
//     ?code=XYZ&state=abc
//
//  2. Token Endpoint:
//     - Client submits code
//     - Server validates binding and expiry
//     - Server invalidates code
//     - Tokens are issued
//
// ---------------------------------------------------------------------
// Storage Model (Recommended)
// ---------------------------------------------------------------------
//
//	type AuthorizationCode struct {
//	    Code                string
//	    ClientID            string
//	    RedirectURI         string
//	    UserID              string
//	    Scope               string
//	    CodeChallenge       string
//	    CodeChallengeMethod string
//	    ExpiresAt           time.Time
//	    Used                bool
//	}
//
// Codes SHOULD be deleted or marked used immediately after successful redemption.
//
// ---------------------------------------------------------------------
// Implementation Notes (Non-Normative)
// ---------------------------------------------------------------------
//
//   - Store codes in memory (single-node) or shared store (Redis/DB)
//     for distributed systems.
//   - Avoid JWT authorization codes unless carefully designed.
//   - Keep TTL small to reduce interception risk.
//   - Always enforce strict redirect_uri string matching.
func GenerateAuthorizationCode(
	clientID string,
	clientSecret string,
	clientRedirectURI string,
	scopes []string,
	userID uint,
	expiry time.Duration,
) (string, error) {
	claims := AuthorizationCodeClaims{
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: clientRedirectURI,
		Scopes:      strings.Join(scopes, ","), // fixme: should the separator be a space?
		ExpiresAt:   time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	code, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", fmt.Errorf("could not generate auth code: %s", err)
	}

	return code, nil
}

func VerifyAuthorizationCode(secret string, code string) (*AuthorizationCodeClaims, error) {
	authClaims, err := DecodeAuthCode(secret, code)
	if err != nil {
		return authClaims, fmt.Errorf("failed to decode auth code: %v", err)
	}

	return authClaims, nil
}

func DecodeAuthCode(secret, code string) (*AuthorizationCodeClaims, error) {
	token, err := jwt.ParseWithClaims(code, &AuthorizationCodeClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not decode auth code: %s", err)
	}

	claims, ok := token.Claims.(*AuthorizationCodeClaims)
	if !ok {
		return nil, fmt.Errorf("could not decode auth code: %s", err)
	}

	//if claims.ExpiresAt < time.Now().Unix() {
	//	return nil, fmt.Errorf("could not decode auth code: token expired")
	//}
	return claims, nil
}
