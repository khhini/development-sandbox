package main

import (
	"fmt"

	"github.com/khhini/development-sandbox/golang/waca-go/internal/config"
	"github.com/khhini/development-sandbox/golang/waca-go/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatal().Msgf("failed to initialize config: %v", err)
	}

	svr, err := InitializeServer(log)
	if err != nil {
		log.Fatal().Msgf("failed to initialize server: %v", err)
	}

	if err := svr.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatal().Msgf("failed to run server: %v", err)
	}
}
