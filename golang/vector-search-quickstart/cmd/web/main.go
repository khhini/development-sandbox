package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jackc/pgx/v5"
	"github.com/khhini/development-sandbox/golang/vector-search-quickstart/internals/domains"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
	client, _ := genai.NewClient(ctx, nil)

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Post("/search", func(c *fiber.Ctx) error {
		searchTerm := strings.ToLower(c.FormValue("search"))

		resutls, err := domains.SearchFood(ctx, conn, client, searchTerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong.")
		}

		return c.Render("results", fiber.Map{"Results": resutls})
	})

	app.Get("/add-food", func(c *fiber.Ctx) error {
		return c.Render("add_food_item", nil)
	})

	app.Post("/add-food", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		description := c.FormValue("description")
		islandOrigin := c.FormValue("island-origin")

		newItem := domains.FoodItem{
			Name:         name,
			Description:  description,
			IslandOrigin: islandOrigin,
		}

		if err := domains.AddFood(ctx, conn, client, newItem); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong.")
		}

		return c.Render("partials/success", fiber.Map{
			"Name": name,
		})
	})

	app.Get("/doc-search", func(c *fiber.Ctx) error {
		return c.Render("doc_search", nil)
	})

	app.Post("/doc-search", func(c *fiber.Ctx) error {
		query := strings.ToLower(c.FormValue("q"))

		items, err := domains.SearchDocs(ctx, conn, client, query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong.")
		}

		return c.Render("doc_results", fiber.Map{
			"Docs": items,
		})
	})

	app.Get("/add-doc", func(c *fiber.Ctx) error {
		return c.Render("add_doc_item", nil)
	})

	app.Post("/add-doc", func(c *fiber.Ctx) error {
		file, err := c.FormFile("document")
		if err != nil {
			return c.Status(400).SendString("<span class='error'>No file uploaded</span>")
		}

		savePath := fmt.Sprintf("./uploads/%s", file.Filename)

		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(500).SendString("<span class='error'>Failed to save file</span>")
		}

		if err := domains.InjectDoc(ctx, conn, client, savePath); err != nil {
			return c.Status(500).SendString("<span class='error'>Failed to add document to vector db</span>")
		}

		return c.SendString(fmt.Sprintf("<span class='success'>File '%s' uploaded successfully!</span>", file.Filename))
	})

	app.Listen(":3000")
}
