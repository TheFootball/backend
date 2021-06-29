package room

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/TheFootball/internal/core/middlewares"
	"github.com/TheFootball/internal/core/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type controller struct{}

func (c *controller) connectRoom(conn *websocket.Conn) {
	var (
		_   int
		msg []byte
		err error
	)
	r := redis.GetRedis()

	roomId := conn.Params("roomId")

	if _, err = r.Get(r.Context(), roomId).Result(); redis.IsNil(err) {
		conn.WriteJSON(notice{Message: "invalid roomId", MessageType: "notice"})
		conn.Close()
		return
	}

	username := conn.Query("username")
	if username == "" {
		conn.WriteJSON(fiber.Map{"message": "username is required", "error": true})
		conn.Close()
		return
	}

	user := newUser(username, "guest", roomId, conn)

	// 채널 퍼블리싱
	if err = user.connect(r.Context(), r); err != nil {
		conn.WriteJSON(fiber.Map{"message": "error occurs", "error": true})
		conn.Close()
		return
	}
	go user.listen()

	conn.WriteJSON(fiber.Map{"message": fmt.Sprintf("Hello, %s! Welcome to %s channel.", username, roomId)})

	// 메시지 받기 루프
	for {
		if _, msg, err = conn.ReadMessage(); err != nil {
			log.Println("read:", err)
			user.stopListenChan <- struct{}{}
			break
		}

		if err := user.chat(r.Context(), r, string(msg)); err != nil {
			log.Println("err:", err)
			break
		}
	}
}

func (c *controller) createRoom(ctx *fiber.Ctx) error {
	roomId := ctx.Params("roomId")
	username := ctx.Query("username", "host")

	r := redis.GetRedis()

	if _, err := r.Get(r.Context(), roomId).Result(); !redis.IsNil(err) {
		return ctx.Status(400).JSON(fiber.Map{"message": "already exist room", "error": true})
	}

	room := room{RoomId: roomId, Host: username, Ongoing: false}
	buf, err := json.Marshal(room)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"message": "error occurs", "error": true})
	}

	if err := r.Set(r.Context(), roomId, buf, 0).Err(); err != nil {
		log.Fatal(err)
		return ctx.Status(400).JSON(fiber.Map{"message": "error occurs", "error": true})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "room created", "room": room})
}

func initController(r fiber.Router) {
	c := &controller{}
	r.Post("/:roomId", c.createRoom)
	r.All("/:roomId", middlewares.Ws, websocket.New(c.connectRoom))
}
