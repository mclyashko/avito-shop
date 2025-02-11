package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/service"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func Authenticate(c *fiber.Ctx, cfg *config.Config, pool *pgxpool.Pool) error {
	ctx := c.Context()

	var req authRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse auth request body, error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request"})
	}

	token, err := service.GetTokenByUsernameAndPassword(ctx, cfg, pool, req.Username, req.Password, cfg.JwtSecretKey)
	if err == service.ErrWrongPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Wrong password"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Cant get token"})
	}

	response := authResponse{
		Token: *token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
