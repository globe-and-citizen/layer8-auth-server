package userH

import (
	"encoding/base64"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) VerifyPhoneNumber(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	message, err := h.uc.VerifyPhoneNumber(c.Request.Context(), userID)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, message, err)
		return
	}

	ginUtils.ReturnOK(c, message, nil)
}

func (h UserHandler) CheckPhoneNumberVerificationCode(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserCheckPhoneNumberVerificationCode](c, h.logger)
	if err != nil {
		return
	}

	status, message, err := h.uc.CheckPhoneNumberVerificationCode(c.Request.Context(), userID, request)
	if err != nil {
		ginUtils.HandleError(c, h.logger, status, message, err)
		return
	}

	ginUtils.ReturnOK(c, "Your phone number is verified successfully! Congratulations!", nil)
}

func (h UserHandler) GenerateTelegramSessionID(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	sessionID, errMsg, err := h.uc.GenerateAndSaveTelegramSessionIDHash(c.Request.Context(), userID)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusInternalServerError, errMsg, err)
		return
	}

	sessionIdDTO := responsedto.UserGetTelegramSessionID{
		SessionID: base64.RawURLEncoding.EncodeToString(sessionID),
	}

	ginUtils.ReturnOK(c, "session id generated", sessionIdDTO)
}
