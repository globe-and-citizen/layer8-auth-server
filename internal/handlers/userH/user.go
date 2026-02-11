package userH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) PrecheckRegister(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserRegisterPrecheck](c, h.logger)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	response, err := h.uc.PrecheckRegister(ctx, request, h.config.ScramIterationCount)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	ginUtils.ReturnCreated(c, "User is successfully registered", response)
}

func (h UserHandler) Register(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserRegister](c, h.logger)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	err = h.uc.Register(ctx, request)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	ginUtils.ReturnCreated(c, "User registered successfully", nil)
}

func (h UserHandler) PrecheckLogin(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserLoginPrecheck](c, h.logger)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	response, ucerr := h.uc.PrecheckLogin(ctx, request)
	if err != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to precheck", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Precheck successful", response)
}

func (h UserHandler) Login(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserLogin](c, h.logger)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	response, ucerr := h.uc.Login(ctx, request)
	if err != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to login", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Login successful", response)
}

func (h UserHandler) GetProfile(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	profileResp, ucerr := h.uc.GetProfile(c.Request.Context(), userID)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to get profile", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Get user profile successful", profileResp)
}

func (h UserHandler) PrecheckResetPassword(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserResetPasswordPrecheck](c, h.logger)
	if err != nil {
		return
	}

	response, ucerr := h.uc.PrecheckResetPassword(c.Request.Context(), request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "User does not exists!", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "User does exist!", response)
}

func (h UserHandler) ResetPassword(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserResetPassword](c, h.logger)
	if err != nil {
		return
	}

	ucerr := h.uc.ResetPassword(c.Request.Context(), request)
	if err != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to reset password", ucerr)
		return
	}

	ginUtils.ReturnCreated(c, "Your password was updated successfully!", nil)
}
