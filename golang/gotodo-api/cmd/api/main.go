package main

import (
	"fmt"
	"log"

	"github.com/khhini/golang-todo-app/internal/config"
	"github.com/khhini/golang-todo-app/internal/server"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	config.SetupLogger(cfg)

	svr := server.NewServer(cfg)

	if err := svr.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		log.Fatalf("http server error: %v", err)
	}
}
