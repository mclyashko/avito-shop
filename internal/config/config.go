package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DbURL                 string
	JwtSecretKey          []byte
	JwtExpirationDuration time.Duration
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file, error: %v", err)
	}

	jwtExpirationString := os.Getenv("JWT_EXPIRATION_DURATION")
	jwtExpirationDuration, err := time.ParseDuration(jwtExpirationString)
	if err != nil {
		log.Fatalf("Cant parse jwtExpirationString %v, error : %v", jwtExpirationString, err)
	}

	return &Config{
		DbURL:                 os.Getenv("DATABASE_URL"),
		JwtSecretKey:          []byte(os.Getenv("JWT_SECRET_KEY")),
		JwtExpirationDuration: jwtExpirationDuration,
	}
}
