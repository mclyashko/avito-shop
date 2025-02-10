package handler

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/service"
)

// AuthRequest структура запроса на аутентификацию
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse структура ответа с токеном
type AuthResponse struct {
	Token string `json:"token"`
}

func Authenticate(c *fiber.Ctx, ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) error {
	var req AuthRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request"})
	}

	token, err := service.GetTokenByUsernameAndPassword(ctx, pool, req.Username, req.Password, cfg.JwtSecretKey)
	if err == service.ErrWrongPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Wrong password"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": err})
	}

	response := AuthResponse{
		Token: *token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
