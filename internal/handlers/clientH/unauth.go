package clientH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) PrecheckRegister(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.ClientRegisterPrecheck](c)
	if err != nil {
		return
	}

	response, err := h.uc.PrecheckRegister(request, h.config.ScramIterationCount)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	utils.ReturnCreated(c, "Client is successfully registered", response)
}

func (h ClientHandler) Register(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.ClientRegister](c)
	if err != nil {
		return
	}

	err = h.uc.Register(request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to register client", err)
		return
	}

	utils.ReturnCreated(c, "Client registered successfully", nil)
}

func (h ClientHandler) PrecheckLogin(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.ClientLoginPrecheck](c)
	if err != nil {
		return
	}

	response, err := h.uc.PrecheckLogin(request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to perform precheck, service error", err)
		return
	}

	utils.ReturnOK(c, "Precheck successful", response)
}

func (h ClientHandler) Login(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.ClientLogin](c)
	if err != nil {
		return
	}

	response, err := h.uc.Login(request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to perform login", err)
		return
	}

	utils.ReturnOK(c, "Login successful", response)
}

func (h ClientHandler) CheckBackendURI(c *gin.Context) {
	request, err := utils.DecodeJSONFromRequest[requestdto.ClientCheckBackendURI](c)
	if err != nil {
		return
	}

	response, err := h.uc.CheckBackendURI(request)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to check backend url", err)
		return
	}

	utils.ReturnOK(c, "Check backend URI successfully", response)
}
