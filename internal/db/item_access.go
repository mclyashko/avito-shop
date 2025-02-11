package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/model"
)

var (
	ErrItemNotFound = fmt.Errorf("item not found")
)

type ItemAccessor interface {
	GetItemByNameTx(ctx context.Context, tx pgx.Tx, itemName string) (*model.Item, error)
}

type ItemAccessorImpl struct {
	*Db
}

func (db *ItemAccessorImpl) GetItemByNameTx(ctx context.Context, tx pgx.Tx, itemName string) (*model.Item, error) {
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
		return nil, err
	}

	return &item, nil
}
