package main

import (
	"fmt"
	"log"

	"github.com/khhini/development-sandbox/golang/waca-go/internal/config"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	svr, err := InitializeServer()
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := svr.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
