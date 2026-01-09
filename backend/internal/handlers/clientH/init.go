package clientH

import (
	"globe-and-citizen/layer8/auth-server/backend/config"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/clientUC"

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
