package bigquery_repo

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/khhini/golang-sandbox/data-pipeline/internal/core/domain"
	"github.com/khhini/golang-sandbox/data-pipeline/internal/core/ports"
	"google.golang.org/api/iterator"
)

type HailReportBQRepository struct {
	client *bigquery.Client
}

func NewHailReportBQRepository(client *bigquery.Client) ports.HailReportRepository {
	return &HailReportBQRepository{
		client: client,
	}
}

func (r *HailReportBQRepository) Get(ctx context.Context, limit int64) ([]domain.HailReport, error) {
	q := r.client.Query(`
		SELECT timestamp, time, size, location, county, state, latitude, longitude, comments, report_point
		FROM bigquery-public-data.noaa_preliminary_severe_storms.hail_reports
		ORDER BY timestamp ASC
		LIMIT @limit
		`)

	q.Parameters = []bigquery.QueryParameter{
		{Name: "limit", Value: limit},
	}

	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed read bigquery: %w", err)
	}

	results := make([]domain.HailReport, 0, it.TotalRows)

	var row domain.HailReport

	var problems []error
	for {
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			problems = append(problems, err)
		}
		results = append(results, row)
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}

	return results, nil
}
