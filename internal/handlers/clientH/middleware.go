package clientH

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) AuthenticateClient(c *gin.Context) {
	token, err := ginUtils.GetBearerToken(c)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	clientID, username, err := h.uc.VerifyClientJWTToken(c.Request.Context(), token)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "Authentication error: invalid token", err)
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
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "authentication failed", fmt.Errorf("authentication failed"))
	}

	//todo update later
	c.Next()
}

func (h ClientHandler) getAuthenticatedUsername(c *gin.Context) (string, error) {
	username := c.GetString(consts.MiddlewareKeyClientUsername)
	if username == "" {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "Failed to get authenticated client username", nil)
		return "", consts.ErrUserUnauthorized
	}

	return username, nil
}

func (h ClientHandler) getAuthenticatedClientID(c *gin.Context) (string, error) {
	clientID := c.GetString(consts.MiddlewareKeyClientClientID)
	if clientID == "" {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "Failed to get authenticated client username", nil)
		return "", consts.ErrUserUnauthorized
	}

	return clientID, nil
}
