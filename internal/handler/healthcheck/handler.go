package healthcheck

import (
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -destination=handler_mock_test.go -package=healthcheck
type Service interface {
	HealthcheckServices() (string, error)
}

type Handler struct {
	*gin.Engine
	Service Service
}

func NewHandler(api *gin.Engine, hca Service) *Handler {
	return &Handler{
		api,
		hca,
	}
}

func (h *Handler) RegisterRoute() {
	h.GET("/health", h.HealthcheckHandlerHTTP)
}
