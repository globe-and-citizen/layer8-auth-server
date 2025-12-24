package clientH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) GetProfile(c *gin.Context) {
	username, err := h.getAuthenticatedUsername(c)
	if err != nil {
		return
	}

	profileResp, err := h.uc.GetProfile(username)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get user profile, user not found", err)
		return
	}

	utils.ReturnOK(c, "Get client profile successfully", profileResp)
}

func (h ClientHandler) GetUsageStatistics(c *gin.Context) {
	clientID, err := h.getAuthenticatedClientID(c)
	if err != nil {
		return
	}

	response, status, msg, err := h.uc.GetUsageStatistics(clientID)
	if err != nil {
		utils.HandleError(c, status, msg, err)
		return
	}

	utils.ReturnOK(c, "Get client usage statistics successful", response)
}

func (h ClientHandler) GetUnpaidAmount(c *gin.Context) {
	clientID, err := h.getAuthenticatedClientID(c)
	if err != nil {
		return
	}

	response, err := h.uc.GetUnpaidAmount(clientID)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to get unpaid amount", err)
		return
	}

	utils.ReturnOK(c, "successfully retrieved client's unpaid amount", response)
}

func (h ClientHandler) UploadNTorCertificate(c *gin.Context) {
	clientID, err := h.getAuthenticatedClientID(c)
	if err != nil {
		return
	}

	req, err := utils.DecodeJSONFromRequest[requestdto.ClientUploadNTorCertificate](c)
	if err != nil {
		return
	}

	err = h.uc.SaveNTorCertificate(clientID, req)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "failed to save the SP x.509 certificate", err)
		return
	}

	utils.ReturnCreated(c, "x.509 certificate was saved successfully", nil)
}
