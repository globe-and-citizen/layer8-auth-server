package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) GetAccessToken(c *gin.Context) {
	req, err := ginUtils.DecodeJSONFromRequest[requestdto.OAuthAccessToken](c, h.logger)
	if err != nil {
		return
	}

	response, ucErr := h.uc.GetAccessToken(c.Request.Context(), req)
	if ucErr != nil {
		handlers.HandlerUCError(c, h.logger, "", ucErr)
		return
	}

	ginUtils.ReturnOK(c, "Successful token retrieval", response)
}
