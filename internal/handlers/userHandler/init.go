package userHandler

import (
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/usecases/userUsecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc     userUsecase.IUserUseCase
	config config.UserConfig
	router *gin.RouterGroup
}

func NewUserHandler(router *gin.RouterGroup, uc userUsecase.IUserUseCase, config config.UserConfig) UserHandler {
	return UserHandler{
		uc:     uc,
		config: config,
		router: router.Group(""),
	}
}

func (h UserHandler) RegisterHandler(authorizationMiddleware gin.HandlerFunc, middlewares ...gin.HandlerFunc) {

	unauthorisedGroup := h.router.Group("users")
	unauthorisedGroup.POST("/register-precheck", h.PrecheckRegister)
	unauthorisedGroup.POST("/register", h.Register)
	unauthorisedGroup.POST("/login-precheck", h.PrecheckLogin)
	unauthorisedGroup.POST("/login", h.Login)

	authorisedGroup := h.router.Group("user")
	authorisedGroup.Use(authorizationMiddleware)
	authorisedGroup.Use(middlewares...)

	authorisedGroup.POST("/verify-email", h.VerifyEmail)
	authorisedGroup.POST("/check-email-verification-code", h.CheckEmailVerificationCode)
}
