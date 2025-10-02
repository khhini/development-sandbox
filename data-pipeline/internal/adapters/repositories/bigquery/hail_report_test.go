package bigquery_repo

import (
	"context"
	"encoding/json"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/khhini/golang-sandbox/data-pipeline/internal/core/domain"
)

func TestHailReportBQRepository(t *testing.T) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "khhini-analytics")
	if err != nil {
		t.Errorf("create client err = %v, want nil", err)
	}

	wantByte := []byte(`{"timestamp":"2021-01-06T17:53:00Z","time":"1753","size":100,"location":"7 NE MILANO","county":"MILAM","state":"TX","latitude":30.76,"longitude":-96.75,"comments":"PUBLIC REPORT OF QUARTER SIZED HAIL NE OF MILANO. (FWD)","report_point":"POINT(-96.75 30.76)"}`)
	var want domain.HailReport
	if err := json.Unmarshal(wantByte, &want); err != nil {
		t.Errorf("unmarshal err = %v, want nil", err)
	}

	repo := NewHailReportBQRepository(client)

	data, err := repo.Get(ctx, 10)
	if err != nil {
		t.Errorf("repo.Get() err = %v, want nil", err)
	}

	if len(data) <= 0 {
		t.Errorf("len(data) = %v, want > 0", len(data))
	}

	if data[0].Time != want.Time {
		t.Errorf("data[0].Time = %v, want %v", data[0].Time, want.Time)
	}
}
