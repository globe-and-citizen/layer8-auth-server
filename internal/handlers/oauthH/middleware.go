package oauthH

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
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

	userID, username, err := h.uc.VerifyOAuthJWTToken(c.Request.Context(), token)
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

	userID, scopes, err := h.uc.VerifyAccessToken(c.Request.Context(), token)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyOAuthScopes, scopes)
	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h OAuthHandler) getAccessTokenUserID(c *gin.Context) (uint, error) {
	userID := c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "Failed to get authorized userID", fmt.Errorf("failed to get authorized userID"))
		return 0, consts.ErrUserUnauthorized
	}

	return userID, nil
}

func (h OAuthHandler) getAccessTokenScopes(c *gin.Context) (string, error) {
	username := c.GetString(consts.MiddlewareKeyOAuthScopes)
	if username == "" {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "Failed to get authorized scopes", fmt.Errorf("failed to get authorized scopes"))
		return "", consts.ErrUserUnauthorized
	}

	return username, nil
}

func (h OAuthHandler) getAuthenticatedUserID(c *gin.Context) (uint, error) {
	userID := c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "authenticate error", fmt.Errorf("Failed to get authenticated user ID from context"))
		return 0, consts.ErrUserUnauthorized
	}

	return userID, nil
}
