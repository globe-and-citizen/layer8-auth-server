package userH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

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
	ucerr := h.uc.VerifyEmail(ctx, userID, request.Email)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to verify email", ucerr)
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

	ucerr := h.uc.CheckEmailVerificationCode(c.Request.Context(), userID, request)
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to verify code", ucerr)
		return
	}

	ginUtils.ReturnOK(c, "Your email was successfully verified!", nil)
}
