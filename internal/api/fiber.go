package api

import (
	"github.com/TheFootball/internal/api/room"
	"github.com/TheFootball/internal/core/middlewares"
	"github.com/TheFootball/internal/core/redis"
	"github.com/gofiber/fiber/v2"
)

func GetFiber() {
	redis.GetRedis()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	socket := app.Group("/ws", middlewares.Ws)
	room.Init(socket)

	app.Listen(":3000")
}
