package userH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h UserHandler) VerifyEmail(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserEmailVerify](c, h.logger)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	err = h.uc.VerifyEmail(ctx, userID, request.Email)
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to verify email", err)
		return
	}

	ginUtils.ReturnOK(c, "Verification email sent", nil)
}

func (h UserHandler) CheckEmailVerificationCode(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	request, err := ginUtils.DecodeJSONFromRequest[requestdto.UserCheckEmailVerificationCode](c, h.logger)
	if err != nil {
		return
	}

	ctx := c.Request.Context()
	ucerr := h.uc.CheckEmailVerificationCode(ctx, userID, request.Code)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to verify code", ucerr)
		return
	}

	ucerr = h.uc.SaveProofOfEmailVerification(ctx, userID, request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to verify code", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Your email was successfully verified!", nil)
}
