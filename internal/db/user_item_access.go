package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/model"
)

func GetUserItemsByUsername(ctx context.Context, pool *pgxpool.Pool, username string) ([]model.UserItem, error) {
	query := `
		SELECT id, user_id, item_name, quantity
		FROM user_item
		WHERE user_id = $1
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