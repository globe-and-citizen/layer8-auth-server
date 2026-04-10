package clientH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/handlers"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) GetNTorCertificate(c *gin.Context) {
	backendDomain, err := utils.GetURLHostPort(c.Query("backend_url"))
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "backend_url is missing or invalid", err)
		return
	}

	response, ucerr := h.uc.GetNTorCertificate(c.Request.Context(), requestdto.ClientGetNTorCertificate{
		BackendURI: backendDomain,
	})
	if ucerr != nil {
		handlers.HandlerUCError(c, h.logger, "Failed to get client ntor certificate", ucerr)
		return
	}

	//utils.ReturnOK(c, "Get client ntor certificate successful", response)
	c.JSON(http.StatusOK, response)
}
