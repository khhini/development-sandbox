package ports

import (
	"context"

	"github.com/khhini/golang-todo-app/internal/core/domain"
	"github.com/khhini/golang-todo-app/internal/core/dto"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetAll(ctx context.Context) ([]*domain.Task, error)
}

type TaskUsecase interface {
	Create(ctx context.Context, input *dto.CreateTaskRequest) error
	GetAll(ctx context.Context) ([]*domain.Task, error)
}
