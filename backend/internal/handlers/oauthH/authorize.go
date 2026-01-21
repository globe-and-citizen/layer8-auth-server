package oauthH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) AuthorizeContext(c *gin.Context) {
	var req requestdto.OAuthAuthorizeContext
	req.ClientID = c.Query("client_id")
	req.Scopes = c.DefaultQuery("scope", string(consts.OAuthScopeReadUser))
	req.RedirectURI = c.Query("redirect_uri")

	response, err := h.uc.AuthorizeContext(req)
	if err != nil {
		utils.HandleError(c, err.StatusCode, err.Description, err.Err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h OAuthHandler) AuthorizeDecision(c *gin.Context) {
	userID, err := h.getAuthenticatedUserID(c)
	if err != nil {
		return
	}

	req, err := utils.DecodeJSONFromRequest[requestdto.OAuthAuthorizeDecision](c)
	if err != nil {
		return
	}

	//var req requestdto.OAuthAuthorizeDecision
	//req.ClientID = c.Query("client_id")
	//req.Scopes = c.DefaultQuery("scope", string(consts.OAuthScopeReadUser))
	//req.RedirectURI = c.Query("redirect_uri")
	//req.ReturnResult = c.DefaultQuery("return_result", "false") == "true"
	//// todo `response_type` is required query param, validate and handle it
	//
	//if c.PostForm("share_display_name") == "true" {
	//	req.Share.DisplayName = true
	//}
	//
	//if c.PostForm("share_color") == "true" {
	//	req.Share.Color = true
	//}
	//
	//if c.PostForm("share_is_email_verified") == "true" {
	//	req.Share.IsEmailVerified = true
	//}
	//
	//if c.PostForm("share_bio") == "true" {
	//	req.Share.Bio = true
	//}

	response, oauthErr := h.uc.AuthorizeDecision(req, userID, h.config.AuthzCodeExpiry)
	if oauthErr != nil {
		//if !req.ReturnResult {
		//	utils.HandleError(c, oauthErr.StatusCode, oauthErr.Description, oauthErr.Err)
		//} else {
		//	c.Redirect(http.StatusSeeOther, "/oauth/error?opt="+string(oauthErr.Code))
		//}
		utils.HandleError(c, oauthErr.StatusCode, oauthErr.Description, oauthErr.Err)
		return
	}

	//if req.ReturnResult {
	//	c.JSON(http.StatusOK, response)
	//	return
	//}
	//
	//c.Redirect(http.StatusSeeOther, response.RedirectURI)
	c.JSON(http.StatusOK, response)
}
