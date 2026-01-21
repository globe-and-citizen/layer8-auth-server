package oauthUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	appError "globe-and-citizen/layer8/auth-server/backend/internal/errors"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/backend/pkg/oauth"
	"net/http"
	"strings"
	"time"
)

func (uc *OAuthUsecase) AuthorizeContext(req requestdto.OAuthAuthorizeContext) (*responsedto.OAuthAuthorizeContext, *appError.OAuthError) {
	client, _, scopes, err := uc.validateAuthorizeParams(req)
	if err != nil {
		return nil, err
	}

	return &responsedto.OAuthAuthorizeContext{
		ClientName: client.Name,
		Scopes:     uc.getAuthorizeScopes(scopes),
	}, nil
}

func (uc *OAuthUsecase) AuthorizeDecision(
	req requestdto.OAuthAuthorizeDecision,
	userID uint,
	authzCodeExpiry time.Duration,
) (*responsedto.OAuthAuthorizeDecision, *appError.OAuthError) {
	client, _, scopes, oauthErr := uc.validateAuthorizeParams(req.OAuthAuthorizeContext)
	if oauthErr != nil {
		return nil, oauthErr
	}

	scopes = append([]consts.OAuthScope{}, scopes...)

	if req.Share.DisplayName {
		scopes = append(scopes, consts.OAuthScopeReadUserDisplayName)
	}
	if req.Share.Color {
		scopes = append(scopes, consts.OAuthScopeReadUserColor)
	}
	if req.Share.Bio {
		scopes = append(scopes, consts.OAuthScopeReadUserBio)
	}
	if req.Share.IsEmailVerified {
		scopes = append(scopes, consts.OAuthScopeReadUserIsEmailVerified)
	}

	code, err := oauth.GenerateAuthorizationCode(req.ClientID, client.Secret,
		client.RedirectURI, consts.OAuthScopesToStringSlice(scopes), userID, authzCodeExpiry)
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorServerError,
			StatusCode:  http.StatusInternalServerError,
			Description: "Failed to generate authorization code",
			Err:         err,
		}
	}

	redirectURL, err := oauth.GenerateAuthURL(req.ClientID, code, client.RedirectURI, consts.OAuthScopesToStringSlice(scopes))
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorServerError,
			StatusCode:  http.StatusInternalServerError,
			Description: "Failed to generate redirect URL",
			Err:         err,
		}
	}

	return &responsedto.OAuthAuthorizeDecision{
		RedirectURI: redirectURL,
		Code:        code,
	}, nil
}

func (uc *OAuthUsecase) validateAuthorizeParams(req requestdto.OAuthAuthorizeContext) (*gormModels.Client, string, []consts.OAuthScope, *appError.OAuthError) {
	client, err := uc.postgres.GetClientByID(req.ClientID) // todo remember to filter errors appropriately
	if err != nil {
		return nil, "", nil, &appError.OAuthError{
			Code:       consts.OAuthErrorInvalidClient,
			StatusCode: 400,
			Err:        fmt.Errorf("client with ID:%s not found: %v", req.ClientID, err),
		}
	}

	if req.RedirectURI != "" && req.RedirectURI != client.RedirectURI {
		return nil, "", nil, &appError.OAuthError{
			Code:       consts.OAuthErrorInvalidRedirectURI,
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("redirect_uri does not match registered URI"),
		}
	}

	var scopes []consts.OAuthScope
	requestedScopes := strings.Split(req.Scopes, ",")
	for _, scope := range requestedScopes {
		if !consts.OAuthScope(scope).IsValid() {
			return nil, "", nil, &appError.OAuthError{
				Code:       consts.OAuthErrorInvalidScope,
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("invalid scope requested: %s", scope),
			}
		}
		scopes = append(scopes, consts.OAuthScope(scope))
	}

	if req.Scopes == "" {
		scopes = append(scopes, consts.OAuthScopeReadUser)
	}

	return &client, req.RedirectURI, scopes, nil
}

func (uc *OAuthUsecase) getAuthorizeScopes(scopes []consts.OAuthScope) []responsedto.OAuthAuthorizeScopes {
	var scopesDesc []responsedto.OAuthAuthorizeScopes
	for _, s := range scopes {
		scopesDesc = append(scopesDesc, responsedto.OAuthAuthorizeScopes{
			Name:        string(s),
			Description: consts.ScopeDescriptions[consts.OAuthScope(s)],
		})
	}
	return scopesDesc
}
