// Package server
package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khhini/development-sandbox/golang/waca-go/internal/adapter/handlers"
)

type HandlerRegistry struct {
	Health *handlers.HealthHandler
}

type ServerOptions struct {
	// log *log.Logger
	H *HandlerRegistry
}

func NewServer(
	opts ServerOptions,
) *fiber.App {
	app := fiber.New(fiber.Config{})

	v1 := app.Group("/api/v1")

	opts.H.Health.Register(v1.Group("/healthz"))

	return app
}
