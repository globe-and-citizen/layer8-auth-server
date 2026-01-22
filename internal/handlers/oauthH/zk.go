package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) GetZkUserMetadata(c *gin.Context) {
	userID, err := h.getAccessTokenUserID(c)
	if err != nil {
		return
	}

	scopes, err := h.getAccessTokenScopes(c)
	if err != nil {
		return
	}

	req := requestdto.OAuthZkMetadata{
		UserID: userID,
		Scopes: scopes,
	}

	zkMetadata, oauthErr := h.uc.GetZkUserMetadata(req)
	if oauthErr != nil {
		utils.HandleError(c, oauthErr.StatusCode, oauthErr.Description, oauthErr.Err)
		return
	}

	utils.ReturnOK(c, "User metadata retrieved successfully", zkMetadata)
}
