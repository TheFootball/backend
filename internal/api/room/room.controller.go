package room

import (
	"log"

	"github.com/TheFootball/internal/core/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type controller struct {
	roomService *service
}

func (c *controller) connectRoom(conn *websocket.Conn) {
	roomId := conn.Params("roomId")
	username := conn.Query("username")

	if err := c.roomService.checkDataValidForConnect(username, roomId); err != nil {
		conn.WriteJSON(exception{Message: err.Error(), Error: true})
		conn.Close()
		return
	}

	user, err := c.roomService.createSubscriber(username, roomId, conn)
	if err != nil {
		conn.WriteJSON(exception{Message: err.Error(), Error: true})
		conn.Close()
		return
	}

	// 메시지 받기 루프
	for {
		var (
			buf []byte
			err error
		)

		if _, buf, err = conn.ReadMessage(); err != nil {
			log.Println("read:", err)
			user.stopListenChan <- struct{}{}
			break
		}

		if err := c.roomService.messageHandler(user, &buf); err != nil {
			log.Println("err:", err)
			break
		}
	}
}

func (c *controller) createRoom(ctx *fiber.Ctx) error {
	roomId := ctx.Params("roomId")
	username := ctx.Query("username", "host")

	room, err := c.roomService.createRoomInstance(username, roomId)
	if err != nil {
		log.Fatal(err.Error())
		return ctx.Status(400).JSON(exception{Error: true, Message: err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "room created", "room": room})
}

func initController(r fiber.Router, s *service) {
	c := &controller{roomService: s}
	r.Post("/:roomId", c.createRoom)
	r.All("/:roomId", middlewares.Ws, websocket.New(c.connectRoom))
}
