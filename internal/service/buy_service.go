package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/db"
)

func BuyItem(ctx context.Context, pool *pgxpool.Pool, username string, itemName string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	item, err := db.GetItemByName(ctx, tx, itemName)
	if err != nil {
		return err
	}

	user, err := db.GetUserByLoginTx(ctx, tx, username)
	if err != nil {
		return err
	}

	if user.Balance < item.Price {
		return ErrInsufficientFunds
	}

	err = db.UpdateUserBalance(ctx, tx, username, -item.Price)
	if err != nil {
		return err
	}

	err = db.AddUserItem(ctx, tx, username, itemName)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
