package clientHandler

import (
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/usecases/clientUsecase"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	config config.ClientConfig
	uc     clientUsecase.IClientUsecase
	router *gin.RouterGroup
}

func NewClientHandler(
	router *gin.RouterGroup,
	config config.ClientConfig,
	uc clientUsecase.IClientUsecase,
) ClientHandler {
	return ClientHandler{
		config: config,
		uc:     uc,
		router: router,
	}
}

func (h ClientHandler) RegisterHandler(middlewares ...gin.HandlerFunc) {
	unauthorisedGroup := h.router.Group("clients")
	unauthorisedGroup.POST("/check-backend-uri", h.CheckBackendURI)
	unauthorisedGroup.POST("/register-precheck", h.PrecheckRegister)
	unauthorisedGroup.POST("/register", h.Register)
	unauthorisedGroup.POST("/login-precheck", h.PrecheckLogin)
	unauthorisedGroup.POST("/login", h.Login)

	authorisedGroup := h.router.Group("client")
	authorisedGroup.Use(middlewares...)
	authorisedGroup.POST("/profile", h.GetProfile)
	authorisedGroup.GET("/usage-stats", h.GetUsageStatistics)
	authorisedGroup.GET("/client-unpaid-amount", h.GetUnpaidAmount)
}

func (h ClientHandler) getAuthenticatedUsername(c *gin.Context) (string, error) {
	username := c.GetString(consts.MiddlewareKeyClientUsername)
	if username == "" {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated client username", nil)
		return "", consts.ErrUserUnauthorized
	}

	return username, nil
}

func (h ClientHandler) getAuthenticatedClientID(c *gin.Context) (string, error) {
	clientID := c.GetString(consts.MiddlewareKeyClientClientID)
	if clientID == "" {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated client username", nil)
		return "", consts.ErrUserUnauthorized
	}

	return clientID, nil
}
