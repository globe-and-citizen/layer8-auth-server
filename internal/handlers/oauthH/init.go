package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/internal/usecases/oauthUC"
	"globe-and-citizen/layer8/auth-server/pkg/log"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	logger log.ILogger
	uc     oauthUC.IOAuthUsecase
	router *gin.RouterGroup
	config config.OAuthConfig
}

func NewOAuthHandler(
	logger log.ILogger,
	router *gin.RouterGroup,
	config config.OAuthConfig,
	oauthuc oauthUC.IOAuthUsecase,
) *OAuthHandler {
	return &OAuthHandler{
		logger: logger,
		uc:     oauthuc,
		router: router,
		config: config,
	}
}

func (h OAuthHandler) RegisterAPIs() {
	h.router.POST("/oauth-login-precheck", h.PrecheckUserLogin)
	h.router.POST("/oauth-login", h.UserLogin)

	oauthGroup := h.router.Group("/oauth")
	oauthGroup.GET("/authorize", h.AuthenticateOAuth, h.GetAuthorizeContext)
	oauthGroup.POST("/authorize", h.AuthenticateOAuth, h.PostAuthorizeDecision)

	oauthGroup.POST("/token", h.GetAccessToken)
	oauthGroup.POST("/zk-metadata", h.AuthenticateClient, h.GetZkUserMetadata)
}
