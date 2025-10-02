package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port int
	Host string

	LogLevel string
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func LoadFromEnv() (Config, error) {
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return Config{}, fmt.Errorf("config: invalid port value '%s': %w", portStr, err)
	}
	host := getEnv("HOST", "0.0.0.0")

	return Config{
		Port: port,
		Host: host,
	}, nil
}
