package ports

import (
	"context"

	"github.com/khhini/golang-sandbox/data-pipeline/internal/core/domain"
)

type HailReportRepository interface {
	Get(ctx context.Context, limit int64) ([]domain.HailReport, error)
}
