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
	"globe-and-citizen/layer8/auth-server/pkg/oauth"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (uc *OAuthUsecase) GetAuthorizeContext(ctx context.Context, req requestdto.OAuthAuthorizeContext) (*responsedto.OAuthAuthorizeContext, *ucerror.UCError) {
	client, scopes, err := uc.validateAuthorizeParams(ctx, req)
	if err != nil {
		return nil, err
	}

	return &responsedto.OAuthAuthorizeContext{
		ClientName: client.Name,
		Scopes:     uc.getAuthorizeScopes(scopes),
	}, nil
}

func (uc *OAuthUsecase) PostAuthorizeDecision(
	ctx context.Context,
	req requestdto.OAuthAuthorizeDecision,
	userID uint,
	authzCodeExpiry time.Duration,
) (*responsedto.OAuthAuthorizeDecision, *ucerror.UCError) {
	client, scopes, ucErr := uc.validateAuthorizeParams(ctx, req.OAuthAuthorizeContext)
	if ucErr != nil {
		return nil, ucErr
	}

	scopes = append([]OAuthScope{}, scopes...)

	if req.Share.DisplayName {
		scopes = append(scopes, ScopeReadUserDisplayName)
	}
	if req.Share.Color {
		scopes = append(scopes, ScopeReadUserColor)
	}
	if req.Share.Bio {
		scopes = append(scopes, ScopeReadUserBio)
	}
	if req.Share.IsEmailVerified {
		scopes = append(scopes, ScopeReadUserIsEmailVerified)
	}

	code, err := oauth.GenerateAuthorizationCode(req.ClientID, uc.config.AuthzCodeSecret,
		client.RedirectURI, ScopesToStringSlice(scopes), userID, authzCodeExpiry)
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("error generating authorization code: %w", err), consts.ErrInternalServer)
	}

	redirectURL, err := oauth.GenerateAuthURL(req.ClientID, code, client.RedirectURI, ScopesToStringSlice(scopes))
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("error generating authURL: %w", err), consts.ErrInternalServer)
	}

	return &responsedto.OAuthAuthorizeDecision{
		RedirectURI: redirectURL,
		Code:        code,
	}, nil
}

// Normative (MUST / REQUIRED):
//   - MUST support HTTP GET.
//   - MAY support POST.
//   - MUST validate required parameters:
//     response_type
//     client_id
//     redirect_uri (if multiple registered)
//     scope
//   - MUST require "openid" scope for OIDC requests.
//   - MUST validate redirect_uri via exact string match.
//   - MUST return errors using redirect-based error response format.
//   - MUST echo `state` exactly if provided.
//   - MUST require and bind `nonce` for OIDC flows issuing ID tokens.
//   - MUST validate PKCE (code_challenge) if present.
func (uc *OAuthUsecase) validateAuthorizeParams(ctx context.Context, req requestdto.OAuthAuthorizeContext) (*gormModels.Client, []OAuthScope, *ucerror.UCError) {
	// 1. Validate response_type
	if req.ResponseType != oauth.ResponseTypeCode {
		return nil, nil, ucerror.New(fmt.Errorf("invalid response type: %s", req.ResponseType), consts.ErrBadRequest)
	}

	// 2. Validate client_id and redirect_uri
	client, err := uc.postgres.GetClientByID(ctx, req.ClientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil,
				ucerror.New(fmt.Errorf("clientID:%s not found: %w", req.ClientID, err), consts.ErrBadRequest)
		}
		return nil, nil,
			ucerror.New(fmt.Errorf("unable to get client by ID:%s: %w", req.ClientID, err), consts.ErrInternalServer)
	}

	if req.RedirectURI != "" && req.RedirectURI != client.RedirectURI {
		return nil, nil,
			ucerror.New(
				fmt.Errorf("%s does not match registered URI:%s", req.ClientID, client.RedirectURI),
				consts.ErrBadRequest,
			)
	}

	// 3. Validate scopes
	var scopes []OAuthScope
	requestedScopes := strings.Split(req.Scopes, ",")
	for _, scope := range requestedScopes {
		if !OAuthScope(scope).IsValid() {
			return nil, nil, ucerror.New(fmt.Errorf("invalid scope requested: %s", scope), consts.ErrBadRequest)
		}
		scopes = append(scopes, OAuthScope(scope))
	}

	if req.Scopes == "" {
		scopes = append(scopes, ScopeReadUser)
	}

	return client, scopes, nil
}

func (uc *OAuthUsecase) getAuthorizeScopes(scopes []OAuthScope) []responsedto.OAuthAuthorizeScopes {
	var scopesDesc []responsedto.OAuthAuthorizeScopes
	for _, s := range scopes {
		scopesDesc = append(scopesDesc, responsedto.OAuthAuthorizeScopes{
			Name:        string(s),
			Description: ScopeDescriptions[s],
		})
	}
	return scopesDesc
}
