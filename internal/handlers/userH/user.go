package userH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) PrecheckRegister(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.UserRegisterPrecheck](c)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	response, err := h.uc.PrecheckRegister(ctx, request, h.config.ScramIterationCount)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	utils.ReturnCreated(c, "User is successfully registered", response)
}

func (h UserHandler) Register(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.UserRegister](c)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	err = h.uc.Register(ctx, request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	utils.ReturnCreated(c, "User registered successfully", nil)
}

func (h UserHandler) PrecheckLogin(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.UserLoginPrecheck](c)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	response, err := h.uc.PrecheckLogin(ctx, request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to perform precheck, service error", err)
		return
	}

	utils.ReturnOK(c, "Precheck successful", response)
}

func (h UserHandler) Login(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.UserLogin](c)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	response, err := h.uc.Login(ctx, request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to perform login", err)
		return
	}

	utils.ReturnOK(c, "Login successful", response)
}

func (h UserHandler) GetProfile(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	profileResp, err := h.uc.GetProfile(ctx, userID)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to get user profile", err)
		return
	}

	utils.ReturnOK(c, "Get user profile successful", profileResp)
}

func (h UserHandler) PrecheckResetPassword(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.UserResetPasswordPrecheck](c)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	response, err := h.uc.PrecheckResetPassword(ctx, request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "User does not exist!", err)
		return
	}

	utils.ReturnOK(c, "User does exist!", response)
}

func (h UserHandler) ResetPassword(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.UserResetPassword](c)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	status, msg, err := h.uc.ResetPassword(ctx, request)
	if err != nil {
		utils.HandleError(c, status, msg, err)
		return
	}

	utils.ReturnCreated(c, "Your password was updated successfully!", nil)
}
