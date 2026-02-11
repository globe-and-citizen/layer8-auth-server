package oauthH

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h OAuthHandler) GetAuthorizeContext(c *gin.Context) {
	var req requestdto.OAuthAuthorizeContext
	req.ClientID = c.Query("client_id")
	req.Scopes = c.DefaultQuery("scope", string(consts.OAuthScopeReadUser))
	req.RedirectURI = c.Query("redirect_uri")

	response, err := h.uc.GetAuthorizeContext(c.Request.Context(), req)
	if err != nil {
		handlers.HandlerUCError(c, h.logger, "", err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h OAuthHandler) PostAuthorizeDecision(c *gin.Context) {
	userID, isErr := h.getAuthenticatedUserID(c)
	if isErr {
		return
	}

	req, err := ginUtils.DecodeJSONFromRequest[requestdto.OAuthAuthorizeDecision](c, h.logger)
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

	response, ucErr := h.uc.PostAuthorizeDecision(c.Request.Context(), req, userID, h.config.AuthzCodeExpiry)
	if ucErr != nil {
		//if !req.ReturnResult {
		//	utils.HandleError(c, oauthErr.StatusCode, oauthErr.Description, oauthErr.Err)
		//} else {
		//	c.Redirect(http.StatusSeeOther, "/oauth/error?opt="+string(oauthErr.Code))
		//}
		//ginUtils.HandleError(c, h.logger, oauthErr.StatusCode, oauthErr.Description, oauthErr.Err)
		handlers.HandlerUCError(c, h.logger, "", ucErr)
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
