package clientH

import (
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/usecases/clientUC"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	config config.ClientConfig
	uc     clientUC.IClientUsecase
	router *gin.RouterGroup
}

func NewClientHandler(
	router *gin.RouterGroup,
	config config.ClientConfig,
	uc clientUC.IClientUsecase,
) ClientHandler {
	return ClientHandler{
		config: config,
		uc:     uc,
		router: router,
	}
}

func (h ClientHandler) RegisterHandler(authentication gin.HandlerFunc, middlewares ...gin.HandlerFunc) {
	unauthenticatedGroup := h.router.Group("clients")
	unauthenticatedGroup.POST("/check-backend-uri", h.CheckBackendURI)
	unauthenticatedGroup.POST("/register-precheck", h.PrecheckRegister)
	unauthenticatedGroup.POST("/register", h.Register)
	unauthenticatedGroup.POST("/login-precheck", h.PrecheckLogin)
	unauthenticatedGroup.POST("/login", h.Login)

	authenticatedGroup := h.router.Group("client")
	authenticatedGroup.Use(authentication)
	authenticatedGroup.Use(middlewares...)
	authenticatedGroup.POST("/profile", h.GetProfile)
	authenticatedGroup.GET("/usage-stats", h.GetUsageStatistics)
	authenticatedGroup.GET("/client-unpaid-amount", h.GetUnpaidAmount)
}
