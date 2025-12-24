package userH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) getAuthenticatedUserID(c *gin.Context) (uint, error) {
	userID := c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated user ID from context", nil)
		return 0, consts.ErrUserUnauthorized
	}

	return userID, nil
}
