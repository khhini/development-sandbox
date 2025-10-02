package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func NewTask(title, description string) Task {
	return Task{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
}
