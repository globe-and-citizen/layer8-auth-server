package oauthH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) GetAccessToken(c *gin.Context) {
	req, err := utils.DecodeJSONFromRequest[requestdto.OAuthAccessToken](c)
	if err != nil {
		return
	}

	response, oauthErr := h.uc.GetAccessToken(req)
	if oauthErr != nil {
		utils.HandleError(c, oauthErr.StatusCode, oauthErr.Description, oauthErr.Err)
		return
	}

	utils.ReturnOK(c, "Successful token retrieval", response)
}
