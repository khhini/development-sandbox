package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/khhini/golang-todo-app/internal/config"
	"github.com/khhini/golang-todo-app/internal/core/domain"
)

type Server struct {
	app  *fiber.App
	cntr *Container
}

func NewServer(cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(logger.New())

	taskMemory := make(map[string]*domain.Task)

	container := NewContainer(
		WithHealthHandler(),
		WithTaskHandler(taskMemory),
	)

	newServer := Server{
		app:  app,
		cntr: container,
	}

	newServer.RegisterRoutes()

	return app
}
