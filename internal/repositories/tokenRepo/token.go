package tokenRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/models"
	gormModels2 "globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/oauth"
)

type ITokenRepository interface {
	GenerateUserJWTToken(user gormModels2.User) (string, error)
	VerifyUserJWTToken(tokenString string) (models.UserClaims, error)
	GenerateClientJWTToken(client gormModels2.Client) (string, error)
	VerifyClientJWTToken(tokenString string) (models.ClientClaims, error)
	GenerateOAuthJWTToken(user gormModels2.User) (string, error)
	VerifyOAuthJWTToken(tokenString string) (models.OAuthAuthenticationClaims, error)
	GenerateOAuthAccessToken(client gormModels2.Client, authClaims oauth.AuthorizationCodeClaims) (string, error)
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
