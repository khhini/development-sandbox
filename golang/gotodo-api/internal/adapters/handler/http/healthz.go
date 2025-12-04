package httphandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (h *HealthHandler) Check(ctx *fiber.Ctx) error {
	if e := log.Debug(); e.Enabled() {
		e.Str("foo", "bar").Msg("Test Debug Logging")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}
