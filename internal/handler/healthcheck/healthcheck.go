package healthcheck

import (
	"ewallet-wallet/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HealthcheckHandlerHTTP(c *gin.Context) {
	msg, err := h.Service.HealthcheckServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, msg, nil)
}
