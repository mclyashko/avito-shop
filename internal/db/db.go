package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mclyashko/avito-shop/internal/config"
)

type Db struct {
	pool *pgxpool.Pool
}

func InitDB(cfg *config.Config) *Db {
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

	return &Db {
		pool: pool,
	}
}

func (db *Db) RunInTransaction(ctx context.Context, txFunc func(tx pgx.Tx) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	if err := txFunc(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
