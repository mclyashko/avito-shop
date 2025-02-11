package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/db"
)

var ErrNegativeSignTransaction = errors.New("negative sign")
var ErrInsufficientFunds = errors.New("insufficient funds")

func SendCoins(ctx context.Context, pool *pgxpool.Pool, sender string, receiver string, amount int64) error {
	if amount <= 0 {
		return ErrNegativeSignTransaction
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	senderUser, err := db.GetUserByLoginTx(ctx, tx, sender)
	if err != nil {
		return err
	}

	_, err = db.GetUserByLoginTx(ctx, tx, receiver)
	if err != nil {
		return err
	}

	if senderUser.Balance < amount {
		return ErrInsufficientFunds
	}

	if err := db.UpdateUserBalanceTx(ctx, tx, sender, -amount); err != nil {
		return err
	}

	if err := db.UpdateUserBalanceTx(ctx, tx, receiver, amount); err != nil {
		return err
	}

	if err := db.InsertCoinTransferTx(ctx, tx, sender, receiver, amount); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
