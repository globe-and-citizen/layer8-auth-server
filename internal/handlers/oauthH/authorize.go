package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"globe-and-citizen/layer8/auth-server/pkg/oauth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAuthorizeContext handles the OAuth2 / OIDC Authorization Endpoint.
//
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
//
// Conditionally Normative:
//   - PKCE REQUIRED for public clients.
//   - Consent screen REQUIRED unless prior consent exists.
//
// Non-Normative:
//   - Endpoint path (e.g., /authorize, /oauth2/auth).
//   - UI design and login screen behavior.
//   - Internal session management.
//   - Storage implementation.
//
// TODO...
func (h OAuthHandler) GetAuthorizeContext(c *gin.Context) {
	var req requestdto.OAuthAuthorizeContext
	req.ResponseType = c.Query(oauth.ParamResponseType)
	req.ClientID = c.Query(oauth.ParamClientID)
	req.Scopes = c.Query(oauth.ParamScope)
	req.RedirectURI = c.Query(oauth.ParamRedirectURI)
	req.State = c.Query(oauth.ParamState)

	response, err := h.uc.GetAuthorizeContext(c.Request.Context(), req)
	if err != nil {
		handlers.HandlerUCError(c, h.logger, "", err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// PostAuthorizeDecision handles the OAuth2 / OIDC Authorization Endpoint.
//
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
//
// Conditionally Normative:
//   - PKCE REQUIRED for public clients.
//   - Consent screen REQUIRED unless prior consent exists.
//
// Non-Normative:
//   - Endpoint path (e.g., /authorize, /oauth2/auth).
//   - UI design and login screen behavior.
//   - Internal session management.
//   - Storage implementation.
//
// TODO...
func (h OAuthHandler) PostAuthorizeDecision(c *gin.Context) {
	userID, isErr := h.getAuthenticatedUserID(c)
	if isErr {
		return
	}

	req, err := ginUtils.DecodeJSONFromRequest[requestdto.OAuthAuthorizeDecision](c, h.logger)
	if err != nil {
		return
	}

	//var req requestdto.OAuthAuthorizeDecision
	req.ResponseType = c.Query(oauth.ParamResponseType)
	req.ClientID = c.Query(oauth.ParamClientID)
	req.Scopes = c.Query(oauth.ParamScope)
	req.RedirectURI = c.Query(oauth.ParamRedirectURI)
	req.State = c.Query(oauth.ParamState)
	req.ReturnResult = c.DefaultQuery("return_result", "false") == "true"

	response, ucErr := h.uc.PostAuthorizeDecision(c.Request.Context(), req, userID)
	if ucErr != nil {
		handlers.HandlerUCError(c, h.logger, "", ucErr)
		return
	}

	c.JSON(http.StatusOK, response)
}
