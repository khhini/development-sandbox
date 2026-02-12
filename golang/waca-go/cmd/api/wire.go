//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/khhini/development-sandbox/golang/waca-go/internal/server"
)

var serverProviderSet = wire.NewSet(
	server.NewServer,
)

func InitializeServer() (*fiber.App, error) {
	wire.Build(serverProviderSet)
	return &fiber.App{}, nil
}
