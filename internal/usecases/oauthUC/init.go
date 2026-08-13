package oauthUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/log"
)

type IOAuthUsecase interface {
	PrecheckUserLogin(ctx context.Context, req requestdto.OAuthUserLoginPrecheck) (*responsedto.OAuthUserLoginPrecheck, *ucerror.UCError)
	UserLogin(ctx context.Context, req requestdto.OAuthUserLogin) (*responsedto.OAuthUserLogin, *ucerror.UCError)
	GetAuthorizeContext(ctx context.Context, req requestdto.OAuthAuthorizeQueries) (*responsedto.OAuthAuthorizeContext, *ucerror.UCError)
	PostAuthorizeDecision(ctx context.Context, req requestdto.OAuthAuthorizeConsent, userID uint) (*responsedto.OAuthAuthorizeConsent, *ucerror.UCError)
	RequestToken(ctx context.Context, req requestdto.OAuthTokenRequest) (*responsedto.OAuthTokenRequest, *ucerror.UCError)
	GetZkUserMetadata(ctx context.Context, req requestdto.OAuthZkMetadata) (*responsedto.OAuthZkMetadata, *ucerror.UCError)
	MdwVerifyUserLoggedInToken(ctx context.Context, tokenString string) (userID uint, userUsername string, err *ucerror.UCError)
	MdwVerifyClientAccessToken(ctx context.Context, tokenString string) (userID uint, scopes string, err *ucerror.UCError)
}

func NewOAuthUsecase(logger log.ILogger, conf config.OAuthConfig, postgres postgresRepo.IPostgresRepository, token tokenRepo.ITokenRepository) IOAuthUsecase {
	return &OAuthUsecase{
		logger:   logger,
		config:   conf,
		postgres: postgres,
		token:    token,
	}
}

type OAuthUsecase struct {
	logger   log.ILogger
	config   config.OAuthConfig
	postgres postgresRepo.IPostgresRepository
	token    tokenRepo.ITokenRepository
}
