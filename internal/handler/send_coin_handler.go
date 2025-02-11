package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/middleware"
	"github.com/mclyashko/avito-shop/internal/service"
)

type sendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}

func SendCoinHandler(c *fiber.Ctx, pool *pgxpool.Pool) error {
	ctx := c.Context()

	claims, ok := c.Locals(middleware.LocalsClaimsKey).(*service.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "No token"})
	}

	username := claims.Username

	var req sendCoinRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body, error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request"})
	}

	err := service.SendCoins(ctx, pool, username, req.ToUser, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrNegativeSignTransaction) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Negative sign transaction"})
		}
		if errors.Is(err, db.ErrUserNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Reciever not found"})
		}
		if errors.Is(err, service.ErrInsufficientFunds) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient funds"})
		}
		log.Printf("Error sending coins from %v to %v, error: %v", username, req.ToUser, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	c.Status(fiber.StatusOK)

	return nil
}
