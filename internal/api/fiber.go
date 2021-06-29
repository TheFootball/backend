package api

import (
	"github.com/TheFootball/internal/api/room"
	"github.com/TheFootball/internal/core/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func GetFiber() {
	redis.GetRedis()

	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	socket := app.Group("/ws")
	room.Init(socket)

	app.Listen(":3000")
}
