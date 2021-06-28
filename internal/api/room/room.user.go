package room

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/websocket/v2"
)

type user struct {
	username       string
	userType       string
	roomId         string
	room           *redis.PubSub // Handler of redis subscribe coomand connection
	conn           *websocket.Conn
	stopListenChan chan struct{}
	messageChan    chan redis.Message
}

func newUser(username string, userType string, roomId string, conn *websocket.Conn) *user {
	u := &user{
		username:       username,
		userType:       userType,
		roomId:         roomId,
		conn:           conn,
		stopListenChan: make(chan struct{}),
		messageChan:    make(chan redis.Message),
	}
	return u
}

func (u *user) connect(ctx context.Context, r *redis.Client) error {
	if _, err := r.SAdd(ctx, u.roomId, u.username).Result(); err != nil {
		return err
	}

	pubSub := r.Subscribe(ctx, u.roomId)
	u.room = pubSub

	return nil
}

func (u *user) listen() {
	for {
		select {
		case msg, ok := <-u.room.Channel():
			if !ok {
				return
			}
			chat := chat{}
			json.Unmarshal([]byte(msg.Payload), &chat)
			u.conn.WriteJSON(chat)

		case <-u.stopListenChan:
			fmt.Println(u.username, " stop listening")
			return
		}
	}
}

func (u *user) disconnect(ctx context.Context) error {
	if u.room != nil {
		if err := u.room.Unsubscribe(ctx, u.roomId); err != nil {
			return err
		}

		if err := u.room.Close(); err != nil {
			return err
		}
	}

	u.stopListenChan <- struct{}{}
	close(u.messageChan)

	return nil
}

func (u *user) chat(ctx context.Context, r *redis.Client, msg string) error {
	chat := chat{
		Sender:  u.username,
		Message: msg,
	}

	buf, err := json.Marshal(chat)
	if err != nil {
		return err
	}

	json := string(buf)

	return r.Publish(ctx, u.roomId, json).Err()
}
