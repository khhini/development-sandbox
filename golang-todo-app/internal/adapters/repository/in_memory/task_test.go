package inmemory

import (
	"context"
	"testing"

	"github.com/khhini/golang-todo-app/internal/core/domain"
)

func TestInMemoryTaskRepository(t *testing.T) {
	tasks := make(map[string]*domain.Task)
	ctx := context.Background()
	repo := NewInMemoryTaskRepository(tasks)

	t.Run("Create New Task", func(t *testing.T) {
		task := domain.NewTask("New TaskNewTask", "Describe New Task")
		err := repo.Create(ctx, &task)
		if err != nil {
			t.Fatalf("Create() failed: %v", err)
		}
	})

	t.Run("Get All Task", func(t *testing.T) {
		tasks, err := repo.GetAll(ctx)
		if err != nil {
			t.Fatalf("GetAll() failed: %v", err)
		}

		if len(tasks) != 1 {
			t.Fatalf("GetAll() got %d, want 1", len(tasks))
		}
	})
}
