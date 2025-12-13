package tokenHandler

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/usecases/tokenUsecase"
	"globe-and-citizen/layer8/auth-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	uc tokenUsecase.ITokenUseCase
}

func NewTokenHandler(uc tokenUsecase.ITokenUseCase) *TokenHandler {
	return &TokenHandler{uc: uc}
}

func (h *TokenHandler) UserAuthentication(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = tokenString[7:] // Remove the "Bearer " prefix

	userID, err := h.uc.VerifyUserJWTToken(tokenString)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Authentication error: invalid token", err)
		return
	}

	c.Set(consts.MiddlewareKeyUserID, userID)
}
