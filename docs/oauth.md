## OAuth and OpenID Connect

### 1. Authorization Request
```
// GetAuthorizeContext handles the OAuth2 / OIDC Authorization Endpoint.
//
// Normative (MUST / REQUIRED):
//   - MUST support HTTP GET.
//   - MAY support POST.
//   - MUST validate required parameters:
//     - response_type
//     - client_id
//     - redirect_uri (if multiple registered)
//     - scope
//   - MUST require "openid" scope for OIDC requests.
//   - MUST validate redirect_uri via exact string match.
//   - MUST return errors using redirect-based error response format.
//   - MUST echo `state` exactly if provided.
//   - MUST require and bind `nonce` for OIDC flows issuing ID tokens.
//   - MUST validate PKCE (code_challenge) if present.
//
// Conditionally Normative:
//   - PKCE REQUIRED for public clients.
//   - Consent screen REQUIRED unless prior consent exists.
//
// Non-Normative:
//   - Endpoint path (e.g., /authorize, /oauth2/auth).
//   - UI design and login screen behavior.
//   - Internal session management.
//   - Storage implementation.
//
```

#### Generate Authorization Code
```
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
//      - client_id
//      - redirect_uri
//      - authenticated user (resource owner)
//      - approved scope
//      - PKCE code_challenge (if present)
//   - Make the code single-use.
//   - Expire the code after a short duration (typically 30–120 seconds).
//   - Reject the code if:
//      - Expired
//      - Already used
//      - client_id does not match
//      - redirect_uri does not match
//      - PKCE verification fails
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
```

### 2. Token Request
```
// GetAccessToken handles the OAuth2 Token Endpoint.
//
// Normative (MUST / REQUIRED):
//   - MUST use HTTP POST.
//   - MUST require Content-Type: application/x-www-form-urlencoded.
//   - MUST require grant_type.
//   - MUST validate grant_type-specific parameters.
//   - MUST authenticate confidential clients.
//   - MUST validate authorization code:
//      - single use
//      - not expired
//      - bound to client_id
//      - bound to redirect_uri
//   - MUST validate PKCE code_verifier if used.
//   - MUST return JSON response.
//   - MUST include access_token and token_type ("Bearer").
//   - MUST issue id_token if scope includes "openid".
//
// Conditionally Normative:
//   - redirect_uri required if present in authorization request.
//   - refresh_token issuance depends on policy.
//
// Non-Normative:
//   - Endpoint path (e.g., /token).
//   - Access token format (opaque vs JWT).
//   - Token lifetime policy.
//
```

#### Generate Access Token
```
// GenerateOAuthAccessToken issues an OAuth2 access token.
//
// Normative:
//   - MUST represent authorization granted.
//   - MUST associate token with client_id.
//   - MUST bind token to resource owner (if applicable).
//   - MUST return token_type = "Bearer".
//
// NOT Normative (OAuth 2.0 does NOT mandate):
//   - Token format (opaque vs JWT).
//   - Claim structure.
//   - Encoding strategy.
//
// Best Practice (Not strictly normative):
//   - Include exp.
//   - Include audience (resource server).
//   - Include scope.
//   - Include jti for revocation.
//
// Non-Normative:
//   - Storage mechanism.
//   - Signing algorithm (unless JWT chosen).
//
```

#### Generate ID Token
```
// OAuthIDTokenClaims creates an OpenID Connect ID Token.
//
// Normative (MUST include claims):
//
//	iss   - issuer identifier (Who issued the token: auth server)
//	sub   - stable subject identifier (Who the user is: user_id)
//	aud   - client_id (Who the token is intended for: client_id)
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
```
