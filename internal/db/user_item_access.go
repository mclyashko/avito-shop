package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/model"
)

func GetUserItemsByUsername(ctx context.Context, pool *pgxpool.Pool, username string) ([]model.UserItem, error) {
	query := `
		SELECT id, user_login, item_name, quantity
		FROM user_item
		WHERE user_login = $1
	`

	rows, err := pool.Query(ctx, query, username)
	if (err != nil) {
		return nil, err
	}
	defer rows.Close()

	var items []model.UserItem
	for rows.Next() {
		var item model.UserItem
		if err := rows.Scan(&item.ID, &item.UserLogin, &item.ItemName, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	return items, nil
}

func InsertUserItemTx(ctx context.Context, tx pgx.Tx, username string, itemName string) error {
	query := `
		INSERT INTO user_item (id, user_login, item_name, quantity)
		VALUES ($1, $2, $3, 1)
		ON CONFLICT (user_login, item_name)
		DO UPDATE SET quantity = user_item.quantity + 1
	`

	_, err := tx.Exec(ctx, query, uuid.New(), username, itemName)
	if err != nil {
		return err
	}

	return nil
}