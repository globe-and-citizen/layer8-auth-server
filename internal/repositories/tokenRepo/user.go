package tokenRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (t TokenRepository) GenerateUserJWTToken(user gormModels.User) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &models.UserClaims{
		Username: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "GlobeAndCitizen",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.userJWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t TokenRepository) VerifyUserJWTToken(tokenString string) (models.UserClaims, error) {
	claims := &models.UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return t.clientJWTSecret, nil
	})
	if err != nil {
		return models.UserClaims{}, err
	}

	if !token.Valid {
		return models.UserClaims{}, fmt.Errorf("invalid token")
	}

	return *claims, nil
}
