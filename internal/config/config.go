package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file, error: %v", err)
	}

	return &Config{
		DbURL: os.Getenv("DATABASE_URL"),
	}
}
