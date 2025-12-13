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

func NewUserHandler(app *gin.Engine, uc userUsecase.IUserUseCase, config config.UserConfig) UserHandler {
	return UserHandler{
		uc:     uc,
		config: config,
		router: app.Group(""),
	}
}

func (h UserHandler) RegisterHandler(middlewares ...gin.HandlerFunc) {
	h.router.POST("users/register-precheck", h.PrecheckRegister)
	h.router.POST("users/register", h.Register)
	h.router.POST("users/login-precheck", h.PrecheckLogin)
	h.router.POST("users/login", h.Login)

	user := h.router.Group("user")
	user.Use(middlewares...)

}
