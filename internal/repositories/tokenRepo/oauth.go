package tokenRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (t TokenRepository) GenerateOAuthJWTToken(user gormModels.User, expiry time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   user.Username, // The value was originally user.ID; it I changed it to Username to avoid type conversion overhead
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
	}

	return t.generateJWTToken(claims, t.oauthJWTSecret)
}

func (t TokenRepository) VerifyOAuthJWTToken(tokenString string) (*models.OAuthAuthenticationClaims, error) {
	claims := &models.OAuthAuthenticationClaims{}
	err := t.verifyJWTToken(tokenString, t.clientJWTSecret, claims)
	return claims, err
}

func (t TokenRepository) GenerateOAuthAccessToken(
	clientID string, scopes string, userID uint, secret []byte, expiry time.Duration,
) (string, error) {
	claims := models.OAuthAccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   clientID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry).UTC()),
		},
		Scopes: scopes,
		UserID: userID,
	}

	return t.generateJWTToken(claims, secret)
}

func (t TokenRepository) VerifyOAuthAccessToken(tokenString string, secret []byte) (*models.OAuthAccessTokenClaims, error) {
	claims := &models.OAuthAccessTokenClaims{}
	err := t.verifyJWTToken(tokenString, secret, claims)
	return claims, err
}

func (t TokenRepository) GenerateOAuthIDToken(
	userID uint, clientID, nonce string, metadata models.OIDCUserProfile, secret []byte, expiry time.Duration,
) (string, error) {
	claims := &models.OAuthIDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.jwtIssuer,
			Subject:   fmt.Sprintf("%d", userID), // Convert uint to string for the 'sub' claim
			Audience:  jwt.ClaimStrings{clientID},
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry).UTC()),
		},
		Nonce:           nonce,
		AuthTime:        time.Now().UTC(),
		OIDCUserProfile: metadata,
	}
	return t.generateJWTToken(claims, secret)
}

func (t TokenRepository) VerifyOAuthIDToken(tokenString string) (*models.OAuthIDTokenClaims, error) {
	claims := &models.OAuthIDTokenClaims{}
	err := t.verifyJWTToken(tokenString, t.oauthJWTSecret, claims)
	return claims, err
}
