package clientHandler

import (
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/usecases/clientUsecase"

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

}
