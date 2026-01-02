package oauthH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) PrecheckUserLogin(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.OAuthUserLoginPrecheck](c)
	if err != nil {
		return
	}

	response, err := h.uc.PrecheckUserLogin(request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to perform precheck, service error", err)
		return
	}

	utils.ReturnOK(c, "Precheck successful", response)
}

func (h OAuthHandler) UserLogin(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.OAuthUserLogin](c)
	if err != nil {
		return
	}

	response, err := h.uc.UserLogin(request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to perform login", err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(consts.OAuthCookieName, response.Token, h.config.CookieMaxAge, "/", "", false, true)
	utils.ReturnOK(c, "Login successful", response.ServerLoginFinalMessage)
}
