package tokenRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"

	layer8Utils "github.com/globe-and-citizen/layer8-utils"
)

type ITokenRepository interface {
	GenerateUserJWTToken(user gormModels.User) (string, error)
	VerifyUserJWTToken(tokenString string) (models.UserClaims, error)
	GenerateClientJWTToken(client gormModels.Client) (string, error)
	VerifyClientJWTToken(tokenString string) (models.ClientClaims, error)
	GenerateClientOauthJWTToken(client gormModels.Client, authClaims layer8Utils.AuthCodeClaims) (string, error)
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
