package handler

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/service"
)

func BuyItemHandler(c *fiber.Ctx, ctx context.Context, pool *pgxpool.Pool) error {
	claims, ok := c.Locals("claims").(*service.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Unauthorized"})
	}

	username := claims.Username

	itemName := c.Params("item")
	if itemName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Item name is required"})
	}

	err := service.BuyItem(ctx, pool, username, itemName)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientFunds) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Insufficient funds"})
		}
		if errors.Is(err, db.ErrItemNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Item not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process purchase"})
	}

	c.Status(fiber.StatusOK)

	return nil
}
