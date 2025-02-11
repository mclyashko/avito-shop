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

const (
	itemParamKey = "item"
)

func BuyItemHandler(c *fiber.Ctx, pool *pgxpool.Pool) error {
	ctx := c.Context()

	claims, ok := c.Locals(middleware.LocalsClaimsKey).(*service.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "No token"})
	}

	username := claims.Username
	itemName := c.Params(itemParamKey)

	err := service.BuyItem(ctx, pool, username, itemName)
	if err != nil {
		if errors.Is(err, db.ErrItemNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Item not found"})
		}
		if errors.Is(err, service.ErrInsufficientFunds) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Insufficient funds"})
		}
		log.Printf("Error buying item %v, error: %v", itemName, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process purchase"})
	}

	c.Status(fiber.StatusOK)

	return nil
}
