package userH

import (
	"globe-and-citizen/layer8/auth-server/backend/config"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/userUC"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc     userUC.IUserUsecase
	config config.UserConfig
	router *gin.RouterGroup
}

func NewUserHandler(router *gin.RouterGroup, uc userUC.IUserUsecase, config config.UserConfig) UserHandler {
	return UserHandler{
		uc:     uc,
		config: config,
		router: router.Group(""),
	}
}

func (h UserHandler) RegisterHandler(authentication gin.HandlerFunc, middlewares ...gin.HandlerFunc) {

	unauthenticatedGroup := h.router.Group("users")
	unauthenticatedGroup.POST("/register-precheck", h.PrecheckRegister)
	unauthenticatedGroup.POST("/register", h.Register)
	unauthenticatedGroup.POST("/login-precheck", h.PrecheckLogin)
	unauthenticatedGroup.POST("/login", h.Login)
	unauthenticatedGroup.POST("/reset-password-precheck", h.PrecheckResetPassword)
	unauthenticatedGroup.POST("/reset-password", h.ResetPassword)

	authenticatedGroup := h.router.Group("user")
	authenticatedGroup.Use(authentication)
	authenticatedGroup.Use(middlewares...)

	authenticatedGroup.GET("/profile", h.GetProfile)
	authenticatedGroup.POST("/verify-email", h.VerifyEmail)
	authenticatedGroup.POST("/check-email-verification-code", h.CheckEmailVerificationCode)
	authenticatedGroup.POST("/verify-phone-number-via-bot", h.VerifyPhoneNumber)
	authenticatedGroup.POST("/check-phone-number-verification-code", h.CheckPhoneNumberVerificationCode)
	authenticatedGroup.GET("/generate-telegram-session-id", h.GenerateTelegramSessionID) // this api was originally POST, but I think GET is more suitable
}
