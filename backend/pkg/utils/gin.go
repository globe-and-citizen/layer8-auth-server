package utils

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/pkg/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, status int, message string, err error) {
	l := log.Get()
	l.Err(err)

	c.AbortWithStatusJSON(status, responsedto.Response{
		IsSuccess: false,
		Message:   message,
		Error:     strings.Split(err.Error(), "\n"),
	})
}

func ReturnOK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, responsedto.Response{
		IsSuccess: true,
		Message:   message,
		Data:      data,
	})
}

func ReturnCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, responsedto.Response{
		IsSuccess: true,
		Message:   message,
		Data:      data,
	})
}

func DecodeJSONFromRequest[T any](c *gin.Context) (T, error) {
	var request T
	err := c.BindJSON(&request)
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid request payload", err)
		return request, err
	}

	return request, nil
}

func GetBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if !strings.HasPrefix(authHeader, consts.TokenTypeBearer) {
		errorMsg := "invalid authorization header"
		return "", fmt.Errorf(errorMsg)
	}

	return authHeader[len(consts.TokenTypeBearer)+1:], nil
}
