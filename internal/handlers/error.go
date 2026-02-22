package handlers

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlerUCError(c *gin.Context, logger log.ILogger, message string, err *ucerror.UCError) {
	if err != nil {
		status := consts.MapErrorToStatusCode(err.Public())

		logger.Error(message, err.Details())
		c.AbortWithStatusJSON(status, ginUtils.Response{
			IsSuccess: false,
			Message:   message,
			Error:     err.Error(),
		})
		return
	}

	logger.Error(message, nil)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ginUtils.Response{
		IsSuccess: false,
		Message:   "",
		Error:     fmt.Errorf("internal server error"),
	})
}
