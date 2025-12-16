package userHandler

import (
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/usecases/userUsecase"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

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
	unauthorisedGroup.POST("/reset-password-precheck", h.PrecheckResetPassword)
	unauthorisedGroup.POST("/reset-password", h.ResetPassword)

	authorisedGroup := h.router.Group("user")
	authorisedGroup.Use(authorizationMiddleware)
	authorisedGroup.Use(middlewares...)

	authorisedGroup.GET("/profile", h.GetProfile)
	authorisedGroup.POST("/verify-email", h.VerifyEmail)
	authorisedGroup.POST("/check-email-verification-code", h.CheckEmailVerificationCode)
	authorisedGroup.POST("/verify-phone-number-via-bot", h.VerifyPhoneNumber)
	authorisedGroup.POST("/check-phone-number-verification-code", h.CheckPhoneNumberVerificationCode)
	authorisedGroup.GET("/generate-telegram-session-id", h.GenerateTelegramSessionID) // this api was originally POST, but I think GET is more suitable
}

func (h UserHandler) getAuthenticatedUserID(c *gin.Context) (uint, error) {
	userID := c.GetUint(consts.MiddlewareKeyUserUserID)

	if userID == 0 {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get authenticated user ID from context", nil)
		return 0, consts.ErrUserUnauthorized
	}

	return userID, nil
}
