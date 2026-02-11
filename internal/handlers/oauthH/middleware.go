package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) AuthenticateOAuth(c *gin.Context) {
	token, err := c.Cookie(consts.OAuthCookieName)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, username, err := h.uc.MdwVerifyUserLoggedInToken(c.Request.Context(), token)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyUserUsername, username)
	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h OAuthHandler) AuthenticateClient(c *gin.Context) {
	token, err := ginUtils.GetBearerToken(c)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, scopes, err := h.uc.MdwVerifyClientAccessToken(c.Request.Context(), token)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyOAuthScopes, scopes)
	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h OAuthHandler) getAccessTokenUserID(c *gin.Context) (userID uint, isError bool) {
	userID = c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		handlers.HandlerUCError(c, h.logger, "failed to get oauth access token userID from context", nil)
		return 0, true
	}

	return userID, false
}

func (h OAuthHandler) getAccessTokenScopes(c *gin.Context) (scopes string, isError bool) {
	scopes = c.GetString(consts.MiddlewareKeyOAuthScopes)
	if scopes == "" {
		handlers.HandlerUCError(c, h.logger, "failed to get oauth access token scopes from context", nil)
		return "", true
	}

	return scopes, false
}

func (h OAuthHandler) getAuthenticatedUserID(c *gin.Context) (userID uint, isError bool) {
	userID = c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		handlers.HandlerUCError(c, h.logger, "failed to get oauth authenticated userID from context", nil)
		return 0, true
	}

	return userID, false
}
