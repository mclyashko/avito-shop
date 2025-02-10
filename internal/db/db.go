package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mclyashko/avito-shop/internal/config"
)

func InitDB(cfg *config.Config) *pgxpool.Pool {
	url := cfg.DbURL

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Failed to create a dbConfig, error: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create DbPool, error: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	return pool
}
