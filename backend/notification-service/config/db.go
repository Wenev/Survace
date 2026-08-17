package config

import (
	"os"
)

func DbURLFromEnv() string {
	dbURL := os.Getenv("DATABASE_URL")
	return dbURL
}
