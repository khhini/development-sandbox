package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/bigquery"
	bigquery_repo "github.com/khhini/golang-sandbox/data-pipeline/internal/adapters/repositories/bigquery"
)

func main() {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "khhini-analytics")
	if err != nil {
		fmt.Printf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	repo := bigquery_repo.NewHailReportBQRepository(client)

	data, err := repo.Get(ctx, 10)
	if err != nil {
		fmt.Printf("failed retrieve data: %v", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("failed marshaling to JSON: %v", err)
	}

	fmt.Print(string(jsonData))
}
