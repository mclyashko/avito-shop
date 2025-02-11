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

	user, err := db.GetUserByLoginTx(ctx, tx, username)
	if err != nil {
		return err
	}

	item, err := db.GetItemByNameTx(ctx, tx, itemName)
	if err != nil {
		return err
	}

	if user.Balance < item.Price {
		return ErrInsufficientFunds
	}

	err = db.UpdateUserBalanceTx(ctx, tx, username, -item.Price)
	if err != nil {
		return err
	}

	err = db.InsertUserItemTx(ctx, tx, username, itemName)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
