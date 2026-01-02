package oauthH

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
