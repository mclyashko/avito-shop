package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mclyashko/avito-shop/internal/db"
)

// InfoResponse структура ответа с информацией о монетах и инвентаре
type InfoResponse struct {
	Coins       int64       `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

// Item структура для предмета
type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// CoinHistory структура для истории монет
type CoinHistory struct {
	Received []Transaction `json:"received"`
	Sent     []Transaction `json:"sent"`
}

// Transaction структура для транзакции
type Transaction struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int64  `json:"amount"`
}

func GetUserInfo(ctx context.Context, pool *pgxpool.Pool, username string) (*InfoResponse, error) {
	user, err := db.GetUserByLogin(ctx, pool, username)
	if err != nil {
		return nil, err
	}

	user_items, err := db.GetUserItemsByUsername(ctx, pool, username)
	if err != nil {
		return nil, err
	}

	inventory := make([]Item, len(user_items))
	for i, user_item := range user_items {
		inventory[i] = Item{
			Type:     user_item.ItemName,
			Quantity: user_item.Quantity,
		}
	}

	recieved, sent, err := db.GetUserTransactionHistory(ctx, pool, username)
	if err != nil {
		return nil, err
	}

	recievedCoins := make([]Transaction, len(recieved))
	for i, coins := range recieved {
		recievedCoins[i] = Transaction{
			FromUser: coins.SenderID,
			Amount:   coins.Amount,
		}
	}

	sentCoins := make([]Transaction, len(sent))
	for i, coins := range sent {
		sentCoins[i] = Transaction{
			ToUser: coins.ReceiverID,
			Amount: coins.Amount,
		}
	}

	return &InfoResponse{
		Coins:     user.Balance,
		Inventory: inventory,
		CoinHistory: CoinHistory{
			Received: recievedCoins,
			Sent:     sentCoins,
		},
	}, nil
}
