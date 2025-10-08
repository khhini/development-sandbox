package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/khhini/golang-todo-app/internal/config"
	"github.com/khhini/golang-todo-app/internal/core/domain"
)

type Server struct {
	host string
	port int
	cntr *Container
}

func NewServer(cfg config.Config) *http.Server {
	taskMemory := make(map[string]*domain.Task)

	container := NewContainer(
		WithHealthHandler(),
		WithTaskHandler(taskMemory),
	)
	newServer := Server{
		host: cfg.Host,
		port: cfg.Port,
		cntr: container,
	}

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", newServer.host, newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
