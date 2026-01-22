package userH

import (
	consts2 "globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) AuthenticateUser(c *gin.Context) {
	token, err := utils.GetBearerToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, username, err := h.uc.VerifyUserJWTToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts2.MiddlewareKeyUserUsername, username)
	c.Set(consts2.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h UserHandler) getAuthenticatedUserID(c *gin.Context) (uint, error) {
	userID := c.GetUint(consts2.MiddlewareKeyUserUserID)

	if userID == 0 {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated user ID from context", nil)
		return 0, consts2.ErrUserUnauthorized
	}

	return userID, nil
}
