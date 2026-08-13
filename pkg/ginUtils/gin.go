package ginUtils

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	IsSuccess bool        `json:"is_success"`
	Message   string      `json:"message"`
	Error     interface{} `json:"errors"`
	Data      interface{} `json:"data"`
}

func HandleError(c *gin.Context, logger log.ILogger, status int, message string, err error) {
	logger.Error(message, err)
	c.AbortWithStatusJSON(status, Response{
		IsSuccess: false,
		Message:   message,
		Error:     strings.Split(err.Error(), "\n"),
	})
}

func ReturnOK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   message,
		Data:      data,
	})
}

func ReturnCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		IsSuccess: true,
		Message:   message,
		Data:      data,
	})
}

func DecodeJSONFromRequest[T any](c *gin.Context, logger log.ILogger) (T, error) {
	var request T
	err := c.BindJSON(&request)
	if err != nil {
		HandleError(c, logger, http.StatusBadRequest, "Invalid request payload", err)
		return request, err
	}

	return request, nil
}

func GetBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	token := strings.Split(authHeader, " ")

	if len(token) != 2 || token[0] != "Bearer" || token[1] == "" {
		return "", fmt.Errorf("invalid authorization header")
	}

	return token[1], nil
}
