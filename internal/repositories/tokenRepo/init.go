package tokenRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ITokenRepository interface {
	GenerateUserJWTToken(user gormModels.User, expiry time.Duration) (string, error)
	VerifyUserJWTToken(tokenString string) (*models.UserClaims, error)
	GenerateClientJWTToken(client gormModels.Client, expiry time.Duration) (string, error)
	VerifyClientJWTToken(tokenString string) (*models.ClientClaims, error)
	GenerateOAuthJWTToken(user gormModels.User, expiry time.Duration) (string, error)
	VerifyOAuthJWTToken(tokenString string) (*models.OAuthAuthenticationClaims, error)

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
	GenerateOAuthAccessToken(clientID string, scopes string, userID uint, secret []byte, expiry time.Duration) (string, error)
	VerifyOAuthAccessToken(tokenString string, secret []byte) (*models.OAuthAccessTokenClaims, error)

	// GenerateOAuthIDToken creates an OpenID Connect ID Token.
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
	GenerateOAuthIDToken(userID uint, clientID, nonce string, metadata models.OIDCUserProfile, secret []byte, expiry time.Duration) (string, error)
	VerifyOAuthIDToken(tokenString string) (*models.OAuthIDTokenClaims, error)
}

func NewTokenRepository(jwtIssuer string, userJWTSecret []byte, clientJWTSecret []byte, oauthJWTSecret []byte) ITokenRepository {
	return &TokenRepository{
		jwtIssuer:       jwtIssuer,
		userJWTSecret:   userJWTSecret,
		clientJWTSecret: clientJWTSecret,
		oauthJWTSecret:  oauthJWTSecret,
	}
}

type TokenRepository struct {
	jwtIssuer       string
	userJWTSecret   []byte
	clientJWTSecret []byte
	oauthJWTSecret  []byte
}

func (t TokenRepository) generateJWTToken(claims jwt.Claims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", utils.StackError(err)
	}

	return tokenString, nil
}

func (t TokenRepository) verifyJWTToken(tokenString string, secret []byte, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return utils.StackError(fmt.Errorf("invalid token"))
	}

	return nil
}
