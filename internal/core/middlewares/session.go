package middlewares

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func GetStore() *session.Store {
	if store != nil {
		return store
	}

	store = session.New()
	return store
}
