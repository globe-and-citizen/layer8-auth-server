package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) GetAccessToken(c *gin.Context) {
	req, err := ginUtils.DecodeJSONFromRequest[requestdto.OAuthAccessToken](c, h.logger)
	if err != nil {
		return
	}

	response, oauthErr := h.uc.GetAccessToken(c.Request.Context(), req)
	if oauthErr != nil {
		ginUtils.HandleError(c, h.logger, oauthErr.StatusCode, oauthErr.Description, oauthErr.Err)
		return
	}

	ginUtils.ReturnOK(c, "Successful token retrieval", response)
}
