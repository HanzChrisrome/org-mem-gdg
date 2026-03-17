package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary Health check
// @Description Returns API liveness status.
// @Tags Health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
