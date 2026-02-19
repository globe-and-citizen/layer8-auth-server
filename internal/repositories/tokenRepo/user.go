package tokenRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (t TokenRepository) GenerateUserJWTToken(user gormModels.User, expiry time.Duration) (string, error) {
	claims := &models.UserClaims{
		Username: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			Issuer:    t.JWTIssuer,
		},
	}

	return t.generateJWTToken(claims, t.userJWTSecret)
}

func (t TokenRepository) VerifyUserJWTToken(tokenString string) (*models.UserClaims, error) {
	claims := &models.UserClaims{}
	err := t.verifyJWTToken(tokenString, t.userJWTSecret, claims)
	return claims, err
}
