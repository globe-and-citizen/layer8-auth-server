package oauthUC

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"time"
)

type IOAuthUsecase interface {
	PrecheckUserLogin(ctx context.Context, req requestdto.OAuthUserLoginPrecheck) (responsedto.OAuthUserLoginPrecheck, *ucerror.UCError)
	UserLogin(ctx context.Context, req requestdto.OAuthUserLogin) (responsedto.OAuthUserLogin, *ucerror.UCError)
	GetAuthorizeContext(ctx context.Context, req requestdto.OAuthAuthorizeContext) (*responsedto.OAuthAuthorizeContext, *ucerror.UCError)
	PostAuthorizeDecision(ctx context.Context, req requestdto.OAuthAuthorizeDecision, userID uint, authzCodeExpiry time.Duration) (*responsedto.OAuthAuthorizeDecision, *ucerror.UCError)
	GetAccessToken(ctx context.Context, req requestdto.OAuthAccessToken) (*responsedto.OAuthAccessToken, *ucerror.UCError)
	GetZkUserMetadata(ctx context.Context, req requestdto.OAuthZkMetadata) (*responsedto.OAuthZkMetadata, *ucerror.UCError)
	MdwVerifyUserLoggedInToken(ctx context.Context, tokenString string) (userID uint, userUsername string, err *ucerror.UCError)
	MdwVerifyClientAccessToken(ctx context.Context, tokenString string) (userID uint, scopes string, err *ucerror.UCError)
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
