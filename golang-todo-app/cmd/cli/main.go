package main

import (
	"github.com/khhini/golang-todo-app/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, _ := config.LoadFromEnv()
	config.SetupLogger(cfg)

	if e := log.Debug(); e.Enabled() {
		e.Str("foo", "bar").Msg("some debug message")
	}
}
