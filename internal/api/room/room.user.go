package room

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rds "github.com/TheFootball/internal/core/redis"
	"github.com/TheFootball/internal/shared/exceptions"
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
	cnt, err := r.SCard(ctx, rds.MemberChannel(u.roomId)).Result()
	if err != nil {
		return err
	}

	if cnt >= 2 {
		return exceptions.ErrFullRoom
	}

	if _, err := r.SAdd(ctx, rds.MemberChannel(u.roomId), u.username).Result(); err != nil {
		return err
	}

	pubSub := r.Subscribe(ctx, u.roomId)
	u.room = pubSub
	return nil
}

func (u *user) listenChat() {
	subscribe := u.room.Channel()

	for {
		select {
		case msg, ok := <-subscribe:
			if !ok {
				return
			}
			notice := notice{}
			json.Unmarshal([]byte(msg.Payload), &notice)
			if notice.MessageType == "notice" {
				u.handleNotice(notice)
			} else {
				chat := chat{}
				json.Unmarshal([]byte(msg.Payload), &chat)
				u.conn.WriteJSON(chat)
			}
		}
	}
}

func (u *user) handleNotice(n notice) {
	r := rds.GetRedis()
	switch n.Message {
	case REMOVED:
		u.conn.WriteJSON(n)
		u.disconnect(r.Context())
	}
}

func (u *user) disconnect(ctx context.Context) error {
	fmt.Println("DISCONNECT ! ", u)
	r := rds.GetRedis()
	if u.userType == "host" {
		r.Del(r.Context(), rds.MemberChannel(u.roomId))
		if err := u.notice(r.Context(), r, REMOVED); err != nil {
			fmt.Println("Notice has error")
		}
	} else {
		r.SRem(r.Context(), rds.MemberChannel(u.roomId), u.username)
	}

	if u.room != nil {
		if err := u.room.Unsubscribe(ctx, u.roomId); err != nil {
			return err
		}

		if err := u.room.Close(); err != nil {
			return err
		}
	}

	close(u.messageChan)

	return nil
}

func (u *user) notice(ctx context.Context, r *redis.Client, msg string) error {
	notice := notice{
		Message:     msg,
		Timestamp:   time.Now().UnixNano(),
		MessageType: "notice",
	}

	buf, err := json.Marshal(notice)
	if err != nil {
		return err
	}

	json := string(buf)
	return r.Publish(ctx, u.roomId, json).Err()
}

func (u *user) chat(ctx context.Context, r *redis.Client, msg string) error {
	chat := chat{
		Sender:      u.username,
		Message:     msg,
		Timestamp:   time.Now().UnixNano(),
		MessageType: "chat",
	}

	buf, err := json.Marshal(chat)
	if err != nil {
		return err
	}

	json := string(buf)

	return r.Publish(ctx, u.roomId, json).Err()
}
