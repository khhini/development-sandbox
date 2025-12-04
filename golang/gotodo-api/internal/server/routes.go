package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) RegisterRoutes() *fiber.App {
	r := s.app

	// r.Use(logger.SetLogger(logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
	// 	return l.Output(gin.DefaultWriter).With().
	// 		Str("path", c.Request.URL.Path).
	// 		Str("method", c.Request.Method).
	// 		Str("ip", c.ClientIP()).
	// 		Logger()
	// })))
	//
	v1 := r.Group("/api/v1")

	v1.Get("/healthz", s.cntr.healthHandler.Check)

	task := v1.Group("/tasks")
	{
		task.Get("", s.cntr.taskHandler.GetAll)
		task.Post("", s.cntr.taskHandler.Create)
	}

	return r
}
