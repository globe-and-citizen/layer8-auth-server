package clientH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) getAuthenticatedUsername(c *gin.Context) (string, error) {
	username := c.GetString(consts.MiddlewareKeyClientUsername)
	if username == "" {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated client username", nil)
		return "", consts.ErrUserUnauthorized
	}

	return username, nil
}

func (h ClientHandler) getAuthenticatedClientID(c *gin.Context) (string, error) {
	clientID := c.GetString(consts.MiddlewareKeyClientClientID)
	if clientID == "" {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated client username", nil)
		return "", consts.ErrUserUnauthorized
	}

	return clientID, nil
}
