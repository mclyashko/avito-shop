package db

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/model"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

func GetItemByNameTx(ctx context.Context, tx pgx.Tx, itemName string) (*model.Item, error) {
	query := `
		SELECT name, price 
		FROM item 
		WHERE name = $1
	`

	row := tx.QueryRow(ctx, query, itemName)

	var item model.Item

	err := row.Scan(&item.Name, &item.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrItemNotFound
		}
		log.Printf("Failed to get user by login: %v", err)
		return nil, err
	}

	log.Printf("Got item: %v", item)
	return &item, nil
}

func InsertUserItemTx(ctx context.Context, tx pgx.Tx, username string, itemName string) error {
	query := `
		INSERT INTO user_item (id, user_id, item_name, quantity)
		VALUES ($1, $2, $3, 1)
		ON CONFLICT (user_id, item_name)
		DO UPDATE SET quantity = user_item.quantity + 1
	`

	_, err := tx.Exec(ctx, query, uuid.New(), username, itemName)
	if err != nil {
		return err
	}

	return nil
}
