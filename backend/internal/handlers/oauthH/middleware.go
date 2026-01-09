package oauthH

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) AuthenticateOAuth(c *gin.Context) {
	token, err := c.Cookie(consts.OAuthCookieName)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, username, err := h.uc.VerifyOAuthJWTToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyUserUsername, username)
	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h OAuthHandler) AuthenticateClient(c *gin.Context) {
	token, err := utils.GetBearerToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, scopes, err := h.uc.VerifyAccessToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
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
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authorized userID", fmt.Errorf("failed to get authorized userID"))
		return 0, consts.ErrUserUnauthorized
	}

	return userID, nil
}

func (h OAuthHandler) getAccessTokenScopes(c *gin.Context) (string, error) {
	username := c.GetString(consts.MiddlewareKeyOAuthScopes)
	if username == "" {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authorized scopes", fmt.Errorf("failed to get authorized scopes"))
		return "", consts.ErrUserUnauthorized
	}

	return username, nil
}

func (h OAuthHandler) getAuthenticatedUserID(c *gin.Context) (uint, error) {
	userID := c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		utils.HandleError(c, http.StatusInternalServerError, "authenticate error", fmt.Errorf("Failed to get authenticated user ID from context"))
		return 0, consts.ErrUserUnauthorized
	}

	return userID, nil
}
