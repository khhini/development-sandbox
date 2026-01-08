package server

import (
	"github.com/jackc/pgx/v5"
	httphandler "github.com/khhini/golang-todo-app/internal/adapters/handler/http"
	"github.com/khhini/golang-todo-app/internal/adapters/repository/sqlc"
	"github.com/khhini/golang-todo-app/internal/core/usecases"
	"github.com/khhini/golang-todo-app/internal/infra/sqlc/tasks"
)

type Container struct {
	healthHandler *httphandler.HealthHandler
	taskHandler   *httphandler.TaskHandler
}

type ContainerOption func(c *Container)

func NewContainer(opts ...ContainerOption) *Container {
	ctr := &Container{}

	for _, opt := range opts {
		opt(ctr)
	}

	return ctr
}

func WithHealthHandler() ContainerOption {
	handler := httphandler.NewHealthHandler()
	return func(ctr *Container) {
		ctr.healthHandler = &handler
	}
}

func WithTaskHandler(conn *pgx.Conn) ContainerOption {
	queries := tasks.New(conn)
	repo := sqlc.NewSqlcTaskRepository(queries)
	uc := usecases.NewTaskService(repo)
	handler := httphandler.NewTaskHandler(uc)
	return func(ctr *Container) {
		ctr.taskHandler = &handler
	}
}
