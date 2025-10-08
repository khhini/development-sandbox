package server

import (
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.Use(logger.SetLogger(logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
		return l.Output(gin.DefaultWriter).With().
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Str("ip", c.ClientIP()).
			Logger()
	})))

	v1 := r.Group("/api/v1")

	v1.GET("/healthz", s.cntr.healthHandler.Check)

	task := v1.Group("/tasks")
	{
		task.GET("", s.cntr.taskHandler.GetAll)
		task.POST("", s.cntr.taskHandler.Create)
	}

	return r
}
