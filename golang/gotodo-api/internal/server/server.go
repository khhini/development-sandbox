package server

import (
	"context"
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/khhini/golang-todo-app/internal/config"
	"github.com/rs/zerolog"
)

type Server struct {
	app  *fiber.App
	cntr *Container
}

func NewServer(cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{})
	ctx := context.Background()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Err(err)
	}

	container := NewContainer(
		WithHealthHandler(),
		WithTaskHandler(conn),
	)

	newServer := Server{
		app:  app,
		cntr: container,
	}

	newServer.RegisterRoutes()

	return app
}
