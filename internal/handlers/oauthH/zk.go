package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) GetZkUserMetadata(c *gin.Context) {
	userID, isErr := h.getAccessTokenUserID(c)
	if isErr {
		return
	}

	scopes, isErr := h.getAccessTokenScopes(c)
	if isErr {
		return
	}

	req := requestdto.OAuthZkMetadata{
		UserID: userID,
		Scopes: scopes,
	}

	zkMetadata, ucErr := h.uc.GetZkUserMetadata(c.Request.Context(), req)
	if ucErr != nil {
		handlers.HandlerUCError(c, h.logger, "", ucErr)
		return
	}

	ginUtils.ReturnOK(c, "User metadata retrieved successfully", zkMetadata)
}
