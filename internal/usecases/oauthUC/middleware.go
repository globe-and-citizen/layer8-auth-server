package oauthUC

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
)

func (uc *OAuthUsecase) MdwVerifyUserLoggedInToken(ctx context.Context, tokenString string) (uint, string, *ucerror.UCError) {
	claims, err := uc.token.VerifyOAuthJWTToken(tokenString)
	if err != nil {
		return 0, "", ucerror.New(fmt.Errorf("failed to verify oauth user logged in token: %w", err), consts.ErrUnauthorized)
	}

	// verify user by username
	user, err := uc.postgres.GetUserByUsername(ctx, claims.Subject)
	if err != nil {
		return 0, "", ucerror.New(fmt.Errorf("failed to find user: %w", err), consts.ErrUnauthorized)
	}

	// todo verify the rest claims

	return user.ID, user.Username, nil
}

func (uc *OAuthUsecase) MdwVerifyClientAccessToken(ctx context.Context, tokenString string) (uint, string, *ucerror.UCError) {
	claims, err := uc.token.ParseOAuthAccessToken(tokenString)
	if err != nil {
		return 0, "", ucerror.New(fmt.Errorf("failed to parse client access token: %w", err), consts.ErrUnauthorized)
	}

	// validate clientID
	client, err := uc.postgres.GetClientByID(ctx, claims.Subject)
	if err != nil {
		return 0, "", ucerror.New(fmt.Errorf("client not found: %w", err), consts.ErrUnauthorized)
	}

	err = uc.token.VerifyOAuthAccessToken(tokenString, []byte(client.Secret))
	if err != nil {
		return 0, "", ucerror.New(fmt.Errorf("invalid access token: %w", err), consts.ErrUnauthorized)
	}

	return claims.UserID, claims.Scopes, nil
}
