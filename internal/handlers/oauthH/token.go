package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

	"github.com/gin-gonic/gin"
)

// GetAccessToken handles the OAuth2 Token Endpoint.
//
// Normative (MUST / REQUIRED):
//   - MUST use HTTP POST.
//   - MUST require Content-Type: application/x-www-form-urlencoded.
//   - MUST require grant_type.
//   - MUST validate grant_type-specific parameters.
//   - MUST authenticate confidential clients.
//   - MUST validate authorization code:
//   - single use
//   - not expired
//   - bound to client_id
//   - bound to redirect_uri
//   - MUST validate PKCE code_verifier if used.
//   - MUST return JSON response.
//   - MUST include access_token and token_type ("Bearer").
//   - MUST issue id_token if scope includes "openid".
//
// Conditionally Normative:
//   - redirect_uri required if present in authorization request.
//   - refresh_token issuance depends on policy.
//
// Non-Normative:
//   - Endpoint path (e.g., /token).
//   - Access token format (opaque vs JWT).
//   - Token lifetime policy.
//
// TODO...
func (h OAuthHandler) GetAccessToken(c *gin.Context) {
	req, err := ginUtils.DecodeJSONFromRequest[requestdto.OAuthAccessToken](c, h.logger) // fixme application/x-www-form-urlencoded.
	if err != nil {
		return
	}

	response, ucErr := h.uc.GetAccessToken(c.Request.Context(), req)
	if ucErr != nil {
		handlers.HandlerUCError(c, h.logger, "", ucErr)
		return
	}

	ginUtils.ReturnOK(c, "Successful token retrieval", response)
}
