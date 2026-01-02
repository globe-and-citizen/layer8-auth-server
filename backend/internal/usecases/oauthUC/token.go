package oauthUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	appError "globe-and-citizen/layer8/auth-server/backend/internal/errors"
	"globe-and-citizen/layer8/auth-server/backend/pkg/oauth"
	"net/http"
)

func (uc *OAuthUsecase) GetAccessToken(req requestdto.OAuthAccessToken) (*responsedto.OAuthAccessToken, *appError.OAuthError) {
	client, err := uc.postgres.GetClientByID(req.ClientID)
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorInvalidClient,
			Description: "incorrect client id or secret",
			StatusCode:  http.StatusUnauthorized,
			Err:         err,
		}
	}

	if client.Secret != req.ClientSecret {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorInvalidClient,
			Description: "incorrect client id or secret",
			StatusCode:  http.StatusUnauthorized,
			Err:         err,
		}
	}

	if client.RedirectURI != req.RedirectURI {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorInvalidRedirectURI,
			Description: "redirect uri mismatch",
			StatusCode:  http.StatusBadRequest,
			Err:         fmt.Errorf("redirect uri mismatch"),
		}
	}

	claims, err := oauth.VerifyAuthorizationCode(req.ClientSecret, req.AuthorizationCode)
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorInvalidAuthzCode,
			Description: "failed to verify authorization code",
			StatusCode:  http.StatusBadRequest,
			Err:         err,
		}
	}

	accessToken, err := uc.token.GenerateOAuthAccessToken(client, *claims)
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts.OAuthErrorServerError,
			Description: "failed to generate access token",
			StatusCode:  http.StatusInternalServerError,
			Err:         err,
		}
	}

	return &responsedto.OAuthAccessToken{
		AccessToken:     accessToken,
		TokenType:       consts.TokenTypeBearer,
		ExpireInMinutes: consts.AccessTokenValidityMinutes,
	}, nil
}
