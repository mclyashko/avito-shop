package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/model"
)

func GetUserTransactionHistory(ctx context.Context, pool *pgxpool.Pool, username string) (recieved []model.CoinTransfer, sent []model.CoinTransfer, err error) {
	query := `
		SELECT id, sender_id, receiver_id, amount
		FROM coin_transfer
		WHERE receiver_id = $1
	`

	rows, err := pool.Query(ctx, query, username)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transfer model.CoinTransfer
		if err := rows.Scan(&transfer.ID, &transfer.SenderID, &transfer.ReceiverID, &transfer.Amount); err != nil {
			return nil, nil, err
		}
		recieved = append(recieved, transfer)
	}

	query = `
		SELECT id, sender_id, receiver_id, amount
		FROM coin_transfer
		WHERE sender_id = $1
	`

	rows, err = pool.Query(ctx, query, username)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transfer model.CoinTransfer
		if err := rows.Scan(&transfer.ID, &transfer.SenderID, &transfer.ReceiverID, &transfer.Amount); err != nil {
			return nil, nil, err
		}
		sent = append(sent, transfer)
	}

	return recieved, sent, nil
}