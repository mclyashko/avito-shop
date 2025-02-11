package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mclyashko/avito-shop/internal/service"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func Authenticate(c *fiber.Ctx, s service.AuthService) error {
	ctx := c.Context()

	var req authRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse auth request body, error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid request"})
	}

	token, err := s.GetTokenByUsernameAndPassword(ctx, req.Username, req.Password)
	if err == service.ErrWrongPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Wrong password"})
	}
	if err != nil {
		log.Printf("Failed to get token for username: %v, error: %v", req.Username, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Cant get token"})
	}

	response := authResponse{
		Token: *token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
