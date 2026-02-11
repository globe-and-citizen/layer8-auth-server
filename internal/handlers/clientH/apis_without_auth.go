package clientH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) PrecheckRegister(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientRegisterPrecheck](c, h.logger)
	if err != nil {
		return
	}

	response, ucerr := h.uc.PrecheckRegister(c.Request.Context(), request, h.config.ScramIterationCount)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to register user", ucerr)
		return
	}

	ginUtils.ReturnCreated(c, "Client is successfully registered", response)
}

func (h ClientHandler) Register(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientRegister](c, h.logger)
	if err != nil {
		return
	}

	ucerr := h.uc.Register(c.Request.Context(), request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to register client", ucerr)
		return
	}

	ginUtils.ReturnCreated(c, "Client registered successfully", nil)
}

func (h ClientHandler) PrecheckLogin(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientLoginPrecheck](c, h.logger)
	if err != nil {
		return
	}

	response, ucerr := h.uc.PrecheckLogin(c.Request.Context(), request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to perform precheck", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Precheck successful", response)
}

func (h ClientHandler) Login(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientLogin](c, h.logger)
	if err != nil {
		return
	}

	response, ucerr := h.uc.Login(c.Request.Context(), request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to perform login", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Login successful", response)
}

func (h ClientHandler) CheckBackendURI(c *gin.Context) {
	request, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientCheckBackendURI](c, h.logger)
	if err != nil {
		return
	}

	response, ucerr := h.uc.CheckBackendURI(c.Request.Context(), request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to check backend url", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Check backend URI successfully", response)
}
