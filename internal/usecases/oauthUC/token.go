package oauthUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/oauth"
	"time"

	"gorm.io/gorm"
)

func (uc *OAuthUsecase) RequestToken(
	ctx context.Context,
	req requestdto.OAuthTokenRequest,
) (*responsedto.OAuthTokenRequest, *ucerror.UCError) {
	client, authzCode, ucerr := uc.validateTokenRequest(ctx, req)
	if ucerr != nil {
		return nil, ucerr
	}

	// delete the authorization code after use, as it can only be used once.
	// We do this before generating the access token to avoid the case where the access token generation succeeds
	// but the authorization code deletion fails, which would leave a valid authorization code that can be reused.
	err := uc.postgres.DeleteOAuthAuthorizationCode(ctx, authzCode.Code)
	if err != nil {
		return nil, ucerror.New(
			fmt.Errorf("failed to delete authorization code: %w", err),
			consts.ErrInternalServer,
		)
	}

	scopes := parseScopes(authzCode.Scopes)
	accessToken, err := uc.token.GenerateOAuthAccessToken(
		client.ID,
		scopes.APIScopes().String(),
		authzCode.UserID,
		[]byte(uc.config.AccessTokenSecret),
		uc.config.AccessTokenExpiry,
	)
	if err != nil {
		return nil, ucerror.New(
			fmt.Errorf("failed to generate access token: %w", err),
			consts.ErrInternalServer,
		)
	}

	res := &responsedto.OAuthTokenRequest{
		AccessToken:      accessToken,
		TokenType:        consts.TokenTypeBearer,
		ExpiresInMinutes: int(uc.config.AccessTokenExpiry.Minutes()),
	}

	if scopes.IsOIDC() {
		res.IDToken, ucerr = uc.generateIDToken(ctx, authzCode.UserID, client.ID, authzCode.Nonce, scopes.OIDCScopes())
		if ucerr != nil {
			return nil, ucerr
		}
	}

	return res, nil
}

func (uc *OAuthUsecase) generateIDToken(ctx context.Context, userID uint, clientID, nonce string, scopes Scopes) (string, *ucerror.UCError) {
	// todo how to do with multiple scopes?
	user, userMetadata, err := uc.postgres.GetUserProfile(ctx, userID)
	if err != nil {
		return "", ucerror.New(
			fmt.Errorf("failed to get user profile for ID token generation: %w", err),
			consts.ErrInternalServer,
		)
	}

	profile := models.OIDCUserProfile{
		Username:    user.Username,
		DisplayName: userMetadata.DisplayName,
		Bio:         userMetadata.Bio,
	}

	idToken, err := uc.token.GenerateOAuthIDToken(
		userID, clientID, nonce, profile, []byte(uc.config.IDTokenSecret), uc.config.IDTokenExpiry,
	)
	if err != nil {
		return "", ucerror.New(
			fmt.Errorf("failed to generate ID token: %w", err),
			consts.ErrInternalServer,
		)
	}

	return idToken, nil
}

// validateTokenRequest checks:
// - grant_type is supported,
// - client ID exists, and if the provided client secret and redirect URI match the stored values for that client ID.
// - authorization code is valid, not expired, and belongs to the client ID.
func (uc *OAuthUsecase) validateTokenRequest(
	ctx context.Context, req requestdto.OAuthTokenRequest,
) (*gormModels.Client, *gormModels.OAuthAuthorizationCode, *ucerror.UCError) {
	if req.GrantType != string(oauth.GrantTypeAuthorizationCode) {
		return nil, nil, ucerror.New(
			fmt.Errorf("unsupported grant type: %s", req.GrantType),
			consts.ErrBadRequest,
		)
	}

	client, err := uc.postgres.GetClientByID(ctx, req.ClientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ucerror.New(
				fmt.Errorf("client with ID:%s not found: %w", req.ClientID, err),
				consts.ErrBadRequest,
			)
		}
		return nil, nil, ucerror.New(
			fmt.Errorf("failed to get client with ID:%s: %w", req.ClientID, err),
			consts.ErrInternalServer,
		)
	}

	if client.Secret != req.ClientSecret {
		return nil, nil, ucerror.New(
			fmt.Errorf("incorrect client secret for client ID:%s", req.ClientID),
			consts.ErrBadRequest,
		)
	}

	if client.RedirectURI != req.RedirectURI {
		return nil, nil, ucerror.New(fmt.Errorf("redirect uri mismatch"), consts.ErrBadRequest)
	}

	authzCode, ucerr := uc.verifyAuthorizationCode(ctx, req.AuthorizationCode, client.ID)
	if ucerr != nil {
		return nil, nil, ucerr
	}

	return client, authzCode, nil
}

// verifyAuthorizationCode checks if the provided authorization code is valid, not expired,
// and belongs to the given client ID.
// It also checks if the user associated with the authorization code still exists.
func (uc *OAuthUsecase) verifyAuthorizationCode(
	ctx context.Context, code string, clientID string,
) (*gormModels.OAuthAuthorizationCode, *ucerror.UCError) {
	authzCode, err := uc.postgres.GetOAuthAuthorizationCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ucerror.New(
				fmt.Errorf("authorization code not found: %w", err),
				consts.ErrBadRequest,
			)
		}
		if errors.Is(err, context.Canceled) {
			return nil, ucerror.New(
				fmt.Errorf("request canceled: %w", err),
				consts.ErrInternalServer,
			)
		}
		return nil, ucerror.New(
			fmt.Errorf("failed to get authorization code: %w", err),
			consts.ErrInternalServer,
		)
	}

	if authzCode.ExpiresAt < time.Now().Unix() {
		return nil, ucerror.New(
			fmt.Errorf("authorization code expired"),
			consts.ErrBadRequest,
		)
	}

	if authzCode.ClientID != clientID {
		return nil, ucerror.New(
			fmt.Errorf("client ID mismatch for authorization code"),
			consts.ErrBadRequest,
		)
	}

	// user might be deleted after the authorization code was issued, so we need to check if the user still exists
	ok, err := uc.postgres.IsUserIDExists(ctx, authzCode.UserID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, ucerror.New(fmt.Errorf("request canceled: %w", err), consts.ErrInternalServer)
		}
		return nil, ucerror.New(
			fmt.Errorf("failed to validate user for authorization code: %w", err),
			consts.ErrInternalServer,
		)
	}

	if !ok {
		return nil, ucerror.New(
			fmt.Errorf("user not found for authorization code"),
			consts.ErrBadRequest,
		)
	}

	return authzCode, nil
}
