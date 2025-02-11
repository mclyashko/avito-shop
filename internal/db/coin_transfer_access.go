package db

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		if err := rows.Scan(&transfer.ID, &transfer.SenderLogin, &transfer.ReceiverLogin, &transfer.Amount); err != nil {
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
		if err := rows.Scan(&transfer.ID, &transfer.SenderLogin, &transfer.ReceiverLogin, &transfer.Amount); err != nil {
			return nil, nil, err
		}
		sent = append(sent, transfer)
	}

	return recieved, sent, nil
}

func InsertCoinTransferTx(ctx context.Context, tx pgx.Tx, sender string, reciever string, amount int64) error {
	query := `
		INSERT INTO coin_transfer (id, sender_id, receiver_id, amount) 
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, query, uuid.New(), sender, reciever, amount)
	if err != nil {
		log.Printf("Failed to insert coin transfer, sender: %v, reciever: %v error: %v", sender, reciever, err)
		return err
	}

	return nil
}
