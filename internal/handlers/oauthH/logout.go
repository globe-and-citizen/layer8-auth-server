package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) UserLogout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		consts.OAuthCookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	ginUtils.ReturnOK(c, "Logout successful", nil)
}
