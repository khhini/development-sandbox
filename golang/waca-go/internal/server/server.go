// Package server
package server

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/khhini/development-sandbox/golang/waca-go/internal/adapter/handlers"
	"github.com/rs/zerolog"
)

type HandlerRegistry struct {
	Health *handlers.HealthHandler
}

func NewServer(
	log *zerolog.Logger,
	h *HandlerRegistry,
) *fiber.App {
	app := fiber.New(fiber.Config{})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: log,
	}))

	v1 := app.Group("/api/v1")

	h.Health.Register(v1.Group("/healthz"))

	return app
}
