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

func (h UserHandler) RegisterHandler() {

	unauthenticatedGroup := h.router.Group("")
	unauthenticatedGroup.POST("/user-register-precheck", h.PrecheckRegister)
	unauthenticatedGroup.POST("/user-register", h.Register)
	unauthenticatedGroup.POST("/user-login-precheck", h.PrecheckLogin)
	unauthenticatedGroup.POST("/user-login", h.Login)
	unauthenticatedGroup.POST("/user-reset-password-precheck", h.PrecheckResetPassword)
	unauthenticatedGroup.POST("/user-reset-password", h.ResetPassword)

	authenticatedGroup := h.router.Group("user")
	authenticatedGroup.Use(h.AuthenticateUser)

	authenticatedGroup.GET("/profile", h.GetProfile)
	authenticatedGroup.POST("/verify-email", h.VerifyEmail)
	authenticatedGroup.POST("/check-email-verification-code", h.CheckEmailVerificationCode)
	authenticatedGroup.POST("/verify-phone-number-via-bot", h.VerifyPhoneNumber)
	authenticatedGroup.POST("/check-phone-number-verification-code", h.CheckPhoneNumberVerificationCode)
	authenticatedGroup.GET("/generate-telegram-session-id", h.GenerateTelegramSessionID) // this api was originally POST, but I think GET is more suitable
}
