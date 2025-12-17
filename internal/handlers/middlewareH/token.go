package middlewareH

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/usecases/middlewareUC"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
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
	tokenString := c.GetHeader("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix

	userID, err := h.uc.VerifyUserJWTToken(tokenString)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	c.Set(consts.MiddlewareKeyUserUserID, userID)
	c.Next()
}

func (h *MiddlewareHandler) AuthenticateClient(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix

	clientID, username, err := h.uc.VerifyClientJWTToken(tokenString)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	c.Set(consts.MiddlewareKeyClientUsername, username)
	c.Set(consts.MiddlewareKeyClientClientID, clientID)
	c.Next()
}
