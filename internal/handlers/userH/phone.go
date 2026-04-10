package userH

import (
	"encoding/base64"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) VerifyPhoneNumber(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	ucerr := h.uc.VerifyPhoneNumber(c.Request.Context(), userID)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Verify phone number failed", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Phone number has been verified!", nil)
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

	ucerr := h.uc.CheckPhoneNumberVerificationCode(c.Request.Context(), userID, request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Check phone number verification code failed", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Your phone number is verified successfully! Congratulations!", nil)
}

func (h UserHandler) GenerateTelegramSessionID(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	sessionID, ucerr := h.uc.GenerateAndSaveTelegramSessionIDHash(c.Request.Context(), userID)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Generate telegram session ID failed", ucerr)
		return
	}

	sessionIdDTO := responsedto.UserGetTelegramSessionID{
		SessionID: base64.RawURLEncoding.EncodeToString(sessionID),
	}

	ginUtils.ReturnOK(c, "session id generated", sessionIdDTO)
}
