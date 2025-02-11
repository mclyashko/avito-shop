package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/middleware"
	"github.com/mclyashko/avito-shop/internal/service"
)

type infoResponse struct {
	Coins       int64       `json:"coins"`
	Inventory   []item      `json:"inventory"`
	CoinHistory coinHistory `json:"coinHistory"`
}

type item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type coinHistory struct {
	Received []transaction `json:"received"`
	Sent     []transaction `json:"sent"`
}

type transaction struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int64  `json:"amount"`
}

func GetInfo(c *fiber.Ctx, pool *pgxpool.Pool) error {
	ctx := c.Context()

	claims, ok := c.Locals(middleware.LocalsClaimsKey).(*service.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "No token"})
	}

	username := claims.Username

	balance, userItems, recievedTransfers, sentTransfers, err := service.GetUserInfo(ctx, pool, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Failed to retrieve user info"})
	}

	inventory := make([]item, len(userItems))
	for i, userItem := range userItems {
		inventory[i] = item{
			Type:     userItem.ItemName,
			Quantity: userItem.Quantity,
		}
	}

	recievedCoins := make([]transaction, len(recievedTransfers))
	for i, recieved := range recievedTransfers {
		recievedCoins[i] = transaction{
			FromUser: recieved.SenderLogin,
			Amount:   recieved.Amount,
		}
	}

	sentCoins := make([]transaction, len(sentTransfers))
	for i, sent := range sentTransfers {
		sentCoins[i] = transaction{
			ToUser: sent.ReceiverLogin,
			Amount: sent.Amount,
		}
	}

	userInfo := infoResponse{
		Coins:     *balance,
		Inventory: inventory,
		CoinHistory: coinHistory{
			Received: recievedCoins,
			Sent:     sentCoins,
		},
	}

	return c.Status(fiber.StatusOK).JSON(userInfo)
}
