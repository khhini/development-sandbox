package server

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) RegisterRoutes() *fiber.App {
	r := s.app

	r.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
	}))

	v1 := r.Group("/api/v1")

	v1.Get("/healthz", s.cntr.healthHandler.Check)

	task := v1.Group("/tasks")
	{
		task.Get("", s.cntr.taskHandler.GetAll)
		task.Post("", s.cntr.taskHandler.Create)
	}

	return r
}
