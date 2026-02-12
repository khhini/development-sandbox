// Package server
package server

import "github.com/gofiber/fiber/v2"

type Server struct{}

func NewServer() *fiber.App {
	app := fiber.New(fiber.Config{})

	return app
}
