package authUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/authRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/tokenRepo"
	"net/http"
)

type IAuthorizationUsecase interface {
	ValidateAndGenerateAccessToken(req requestdto.AuthorizationToken) (responsedto.OauthToken, int, string, error)
}

func NewAuthorizationUsecase(
	postgres postgresRepo.IPostgresRepository,
	oauth authRepo.IAuthorizationRepository,
	token tokenRepo.ITokenRepository,
) IAuthorizationUsecase {
	return &AuthorizationUsecase{
		postgres: postgres,
		oauth:    oauth,
		token:    token,
	}
}

type AuthorizationUsecase struct {
	postgres postgresRepo.IPostgresRepository
	oauth    authRepo.IAuthorizationRepository
	token    tokenRepo.ITokenRepository
}

func (uc *AuthorizationUsecase) ValidateAndGenerateAccessToken(req requestdto.AuthorizationToken) (responsedto.OauthToken, int, string, error) {
	client, err := uc.postgres.GetClientByID(req.ClientUUID)
	if err != nil {
		return responsedto.OauthToken{}, http.StatusUnauthorized, "failed to authenticate client", err
	}

	if client.Secret != req.ClientSecret {
		return responsedto.OauthToken{}, http.StatusUnauthorized, "failed to authenticate client", fmt.Errorf("provided secret value is invalid")
	}

	claims, err := uc.oauth.VerifyAuthorizationCode(req.ClientSecret, req.AuthorizationCode)
	if err != nil {
		return responsedto.OauthToken{}, http.StatusBadRequest, "the authorization code is invalid", err
	}

	accessToken, err := uc.token.GenerateClientAuthJWTToken(client, *claims)
	if err != nil {
		return responsedto.OauthToken{}, http.StatusInternalServerError, "internal error when generating the access token", err
	}

	return responsedto.OauthToken{
		AccessToken:      accessToken,
		TokenType:        consts.TokenTypeBearer,
		ExpiresInMinutes: consts.AccessTokenValidityMinutes,
	}, http.StatusOK, "", nil
}
