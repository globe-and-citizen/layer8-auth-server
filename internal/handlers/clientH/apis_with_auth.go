package clientH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) GetProfile(c *gin.Context) {
	username, err := h.getAuthenticatedUsername(c)
	if err != nil {
		return
	}

	profileResp, err := h.uc.GetProfile(c.Request.Context(), username)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "Failed to get user profile, user not found", err)
		return
	}

	ginUtils.ReturnOK(c, "Get client profile successfully", profileResp)
}

func (h ClientHandler) GetUsageStatistics(c *gin.Context) {
	clientID, err := h.getAuthenticatedClientID(c)
	if err != nil {
		return
	}

	response, status, msg, err := h.uc.GetUsageStatistics(c.Request.Context(), clientID)
	if err != nil {
		ginUtils.HandleError(c, h.logger, status, msg, err)
		return
	}

	ginUtils.ReturnOK(c, "Get client usage statistics successful", response)
}

func (h ClientHandler) GetUnpaidAmount(c *gin.Context) {
	clientID, err := h.getAuthenticatedClientID(c)
	if err != nil {
		return
	}

	response, err := h.uc.GetUnpaidAmount(c.Request.Context(), clientID)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "Failed to get unpaid amount", err)
		return
	}

	ginUtils.ReturnOK(c, "successfully retrieved client's unpaid amount", response)
}

func (h ClientHandler) UploadNTorCertificate(c *gin.Context) {
	clientID, err := h.getAuthenticatedClientID(c)
	if err != nil {
		return
	}

	req, err := ginUtils.DecodeJSONFromRequest[requestdto.ClientUploadNTorCertificate](c, h.logger)
	if err != nil {
		return
	}

	err = h.uc.SaveNTorCertificate(c.Request.Context(), clientID, req)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, "failed to save the SP x.509 certificate", err)
		return
	}

	ginUtils.ReturnCreated(c, "x.509 certificate was saved successfully", nil)
}
