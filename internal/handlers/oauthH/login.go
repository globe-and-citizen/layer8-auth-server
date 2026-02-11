package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) PrecheckUserLogin(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.OAuthUserLoginPrecheck](c, h.logger)
	if err != nil {
		return
	}

	response, err := h.uc.PrecheckUserLogin(c.Request.Context(), request)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to perform precheck, service error", err)
		return
	}

	ginUtils.ReturnOK(c, "Precheck successful", response)
}

func (h OAuthHandler) UserLogin(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.OAuthUserLogin](c, h.logger)
	if err != nil {
		return
	}

	response, ucErr := h.uc.UserLogin(c.Request.Context(), request)
	if ucErr != nil {
		handlers.HandlerUCError(c, h.logger, "", ucErr)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(consts.OAuthCookieName, response.Token, h.config.CookieMaxAge, "/", "", false, true)
	ginUtils.ReturnOK(c, "Login successful", response.ServerLoginFinalMessage)
}
