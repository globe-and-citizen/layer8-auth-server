package tokenRepo

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/models"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/backend/pkg/oauth"
)

type ITokenRepository interface {
	GenerateUserJWTToken(user gormModels.User) (string, error)
	VerifyUserJWTToken(tokenString string) (models.UserClaims, error)
	GenerateClientJWTToken(client gormModels.Client) (string, error)
	VerifyClientJWTToken(tokenString string) (models.ClientClaims, error)
	GenerateOAuthJWTToken(user gormModels.User) (string, error)
	VerifyOAuthJWTToken(tokenString string) (models.OAuthAuthenticationClaims, error)
	GenerateOAuthAccessToken(client gormModels.Client, authClaims oauth.AuthorizationCodeClaims) (string, error)
	ParseOAuthAccessToken(tokenString string) (*models.ClientAccessTokenClaims, error)
	VerifyOAuthAccessToken(tokenString string, clientSecret []byte) error
}

func NewTokenRepository(userJWTSecret []byte, clientJWTSecret []byte, oauthJWTSecret []byte) ITokenRepository {
	return &TokenRepository{
		userJWTSecret:   userJWTSecret,
		clientJWTSecret: clientJWTSecret,
		oauthJWTSecret:  oauthJWTSecret,
	}
}

type TokenRepository struct {
	userJWTSecret   []byte
	clientJWTSecret []byte
	oauthJWTSecret  []byte
}
