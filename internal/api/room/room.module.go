package room

import (
	"github.com/gofiber/fiber/v2"
)

func Init(r fiber.Router) {
	initController(r)
}
