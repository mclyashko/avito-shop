package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/service"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Missing or invalid Authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Invalid Authorization header"})
		}

		tokenString := parts[1]

		claims := &service.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return cfg.JwtSecretKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Invalid or expired token"})
		}

		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Token expired"})
		}

		return c.Next()
	}
}
