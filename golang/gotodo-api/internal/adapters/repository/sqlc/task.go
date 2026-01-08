package sqlc

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khhini/golang-todo-app/internal/core/domain"
	"github.com/khhini/golang-todo-app/internal/core/ports"
	"github.com/khhini/golang-todo-app/internal/infra/sqlc/tasks"
)

type SqlcTaskRepository struct {
	queries *tasks.Queries
}

func NewSqlcTaskRepository(q *tasks.Queries) ports.TaskRepository {
	return &SqlcTaskRepository{
		queries: q,
	}
}

func (r *SqlcTaskRepository) Create(ctx context.Context, input *domain.Task) error {
	r.queries.Create(ctx, tasks.CreateParams{
		Title:       input.Title,
		Description: pgtype.Text{String: input.Description, Valid: true},
	})
	return nil
}

func (r *SqlcTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.queries.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Task, len(rows))

	for i, row := range rows {
		result[i] = &domain.Task{
			ID:          strconv.Itoa(int(row.ID)),
			Title:       row.Title,
			Description: row.Description.String,
			Completed:   row.Completed.Bool,
			CreatedAt:   row.CreatedAt.Time,
			CompletedAt: &row.CompletedAt.Time,
		}
	}
	return result, nil
}
