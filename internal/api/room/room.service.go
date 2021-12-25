package room

import (
	"encoding/json"

	rds "github.com/TheFootball/internal/core/redis"
	"github.com/TheFootball/internal/shared/exceptions"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/websocket/v2"
)

type service struct {
	rdb *redis.Client
}

func (s *service) createRoomInstance(username string, roomId string) (*room, error) {
	if _, err := s.rdb.Get(s.rdb.Context(), roomId).Result(); !rds.IsNil(err) {
		return nil, exceptions.ErrDuplicated(roomId)
	}

	room := room{RoomId: roomId, Host: username, Ongoing: false}
	buf, err := json.Marshal(room)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(s.rdb.Context(), roomId, buf, 0).Err(); err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *service) checkDataValidForConnect(username string, roomId string) error {
	if username == "" || roomId == "" {
		return exceptions.ErrInvalidDTO("username", "or roomId")
	}

	if _, err := s.rdb.Get(s.rdb.Context(), roomId).Result(); rds.IsNil(err) {
		return exceptions.ErrInvalidDTO("roomId")
	}

	return nil
}

func (s *service) createSubscriber(username string, roomId string, conn *websocket.Conn) (*user, error) {
	user := newUser(username, "guest", roomId, conn)

	// 채널 구독
	if err := user.connect(s.rdb.Context(), s.rdb); err != nil {
		return nil, err
	}
	go user.listenChat()

	return user, nil
}

func (s *service) disconnect(u *user) error {
	return u.disconnect(s.rdb.Context())
}

func (s *service) messageHandler(u *user, buf *[]byte) error {
	msg := message{}

	if err := json.Unmarshal(*buf, &msg); err != nil {
		return err
	}

	if msg.Type == "command" && u.userType == "host" {
		// 게임 시작하기
		return nil
	} else {
		return u.chat(s.rdb.Context(), s.rdb, msg.Message)
	}
}
