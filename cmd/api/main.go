package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/mclyashko/avito-shop/internal/config"
	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/handler"
	"github.com/mclyashko/avito-shop/internal/middleware"
	"github.com/mclyashko/avito-shop/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	database := db.InitDB(cfg)
	
	cta := &db.CoinTransferAccessorImp{
		Db: database,
	}
	ia := &db.ItemAccessorImpl{
		Db: database,
	}
	ua := &db.UserAccessorImpl{
		Db: database,
	}
	uia := &db.UserItemAccessorImpl{
		Db: database,
	}

	s := service.NewService(cfg, database)

	as := &service.AuthServiceImpl{
		Service: s,
		UserAccessor: ua,
	}
	bs := &service.BuyServiceImpl {
		Service: s,
		UserAccessor: ua,
		ItemAccessor: ia,
		UserItemAccessor: uia,
	}
	scs := &service.SendCoinsServiceImpl {
		Service: s,
		UserAccessor: ua,
		CoinTransferAccessor: cta,
	}
	uis := &service.UserInfoServiceImp {
		Service: s,
		UserAccessor: ua,
		UserItemAccessor: uia,
		CoinTransferAccessor: cta,
	}

	app := fiber.New()

	app.Post("/api/auth", func(c *fiber.Ctx) error {
		return handler.Authenticate(c, as)
	})

	app.Get("/api/info", middleware.AuthMiddleware(cfg), func(c *fiber.Ctx) error {
		return handler.GetInfo(c, uis)
	})

	app.Post("/api/sendCoin", middleware.AuthMiddleware(cfg), func(c *fiber.Ctx) error {
		return handler.SendCoinHandler(c, scs)
	})

	app.Get("/api/buy/:item", middleware.AuthMiddleware(cfg), func(c *fiber.Ctx) error {
		return handler.BuyItemHandler(c, bs)
	})

	log.Println("Server is running on :8080")
	log.Fatal(app.Listen(":8080"))
}
