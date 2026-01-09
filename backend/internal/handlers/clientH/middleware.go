package clientH

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) AuthenticateClient(c *gin.Context) {
	token, err := utils.GetBearerToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	clientID, username, err := h.uc.VerifyClientJWTToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyClientUsername, username)
	c.Set(consts.MiddlewareKeyClientClientID, clientID)
	c.Next()
}

func (h ClientHandler) AuthenticateForwardProxy(c *gin.Context) {
	_, _, ok := c.Request.BasicAuth()
	if !ok {
		utils.HandleError(c, http.StatusUnauthorized, "authentication failed", fmt.Errorf("authentication failed"))
	}

	//todo update later
	c.Next()
}

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
