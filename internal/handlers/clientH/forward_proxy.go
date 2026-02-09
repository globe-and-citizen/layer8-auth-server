package clientH

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h ClientHandler) GetNTorCertificate(c *gin.Context) {
	backendDomain, err := utils.GetURLHostPort(c.Query("backend_url"))
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "invalid backend url", err)
		return
	}

	response, err := h.uc.GetNTorCertificate(c.Request.Context(), requestdto.ClientGetNTorCertificate{
		BackendURI: backendDomain,
	})
	if err != nil {
		ginUtils.HandleError(c, h.logger, http.StatusBadRequest, "Failed to get client ntor certificate", err)
		return
	}

	//utils.ReturnOK(c, "Get client ntor certificate successful", response)
	c.JSON(http.StatusOK, response)
}
