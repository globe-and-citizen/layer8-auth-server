package oauthUC

import (
	"fmt"
	consts2 "globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	appError "globe-and-citizen/layer8/auth-server/internal/errors"
	"globe-and-citizen/layer8/auth-server/pkg/oauth"
	"net/http"
)

func (uc *OAuthUsecase) GetAccessToken(req requestdto.OAuthAccessToken) (*responsedto.OAuthAccessToken, *appError.OAuthError) {
	client, err := uc.postgres.GetClientByID(req.ClientID)
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts2.OAuthErrorInvalidClient,
			Description: "incorrect client id or secret",
			StatusCode:  http.StatusUnauthorized,
			Err:         err,
		}
	}

	if client.Secret != req.ClientSecret {
		return nil, &appError.OAuthError{
			Code:        consts2.OAuthErrorInvalidClient,
			Description: "incorrect client id or secret",
			StatusCode:  http.StatusUnauthorized,
			Err:         err,
		}
	}

	if client.RedirectURI != req.RedirectURI {
		return nil, &appError.OAuthError{
			Code:        consts2.OAuthErrorInvalidRedirectURI,
			Description: "redirect uri mismatch",
			StatusCode:  http.StatusBadRequest,
			Err:         fmt.Errorf("redirect uri mismatch"),
		}
	}

	claims, err := oauth.VerifyAuthorizationCode(req.ClientSecret, req.AuthorizationCode)
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts2.OAuthErrorInvalidAuthzCode,
			Description: "failed to verify authorization code",
			StatusCode:  http.StatusBadRequest,
			Err:         err,
		}
	}

	accessToken, err := uc.token.GenerateOAuthAccessToken(client, *claims)
	if err != nil {
		return nil, &appError.OAuthError{
			Code:        consts2.OAuthErrorServerError,
			Description: "failed to generate access token",
			StatusCode:  http.StatusInternalServerError,
			Err:         err,
		}
	}

	return &responsedto.OAuthAccessToken{
		AccessToken:     accessToken,
		TokenType:       consts2.TokenTypeBearer,
		ExpireInMinutes: consts2.AccessTokenValidityMinutes,
	}, nil
}
