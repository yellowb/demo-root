package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Addr   string
	DBPath string
}

func Load() Config {
	return Config{
		Addr:   normalizeAddr(firstNonEmpty(os.Getenv("TODO_APP_ADDR"), os.Getenv("PORT"), "8080")),
		DBPath: firstNonEmpty(os.Getenv("TODO_DB_PATH"), filepath.Join("data", "todos.db")),
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

func normalizeAddr(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ":8080"
	}
	if strings.HasPrefix(trimmed, ":") {
		return trimmed
	}
	if strings.Contains(trimmed, ":") {
		return trimmed
	}

	return ":" + trimmed
}
