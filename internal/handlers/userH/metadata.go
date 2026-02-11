package userH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) UpdateMetadata(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserMetadataUpdate](c, h.logger)
	if err != nil {
		return
	}

	ucerr := h.uc.UpdateUserMetadata(c.Request.Context(), userID, request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to update user's metadata!", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "User's metadata updated successfully", nil)
}
