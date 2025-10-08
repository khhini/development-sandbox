package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (h *HealthHandler) Check(ctx *gin.Context) {
	if e := log.Debug(); e.Enabled() {
		e.Str("foo", "bar").Msg("Test Debug Logging")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
