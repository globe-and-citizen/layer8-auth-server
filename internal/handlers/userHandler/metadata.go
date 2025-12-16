package userHandler

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) UpdateUserMetadata(c *gin.Context) {
	userID := c.GetUint(consts.MiddlewareKeyUserUserID)

	request, err := utils.DecodeJSONFromRequest[requestdto.UserMetadataUpdate](c)
	if err != nil {
		return
	}

	err = h.uc.UpdateUserMetadata(userID, request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to update user's metadata", err)
		return
	}

	utils.ReturnOK(c, "User's metadata updated successfully", nil)
}
