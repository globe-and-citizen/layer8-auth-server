package userHandler

import (
	"encoding/base64"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) VerifyPhoneNumber(c *gin.Context) {
	userID := c.GetUint(consts.MiddlewareKeyUserID)

	message, err := h.uc.VerifyPhoneNumber(userID)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, message, err)
		return
	}

	utils.ReturnOK(c, message, nil)
}

func (h UserHandler) CheckPhoneNumberVerificationCode(c *gin.Context) {
	userID := c.GetUint(consts.MiddlewareKeyUserID)

	request, err := utils.DecodeJSONFromRequest[requestdto.UserCheckPhoneNumberVerificationCode](c)
	if err != nil {
		return
	}

	status, message, err := h.uc.CheckPhoneNumberVerificationCode(userID, request)
	if err != nil {
		utils.HandleError(c, status, message, err)
		return
	}

	utils.ReturnOK(c, "Your phone number is verified successfully! Congratulations!", nil)
}

func (h UserHandler) GenerateTelegramSessionID(c *gin.Context) {
	userID := c.GetUint(consts.MiddlewareKeyUserID)

	sessionID, errMsg, err := h.uc.GenerateAndSaveTelegramSessionIDHash(userID)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, errMsg, err)
		return
	}

	sessionIdDTO := responsedto.UserGetTelegramSessionID{
		SessionID: base64.RawURLEncoding.EncodeToString(sessionID),
	}

	utils.ReturnOK(c, "session id generated", sessionIdDTO)
}
