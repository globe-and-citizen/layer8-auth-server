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

	// POST /token
	// Content-Type: application/x-www-form-urlencoded
	//
	// grant_type=authorization_code&
	// code=abc123&
	// redirect_uri=https://client.com/callback&
	// client_id=client1&
	// client_secret=secret
	oauthGroup.POST("/token", h.RequestToken)
	oauthGroup.POST("/zk-metadata", h.AuthenticateClient, h.GetZkUserMetadata)
}
