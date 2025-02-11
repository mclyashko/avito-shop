package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/handler"
	"github.com/mclyashko/avito-shop/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	pool := db.InitDB(cfg)

	app := fiber.New()

	app.Post("/api/auth", func(c *fiber.Ctx) error {
		return handler.Authenticate(c, cfg, pool)
	})

	app.Get("/api/info", middleware.AuthMiddleware(cfg), func(c *fiber.Ctx) error {
		return handler.GetInfo(c, pool)
	})

	app.Post("/api/sendCoin", middleware.AuthMiddleware(cfg), func (c *fiber.Ctx) error  {
		return handler.SendCoinHandler(c, pool)
	})

	app.Get("/api/buy/:item", middleware.AuthMiddleware(cfg), func  (c *fiber.Ctx) error {
		return handler.BuyItemHandler(c, pool)
	})

	log.Println("Server is running on :8080")
	log.Fatal(app.Listen(":8080"))
}
