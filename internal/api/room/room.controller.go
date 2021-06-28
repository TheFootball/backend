package room

import (
	"fmt"
	"log"

	"github.com/TheFootball/internal/core/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type controller struct {
}

func (c *controller) connectRoom(conn *websocket.Conn) {
	var (
		_   int
		msg []byte
		err error
	)

	roomId := conn.Params("roomId")

	username := conn.Query("username")
	if username == "" {
		conn.WriteJSON(fiber.Map{"message": "username is required", "error": true})
		conn.Close()
		return
	}

	user := newUser(username, "guest", roomId, conn)
	r := redis.GetRedis()
	user.connect(r.Context(), r)
	go user.listen()

	conn.WriteJSON(fiber.Map{"message": fmt.Sprintf("Hello, %s! Welcome to %s channel.", username, roomId)})

	for {
		if _, msg, err = conn.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		if err := user.chat(r.Context(), r, string(msg)); err != nil {
			log.Println("err", err)
			break
		}
	}
}

func (c *controller) createRoom(conn *websocket.Conn) {
	// 이제 이거 끝내야지~
}

func initController(r fiber.Router) {
	c := &controller{}
	r.Get("/:roomId", websocket.New(c.connectRoom))
	r.Post("/:roomId", websocket.New(c.createRoom))
}
