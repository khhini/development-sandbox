package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Article struct {
	ID      int
	Title   string
	Content string
	Rank    float64
}

func searchArticles(db *sqlx.DB, query string) ([]Article, error) {
	var articles []Article

	sqlQuery := `
		SELECT
			id,
			title,
			content,
			ts_rank(search_vector, websearch_to_tsquery('english', $1)) as rank
		FROM articles
		WHERE search_vector @@ websearch_to_tsquery('english', $1)
		ORDER BY rank DESC
		LIMIT 20
	`

	rows, err := db.Query(sqlQuery, query)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	for rows.Next() {
		var article Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Rank); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func main() {
	dbConn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", dbConn)
	if err != nil {
		log.Fatalf("failed connecting to database: %v", err)
	}

	articles, err := searchArticles(db, "Tolkien")
	if err != nil {
		log.Fatalf("failed searching articles: %v", err)
	}

	fmt.Println(articles)
}
