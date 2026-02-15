package userH

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) AuthenticateUser(c *gin.Context) {
	token, err := ginUtils.GetBearerToken(c)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, username, ucerr := h.uc.MdwVerifyUserJWTToken(c.Request.Context(), token)
	if err != nil {
		handlers.HandlerUCError(c, h.logger, "invalid token", ucerr)
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
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "Failed to get authenticated user ID from context", nil)
		return 0, consts.ErrUnauthorized
	}

	return userID, nil
}
