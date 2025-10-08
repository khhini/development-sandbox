package server

import (
	httphandler "github.com/khhini/golang-todo-app/internal/adapters/handler/http"
	inmemory "github.com/khhini/golang-todo-app/internal/adapters/repository/in_memory"
	"github.com/khhini/golang-todo-app/internal/core/domain"
	"github.com/khhini/golang-todo-app/internal/core/usecases"
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

func WithTaskHandler(tasksMemory map[string]*domain.Task) ContainerOption {
	repo := inmemory.NewInMemoryTaskRepository(tasksMemory)
	uc := usecases.NewTaskService(repo)
	handler := httphandler.NewTaskHandler(uc)

	return func(ctr *Container) {
		ctr.taskHandler = &handler
	}
}
