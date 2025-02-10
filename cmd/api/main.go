package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/handler"
)

func main() {
	cfg := config.LoadConfig()

	_ = db.InitDB(cfg)

	app := fiber.New()

	app.Get("/", handler.HelloHandler)

	log.Println("Server is running on :8080")
	log.Fatal(app.Listen(":8080"))
}
