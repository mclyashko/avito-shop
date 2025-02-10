package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/handler"
	"github.com/mclyashko/avito-shop/internal/middleware"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()

	pool := db.InitDB(cfg)

	app := fiber.New()

	app.Get("/", handler.HelloHandler)
	app.Get("/authed", middleware.AuthMiddleware(cfg), handler.HelloHandler)

	app.Post("/api/auth", func(c *fiber.Ctx) error {
		return handler.Authenticate(c, ctx, cfg, pool)
	})

	log.Println("Server is running on :8080")
	log.Fatal(app.Listen(":8080"))
}
