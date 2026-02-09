package clientH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) PrecheckRegister(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientRegisterPrecheck](c, h.logger)
	if err != nil {
		return
	}

	response, err := h.uc.PrecheckRegister(c.Request.Context(), request, h.config.ScramIterationCount)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	ginUtils.ReturnCreated(c, "Client is successfully registered", response)
}

func (h ClientHandler) Register(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientRegister](c, h.logger)
	if err != nil {
		return
	}

	err = h.uc.Register(c.Request.Context(), request)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to register client", err)
		return
	}

	ginUtils.ReturnCreated(c, "Client registered successfully", nil)
}

func (h ClientHandler) PrecheckLogin(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientLoginPrecheck](c, h.logger)
	if err != nil {
		return
	}

	response, err := h.uc.PrecheckLogin(c.Request.Context(), request)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to perform precheck, service error", err)
		return
	}

	ginUtils.ReturnOK(c, "Precheck successful", response)
}

func (h ClientHandler) Login(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientLogin](c, h.logger)
	if err != nil {
		return
	}

	response, err := h.uc.Login(c.Request.Context(), request)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to perform login", err)
		return
	}

	ginUtils.ReturnOK(c, "Login successful", response)
}

func (h ClientHandler) CheckBackendURI(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientCheckBackendURI](c, h.logger)
	if err != nil {
		return
	}

	response, err := h.uc.CheckBackendURI(c.Request.Context(), request)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to check backend url", err)
		return
	}

	ginUtils.ReturnOK(c, "Check backend URI successfully", response)
}
