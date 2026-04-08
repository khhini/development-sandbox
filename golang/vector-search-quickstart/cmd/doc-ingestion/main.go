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

var embeddingOutputDimention = int32(768)

var filePaths = []string{
	"./uploads/mongodb_7_competitive_advantages.pdf",
}

func seedDatabase(ctx context.Context, client *genai.Client, conn *pgx.Conn) {
	for _, file := range filePaths {
		if err := domains.InjectDoc(ctx, conn, client, file); err != nil {
			log.Printf("failed to save %s to db: %v", file, err)
		} else {
			fmt.Println("indexing success: %s \n", file)
		}
	}
}

func main() {
	ctx := context.Background()

	client, _ := genai.NewClient(ctx, nil)

	conn, _ := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	defer conn.Close(ctx)
	_ = pgxvector.RegisterTypes(ctx, conn)

	seedDatabase(ctx, client, conn)
}
