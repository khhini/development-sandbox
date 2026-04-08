package domains

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/khhini/development-sandbox/golang/vector-search-quickstart/internals/config"
	"github.com/ledongthuc/pdf"
	"github.com/pgvector/pgvector-go"
	"google.golang.org/genai"
)

type DocSearchResult struct {
	Name       string
	Content    string
	PageNumber int32
	Similarity float64
}

func SearchDocs(ctx context.Context, conn *pgx.Conn, client *genai.Client, query string) ([]DocSearchResult, error) {
	result, err := client.Models.EmbedContent(ctx, "gemini-embedding-001", []*genai.Content{
		genai.NewContentFromText(query, genai.RoleUser),
	}, &genai.EmbedContentConfig{OutputDimensionality: &config.OutputDimensionality})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %v", err)
	}

	queryVector := pgvector.NewVector(result.Embeddings[0].Values)

	fmt.Printf("searching for: %s \n\n", query)

	rows, _ := conn.Query(ctx, `
		SELECT document_name, content, page_number, 1 - (embedding <=> $1) as similarity
		FROM document_chunks
		WHERE 1 - (embedding <=> $1) > 0.6
		ORDER BY similarity DESC LIMIT 3`, queryVector)
	defer rows.Close()

	results := make([]DocSearchResult, 0, 3)
	for rows.Next() {
		var name, content string
		var page_number int32
		var sim float64
		if err := rows.Scan(&name, &content, &page_number, &sim); err != nil {
			return nil, fmt.Errorf("failed to parsing results: %v", err)
		}
		results = append(results, DocSearchResult{
			name,
			content,
			page_number,
			sim * 100,
		})

	}
	return results, nil
}

func chunkText(text string, chunkSize int, overlap int) []string {
	var chunks []string
	runes := []rune(text)
	for i := 0; i < len(runes); i += chunkSize - overlap {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
		if end == len(runes) {
			break
		}
	}
	return chunks
}

func InjectDoc(ctx context.Context, conn *pgx.Conn, client *genai.Client, filePath string) error {
	_, r, err := pdf.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to load uploaded file: %v", err)
	}

	var fullText string
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		content, _ := p.GetPlainText(nil)
		fullText += content
	}

	fmt.Println(fullText)

	chunks := chunkText(fullText, 500, 100)

	for i, chunk := range chunks {
		result, err := client.Models.EmbedContent(ctx, "gemini-embedding-001", []*genai.Content{
			genai.NewContentFromText(chunk, genai.RoleUser),
		}, &genai.EmbedContentConfig{OutputDimensionality: &config.OutputDimensionality})
		if err != nil {
			return fmt.Errorf("failed to generate embedding: %v", err)
		}

		vector := pgvector.NewVector(result.Embeddings[0].Values)

		if _, err := conn.Exec(ctx, `
			INSERT INTO document_chunks (document_name, content, page_number, embedding)
			VALUES ($1, $2, $3, $4)`,
			filePath, chunk, (i/2)+1, vector); err != nil {
			return fmt.Errorf("faield to add document_chunks to vector db: %v", err)
		}
	}

	return nil
}
