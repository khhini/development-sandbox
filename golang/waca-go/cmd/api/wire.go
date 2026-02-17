//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/khhini/development-sandbox/golang/waca-go/internal/adapter/handlers"
	"github.com/khhini/development-sandbox/golang/waca-go/internal/server"
	"github.com/rs/zerolog"
)

var providerSet = wire.NewSet(
	handlers.NewHealthHandler,
	wire.Struct(new(server.HandlerRegistry), "*"),
	server.NewServer,
)

func InitializeServer(logger *zerolog.Logger) (*fiber.App, error) {
	wire.Build(providerSet)
	return &fiber.App{}, nil
}
