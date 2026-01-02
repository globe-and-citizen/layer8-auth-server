package middlewareH

import (
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/middlewareUC"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MiddlewareHandler struct {
	uc middlewareUC.IMiddlewareUsecase
}

func NewMiddlewareHandler(uc middlewareUC.IMiddlewareUsecase) *MiddlewareHandler {
	return &MiddlewareHandler{uc: uc}
}

func (h *MiddlewareHandler) AuthenticateUser(c *gin.Context) {
	token, err := utils.GetBearerToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, username, err := h.uc.VerifyUserJWTToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyUserUsername, username)
	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h *MiddlewareHandler) AuthenticateClient(c *gin.Context) {
	token, err := utils.GetBearerToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	clientID, username, err := h.uc.VerifyClientJWTToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyClientUsername, username)
	c.Set(consts.MiddlewareKeyClientClientID, clientID)
	c.Next()
}

func (h *MiddlewareHandler) AuthenticateOAuth(c *gin.Context) {
	token, err := c.Cookie(consts.OAuthCookieName)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, username, err := h.uc.VerifyOAuthJWTToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyUserUsername, username)
	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h *MiddlewareHandler) ValidateAccessToken(c *gin.Context) {
	token, err := utils.GetBearerToken(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: missing token", err)
		return
	}

	userID, scopes, err := h.uc.VerifyAccessToken(token)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	// save claims in context for further handlers
	c.Set(consts.MiddlewareKeyOAuthScopes, scopes)
	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}
