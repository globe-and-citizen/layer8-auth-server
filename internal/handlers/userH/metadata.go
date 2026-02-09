package userH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

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

	err = h.uc.UpdateUserMetadata(c.Request.Context(), userID, request)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to update user's metadata", err)
		return
	}

	ginUtils.ReturnOK(c, "User's metadata updated successfully", nil)
}
