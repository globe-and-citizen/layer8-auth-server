package tokenRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (t TokenRepository) GenerateClientJWTToken(client gormModels.Client, expiry time.Duration) (string, error) {
	claims := &models.ClientClaims{
		Username: client.Username,
		ClientID: client.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			Issuer:    t.jwtIssuer,
		},
	}

	return t.generateJWTToken(claims, t.clientJWTSecret)
}

func (t TokenRepository) VerifyClientJWTToken(tokenString string) (*models.ClientClaims, error) {
	claims := &models.ClientClaims{}
	err := t.verifyJWTToken(tokenString, t.clientJWTSecret, claims)
	return claims, err
}
