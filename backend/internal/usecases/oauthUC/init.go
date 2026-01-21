package oauthUC

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/internal/errors"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/tokenRepo"
	"time"
)

type IOAuthUsecase interface {
	PrecheckUserLogin(req requestdto.OAuthUserLoginPrecheck) (responsedto.OAuthUserLoginPrecheck, error)
	UserLogin(req requestdto.OAuthUserLogin) (responsedto.OAuthUserLogin, error)
	AuthorizeContext(req requestdto.OAuthAuthorizeContext) (*responsedto.OAuthAuthorizeContext, *errors.OAuthError)
	AuthorizeDecision(req requestdto.OAuthAuthorizeDecision, userID uint, authzCodeExpiry time.Duration) (*responsedto.OAuthAuthorizeDecision, *errors.OAuthError)
	GetAccessToken(req requestdto.OAuthAccessToken) (*responsedto.OAuthAccessToken, *errors.OAuthError)
	GetZkUserMetadata(req requestdto.OAuthZkMetadata) (*responsedto.OAuthZkMetadata, *errors.OAuthError)
	VerifyOAuthJWTToken(tokenString string) (userID uint, userUsername string, err error)
	VerifyAccessToken(tokenString string) (userID uint, scopes string, err error)
}

func NewOAuthUsecase(postgres postgresRepo.IPostgresRepository, token tokenRepo.ITokenRepository) IOAuthUsecase {
	return &OAuthUsecase{
		postgres: postgres,
		token:    token,
	}
}

type OAuthUsecase struct {
	postgres postgresRepo.IPostgresRepository
	token    tokenRepo.ITokenRepository
}
