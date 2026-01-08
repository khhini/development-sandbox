package main

import (
	"context"
	"log"
	"os"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khhini/development-sandbox/golang/sqlc-quickstart/tutorial"
)

func run() error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Println(err)
		}
	}()

	queries := tutorial.New(conn)

	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}

	log.Println(authors)

	insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Kiki",
		Bio:  pgtype.Text{String: "Owner of khhini.dev", Valid: true},
	})
	if err != nil {
		return err
	}

	log.Println(insertedAuthor)

	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return err
	}

	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
