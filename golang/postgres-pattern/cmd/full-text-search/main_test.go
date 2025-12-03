package main

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestSearchArticles(t *testing.T) {
	t.Parallel()

	dbConn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", dbConn)
	if err != nil {
		t.Errorf("db connection return error: %q", err)
	}

	tests := map[string]struct {
		input string
		count int
	}{
		"search for hobbit":  {"Hobbit", 1},
		"search for tolkien": {"Tolkien", 2},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Log(name)
			got, err := searchArticles(db, test.input)
			if err != nil {
				t.Errorf("searchArticles(%q) return error: %v", test.input, err)
			}

			if len(got) != test.count {
				t.Errorf("searchArticles(%q) got %d entry, expected %d entry", test.input, len(got), test.count)
			}
		})
	}
}
