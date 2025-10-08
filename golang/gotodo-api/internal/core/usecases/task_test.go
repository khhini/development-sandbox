package usecases

import (
	"context"
	"testing"

	inmemory "github.com/khhini/golang-todo-app/internal/adapters/repository/in_memory"
	"github.com/khhini/golang-todo-app/internal/core/domain"
	"github.com/khhini/golang-todo-app/internal/core/dto"
)

func TestTaskService(t *testing.T) {
	tasks := make(map[string]*domain.Task)
	ctx := context.Background()
	repo := inmemory.NewInMemoryTaskRepository(tasks)
	svc := NewTaskService(repo)

	t.Run("Create Task", func(t *testing.T) {
		newTaskRequest := dto.CreateTaskRequest{
			Title:       "New Created Task",
			Description: "Describe New Created Task",
		}
		if err := svc.Create(ctx, &newTaskRequest); err != nil {
			t.Fatalf("Create() failed: %v", err)
		}
	})

	t.Run("Get All Tasks", func(t *testing.T) {
		tasks, err := svc.GetAll(ctx)
		if err != nil {
			t.Fatalf("GetAll() failed: %v", err)
		}

		if len(tasks) != 1 {
			t.Fatalf("GetAll() got %d, want 1", len(tasks))
		}
	})
}
