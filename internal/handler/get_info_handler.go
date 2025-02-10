package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/service"
)

func GetInfo(c *fiber.Ctx, ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) error {
	claims, ok := c.Locals("claims").(*service.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Unauthorized"})
	}

	username := claims.Username

	userInfo, err := service.GetUserInfo(ctx, pool, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Failed to retrieve user info"})
	}

	return c.Status(fiber.StatusOK).JSON(userInfo)
}
