package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port   string
	DBPath string
}

func Load() *Config {
	return &Config{
		Port:   getEnv("PORT", "8080"),
		DBPath: getEnv("DB_PATH", "./tms.db"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
