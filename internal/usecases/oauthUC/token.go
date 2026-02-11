package oauthUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/oauth"

	"gorm.io/gorm"
)

func (uc *OAuthUsecase) GetAccessToken(
	ctx context.Context,
	req requestdto.OAuthAccessToken,
) (*responsedto.OAuthAccessToken, *ucerror.UCError) {
	client, err := uc.postgres.GetClientByID(ctx, req.ClientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ucerror.New(fmt.Errorf("client with ID:%s not found: %w", req.ClientID, err), consts.ErrBadRequest)
		}
		return nil, ucerror.New(fmt.Errorf("failed to get client with ID:%s: %w", req.ClientID, err), consts.ErrInternalServer)
	}

	if client.Secret != req.ClientSecret {
		return nil, ucerror.New(fmt.Errorf("incorrect client secret for client ID:%s", req.ClientID), consts.ErrBadRequest)
	}

	if client.RedirectURI != req.RedirectURI {
		return nil, ucerror.New(fmt.Errorf("redirect uri mismatch"), consts.ErrBadRequest)
	}

	claims, err := oauth.VerifyAuthorizationCode(req.ClientSecret, req.AuthorizationCode)
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("failed to verify authorization code: %w", err), consts.ErrBadRequest)
	}

	accessToken, err := uc.token.GenerateOAuthAccessToken(client, *claims)
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("failed to generate access token: %w", err), consts.ErrInternalServer)
	}

	return &responsedto.OAuthAccessToken{
		AccessToken:      accessToken,
		TokenType:        consts.TokenTypeBearer,
		ExpiresInMinutes: consts.AccessTokenValidityMinutes,
	}, nil
}
