package oauthH

import (
	"globe-and-citizen/layer8/auth-server/backend/config"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/oauthUC"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

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

func (h OAuthHandler) getAuthenticatedUsername(c *gin.Context) (string, error) {
	username := c.GetString(consts.MiddlewareKeyUserUsername)
	if username == "" {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated user username", nil)
		return "", consts.ErrUserUnauthorized
	}

	return username, nil
}

func (h OAuthHandler) getAuthenticatedUserID(c *gin.Context) (uint, error) {
	userID := c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated user ID from context", nil)
		return 0, consts.ErrUserUnauthorized
	}

	return userID, nil
}

func (h OAuthHandler) RegisterHandlers() {
	h.router.POST("/oauth-login", h.UserLogin)
	h.router.POST("/oauth-login-precheck", h.PrecheckUserLogin)

	oauthGroup := h.router.Group("/oauth")
	oauthGroup.GET("/authorize", h.AuthenticateOAuth, h.AuthorizeContext)
	oauthGroup.POST("/authorize", h.AuthenticateOAuth, h.AuthorizeDecision)

	oauthGroup.POST("/token", h.GetAccessToken)
	oauthGroup.POST("/zk-metadata", h.ValidateAccessToken, h.GetZkUserMetadata)
}
