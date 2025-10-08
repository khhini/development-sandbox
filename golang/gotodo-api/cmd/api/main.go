package main

import (
	"log"
	"net/http"

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

	err = svr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %v", err)
	}
}
