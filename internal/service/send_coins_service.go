package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/db"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

func SendCoins(ctx context.Context, pool *pgxpool.Pool, sender string, receiver string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	user, err := db.GetUserByLoginTx(ctx, tx, sender)

	if user.Balance < amount {
		return ErrInsufficientFunds
	}

	if err := db.UpdateUserBalance(ctx, tx, sender, -amount); err != nil {
		return err
	}

	if err := db.UpdateUserBalance(ctx, tx, receiver, amount); err != nil {
		return err
	}

	if err := db.InsertCoinTransfer(ctx, tx, sender, receiver, amount); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
