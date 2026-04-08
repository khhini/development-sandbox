package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/khhini/development-sandbox/golang/vector-search-quickstart/internals/domains"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	connStr := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
	pgxvector.RegisterTypes(ctx, conn)

	client, _ := genai.NewClient(ctx, nil)

	userQuery := "Mongodb Scalability & reliability"

	results, _ := domains.SearchDocs(ctx, conn, client, userQuery)

	for _, r := range results {
		fmt.Printf("Name: %s\nDescription: %s\nSimilarity: %.2f\nPageNumber: %d\n\n", r.Name, r.Content, r.Similarity, r.PageNumber)
	}
}
