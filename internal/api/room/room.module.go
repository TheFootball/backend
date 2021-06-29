package room

import (
	"github.com/TheFootball/internal/core/redis"
	"github.com/gofiber/fiber/v2"
)

func Init(r fiber.Router) {
	s := &service{rdb: redis.GetRedis()}
	initController(r, s)
}
