package tokenRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
)

type ITokenRepository interface {
	GenerateUserJWTToken(user gormModels.User) (string, error)
	VerifyUserJWTToken(tokenString string) (models.UserClaims, error)
	GenerateClientJWTToken(client gormModels.Client) (string, error)
	VerifyClientJWTToken(tokenString string) (models.ClientClaims, error)
}

func NewTokenRepository(userJWTSecret []byte, clientJWTSecret []byte) ITokenRepository {
	return &TokenRepository{
		userJWTSecret:   userJWTSecret,
		clientJWTSecret: clientJWTSecret,
	}
}

type TokenRepository struct {
	userJWTSecret   []byte
	clientJWTSecret []byte
}
