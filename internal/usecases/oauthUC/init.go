package oauthUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/errors"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"time"
)

type IOAuthUsecase interface {
	PrecheckUserLogin(ctx context.Context, req requestdto.OAuthUserLoginPrecheck) (responsedto.OAuthUserLoginPrecheck, error)
	UserLogin(ctx context.Context, req requestdto.OAuthUserLogin) (responsedto.OAuthUserLogin, error)
	AuthorizeContext(ctx context.Context, req requestdto.OAuthAuthorizeContext) (*responsedto.OAuthAuthorizeContext, *errors.OAuthError)
	AuthorizeDecision(ctx context.Context, req requestdto.OAuthAuthorizeDecision, userID uint, authzCodeExpiry time.Duration) (*responsedto.OAuthAuthorizeDecision, *errors.OAuthError)
	GetAccessToken(ctx context.Context, req requestdto.OAuthAccessToken) (*responsedto.OAuthAccessToken, *errors.OAuthError)
	GetZkUserMetadata(ctx context.Context, req requestdto.OAuthZkMetadata) (*responsedto.OAuthZkMetadata, *errors.OAuthError)
	VerifyOAuthJWTToken(ctx context.Context, tokenString string) (userID uint, userUsername string, err error)
	VerifyAccessToken(ctx context.Context, tokenString string) (userID uint, scopes string, err error)
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
