package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/db"
)

type Service struct {
	cfg *config.Config
	db  *db.Db
}

func NewService(cfg *config.Config, db *db.Db) *Service{
	return &Service{
		cfg: cfg,
		db: db,
	}
}

func (s *Service) RunWithTx(ctx context.Context, txFunc func(tx pgx.Tx) error) error {
	return s.db.RunInTransaction(ctx, txFunc)
}
