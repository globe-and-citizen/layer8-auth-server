package clientH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) GetProfile(c *gin.Context) {
	username, err := h.getAuthenticatedUsername(c)
	if err != nil {
		return
	}

	profileResp, ucerr := h.uc.GetProfile(c.Request.Context(), username)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to get user profile", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Get client profile successful", profileResp)
}

func (h ClientHandler) GetUsageStatistics(c *gin.Context) {
	clientID, err := h.getAuthenticatedClientID(c)
	if err != nil {
		return
	}

	response, ucerr := h.uc.GetUsageStatistics(c.Request.Context(), clientID)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Get client usage statistics successful", response)
}

func (h ClientHandler) GetUnpaidAmount(c *gin.Context) {
	clientID, err := h.getAuthenticatedClientID(c)
	if err != nil {
		return
	}

	response, ucerr := h.uc.GetUnpaidAmount(c.Request.Context(), clientID)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to get unpaid amount", ucerr)
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

	ucerr := h.uc.SaveNTorCertificate(c.Request.Context(), clientID, req)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "failed to save the SP x.509 certificate", ucerr)
		return
	}

	ginUtils.ReturnCreated(c, "x.509 certificate was saved successfully", nil)
}
