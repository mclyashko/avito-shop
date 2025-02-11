package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/model"
)

func GetUserInfo(ctx context.Context, pool *pgxpool.Pool, username string) (balance *int64, userItems []model.UserItem, recievedTransfers []model.CoinTransfer, sentTransfers []model.CoinTransfer, err error) {
	user, err := db.GetUserByLogin(ctx, pool, username)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	balance = &user.Balance

	userItems, err = db.GetUserItemsByUsername(ctx, pool, username)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	recievedTransfers, sentTransfers, err = db.GetUserTransactionHistory(ctx, pool, username)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return balance, userItems, recievedTransfers, sentTransfers, nil
}
