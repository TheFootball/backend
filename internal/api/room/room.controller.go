package room

import (
	"fmt"
	"log"

	"github.com/TheFootball/internal/core/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type controller struct {
	roomService *service
}

func (c *controller) connectRoom(conn *websocket.Conn) {
	store := middlewares.GetStore()
	buf, err := store.Storage.Get(conn.Cookies("session_id"))
	fmt.Println(string(buf))
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
			c.roomService.disconnect(user)
			break
		}

		if err := c.roomService.messageHandler(user, &buf); err != nil {
			log.Println("err:", err)
			break
		}
	}
}

func (c *controller) createRoom(ctx *fiber.Ctx) error {
	store := middlewares.GetStore()
	sess, err := store.Get(ctx)

	if err != nil {
		panic(err)
	}

	roomId := ctx.Params("roomId")
	username := ctx.Query("username", "host")

	room, err := c.roomService.createRoomInstance(username, roomId)
	if err != nil {
		log.Fatal(err.Error())
		return ctx.Status(400).JSON(exception{Error: true, Message: err.Error()})
	}

	fmt.Print(room.RoomId)

	sess.Set("roomId", room.RoomId)
	if err = sess.Save(); err != nil {
		panic(err)
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "room created", "room": room})
}

func initController(r fiber.Router, s *service) {
	c := &controller{roomService: s}
	r.Post("/:roomId", c.createRoom)
	r.All("/:roomId", middlewares.Ws, websocket.New(c.connectRoom))
}
