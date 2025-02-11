package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/db"
)

type Service interface {
	RunWithTx(ctx context.Context, txFunc func(tx pgx.Tx) error) error
	GetConfig() *config.Config
}

type basicServiceImpl struct {
	db  db.Db
	cfg *config.Config
}

func (s *basicServiceImpl) RunWithTx(ctx context.Context, txFunc func(tx pgx.Tx) error) error {
	return s.db.RunInTransaction(ctx, txFunc)
}

func (s *basicServiceImpl) GetConfig() *config.Config {
	return s.cfg
}

func NewBasicService(db db.Db, cfg *config.Config) Service {
	return &basicServiceImpl{
		db:  db,
		cfg: cfg,
	}
}
