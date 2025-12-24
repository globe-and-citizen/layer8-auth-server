package userH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) VerifyEmail(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	request, err := utils.DecodeJSONFromRequest[requestdto.UserEmailVerify](c)
	if err != nil {
		return
	}

	err = h.uc.VerifyEmail(userID, request.Email)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to verify email", err)
		return
	}

	utils.ReturnOK(c, "Verification email sent", nil)
}

func (h UserHandler) CheckEmailVerificationCode(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	request, err := utils.DecodeJSONFromRequest[requestdto.UserCheckEmailVerificationCode](c)
	if err != nil {
		return
	}

	err = h.uc.CheckEmailVerificationCode(userID, request.Code)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Failed to verify code", err)
		return
	}

	errMsg, err := h.uc.SaveProofOfEmailVerification(userID, request)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, errMsg, err)
	}

	utils.ReturnOK(c, "Your email was successfully verified!", nil)
}
