package oauthUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"time"

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

	authzCode, ucerr := uc.verifyAuthorizationCode(ctx, req.AuthorizationCode)
	if ucerr != nil {
		return nil, ucerr
	}

	accessToken, err := uc.token.GenerateOAuthAccessToken(
		client.ID, authzCode.Scopes, authzCode.UserID, []byte(uc.config.AccessTokenSecret), uc.config.AccessTokenExpiry,
	)
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("failed to generate access token: %w", err), consts.ErrInternalServer)
	}

	err = uc.postgres.DeleteOAuthAuthorizationCode(ctx, authzCode.Code)
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("failed to delete authorization code: %w", err), consts.ErrInternalServer)
	}

	return &responsedto.OAuthAccessToken{
		AccessToken:      accessToken,
		TokenType:        consts.TokenTypeBearer,
		ExpiresInMinutes: int(uc.config.AccessTokenExpiry.Minutes()),
	}, nil
}

func (uc *OAuthUsecase) verifyAuthorizationCode(
	ctx context.Context, code string,
) (*gormModels.OAuthAuthorizationCode, *ucerror.UCError) {
	authzCode, err := uc.postgres.GetOAuthAuthorizationCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ucerror.New(fmt.Errorf("authorization code not found: %w", err), consts.ErrBadRequest)
		}
		if errors.Is(err, context.Canceled) {
			return nil, ucerror.New(fmt.Errorf("request canceled: %w", err), consts.ErrInternalServer)
		}
		return nil, ucerror.New(fmt.Errorf("failed to get authorization code: %w", err), consts.ErrInternalServer)
	}

	if authzCode.ExpiresAt < time.Now().Unix() {
		return nil, ucerror.New(fmt.Errorf("authorization code expired"), consts.ErrBadRequest)
	}

	return authzCode, nil
}
