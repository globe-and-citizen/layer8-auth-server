package userH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
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
	c.Set(consts.MiddlewareKeyUserUsername, username)
	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h UserHandler) getAuthenticatedUserID(c *gin.Context) (uint, error) {
	userID := c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated user ID from context", nil)
		return 0, consts.ErrUserUnauthorized
	}

	return userID, nil
}
