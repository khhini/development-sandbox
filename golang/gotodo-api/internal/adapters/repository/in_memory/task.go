package inmemory

import (
	"context"
	"sync"

	"github.com/khhini/golang-todo-app/internal/core/domain"
	"github.com/khhini/golang-todo-app/internal/core/ports"
)

type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewInMemoryTaskRepository(tasks map[string]*domain.Task) ports.TaskRepository {
	return &InMemoryTaskRepository{
		tasks: tasks,
	}
}

func (r *InMemoryTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return ports.ErrDuplicateID
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allTasks := make([]*domain.Task, 0, len(r.tasks))

	for _, task := range r.tasks {
		allTasks = append(allTasks, task)
	}

	return allTasks, nil
}
