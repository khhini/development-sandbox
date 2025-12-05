package server

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/khhini/golang-todo-app/internal/config"
	"github.com/khhini/golang-todo-app/internal/core/domain"
	"github.com/rs/zerolog"
)

type Server struct {
	app  *fiber.App
	cntr *Container
}

func NewServer(cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{})

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

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
