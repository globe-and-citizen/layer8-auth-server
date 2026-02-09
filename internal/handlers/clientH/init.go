package clientH

import (
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/internal/usecases/clientUC"
	"globe-and-citizen/layer8/auth-server/pkg/log"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	logger log.ILogger
	config config.ClientConfig
	uc     clientUC.IClientUsecase
	router *gin.RouterGroup
}

func NewClientHandler(
	logger log.ILogger,
	router *gin.RouterGroup,
	config config.ClientConfig,
	uc clientUC.IClientUsecase,
) ClientHandler {
	return ClientHandler{
		logger: logger,
		config: config,
		uc:     uc,
		router: router,
	}
}

func (h ClientHandler) RegisterAPIs() {
	unauthenticatedGroup := h.router.Group("")
	unauthenticatedGroup.POST("/check-backend-uri", h.CheckBackendURI)
	unauthenticatedGroup.POST("/client-register-precheck", h.PrecheckRegister)
	unauthenticatedGroup.POST("/client-register", h.Register)
	unauthenticatedGroup.POST("/client-login-precheck", h.PrecheckLogin)
	unauthenticatedGroup.POST("/client-login", h.Login)

	authenticatedGroup := h.router.Group("client")
	authenticatedGroup.Use(h.AuthenticateClient)
	authenticatedGroup.GET("/profile", h.GetProfile)
	authenticatedGroup.GET("/usage-stats", h.GetUsageStatistics)
	authenticatedGroup.GET("/unpaid-amount", h.GetUnpaidAmount)
	authenticatedGroup.POST("/upload-certificate", h.UploadNTorCertificate)

	extGroup := h.router.Group("ext")
	extGroup.GET("/client-cert", h.AuthenticateForwardProxy, h.GetNTorCertificate)
}
