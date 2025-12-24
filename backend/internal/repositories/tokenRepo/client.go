package tokenRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/models"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (t TokenRepository) GenerateClientJWTToken(client gormModels.Client) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.ClientClaims{
		Username: client.Username,
		ClientID: client.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "GlobeAndCitizen",
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(t.clientJWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t TokenRepository) VerifyClientJWTToken(tokenString string) (models.ClientClaims, error) {
	claims := &models.ClientClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return t.clientJWTSecret, nil
	})
	if err != nil {
		return models.ClientClaims{}, err
	}

	if !token.Valid {
		return models.ClientClaims{}, fmt.Errorf("invalid token")
	}

	return *claims, nil
}
