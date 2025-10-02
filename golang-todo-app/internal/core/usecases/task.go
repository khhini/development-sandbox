package usecases

import (
	"context"

	"github.com/khhini/golang-todo-app/internal/core/domain"
	"github.com/khhini/golang-todo-app/internal/core/dto"
	"github.com/khhini/golang-todo-app/internal/core/ports"
)

type TaskUsecase struct {
	repo ports.TaskRepository
}

func NewTaskService(repo ports.TaskRepository) ports.TaskUsecase {
	return &TaskUsecase{
		repo: repo,
	}
}

func (uc *TaskUsecase) Create(ctx context.Context, input *dto.CreateTaskRequest) error {
	task := domain.NewTask(input.Title, input.Description)
	return uc.repo.Create(ctx, &task)
}

func (uc *TaskUsecase) GetAll(ctx context.Context) ([]*domain.Task, error) {
	return uc.repo.GetAll(ctx)
}
