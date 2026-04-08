package domains

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/khhini/development-sandbox/golang/vector-search-quickstart/internals/config"
	"github.com/pgvector/pgvector-go"
	"google.golang.org/genai"
)

type FoodItem struct {
	Name         string
	Description  string
	IslandOrigin string
}

type FoodSearchResult struct {
	Name        string
	Description string
	Similarity  float64
}

func SearchFood(ctx context.Context, conn *pgx.Conn, client *genai.Client, query string) ([]FoodSearchResult, error) {
	result, err := client.Models.EmbedContent(ctx, "gemini-embedding-001", []*genai.Content{
		genai.NewContentFromText(query, genai.RoleUser),
	}, &genai.EmbedContentConfig{OutputDimensionality: &config.OutputDimensionality})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %v", err)
	}

	queryVector := pgvector.NewVector(result.Embeddings[0].Values)

	fmt.Printf("searching for: %s \n\n", query)

	rows, _ := conn.Query(ctx, `
		SELECT name, description, 1 - (embedding <=> $1) as similarity
		FROM indonesian_foods
		ORDER BY similarity DESC LIMIT 3`, queryVector)
	defer rows.Close()

	results := make([]FoodSearchResult, 0, 3)
	for rows.Next() {
		var name, desc string
		var sim float64
		if err := rows.Scan(&name, &desc, &sim); err != nil {
			return nil, fmt.Errorf("failed to parsing results: %v", err)
		}
		results = append(results, FoodSearchResult{
			name,
			desc,
			sim * 100,
		})

	}
	return results, nil
}

func AddFood(ctx context.Context, conn *pgx.Conn, client *genai.Client, food FoodItem) error {
	contentToEmbed := fmt.Sprintf("Nama Makanan: %s, Deskripsi: %s, Asal: %s.",
		food.Name, food.Description, food.IslandOrigin)

	// Generationg vector using gemini-embedding-001
	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		[]*genai.Content{genai.NewContentFromText(contentToEmbed, genai.RoleUser)},
		&genai.EmbedContentConfig{OutputDimensionality: &config.OutputDimensionality},
	)
	if err != nil {
		return fmt.Errorf("failed embedding %s: %v", food.Name, err)
	}

	// Retrive Vector From results
	vector := pgvector.NewVector(result.Embeddings[0].Values)

	_, err = conn.Exec(ctx, `
			INSERT INTO indonesian_foods (name, description, island_origin, embedding)
			VALUES ($1, $2, $3, $4)`,
		food.Name, food.Description, food.IslandOrigin, vector)
	if err != nil {
		return fmt.Errorf("failed to save %s to db: %v", food.Name, err)
	}
	return nil
}
