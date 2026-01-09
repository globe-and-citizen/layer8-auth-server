package oauthH

import (
	"globe-and-citizen/layer8/auth-server/backend/config"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/oauthUC"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	uc     oauthUC.IOAuthUsecase
	router *gin.RouterGroup
	config config.OAuthConfig
}

func NewOAuthHandler(router *gin.RouterGroup, config config.OAuthConfig, oauthuc oauthUC.IOAuthUsecase) *OAuthHandler {
	return &OAuthHandler{
		uc:     oauthuc,
		router: router,
		config: config,
	}
}

func (h OAuthHandler) RegisterAPIs() {
	h.router.POST("/oauth-login", h.UserLogin)
	h.router.POST("/oauth-login-precheck", h.PrecheckUserLogin)

	oauthGroup := h.router.Group("/oauth")
	oauthGroup.GET("/authorize", h.AuthenticateOAuth, h.AuthorizeContext)
	oauthGroup.POST("/authorize", h.AuthenticateOAuth, h.AuthorizeDecision)

	oauthGroup.POST("/token", h.GetAccessToken)
	oauthGroup.POST("/zk-metadata", h.AuthenticateClient, h.GetZkUserMetadata)
}
