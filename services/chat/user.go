package chat

import (
	"fmt"
	"time"
)

type (
	UserId string
	user   struct {
		id       UserId
		name     string
		rooms    map[RoomId]*room
		receiver chan message
		service  *ChatService
	}
	User struct {
		Id       UserId
		Name     string
		internal *user
	}
)

func (u *User) Receiver() <-chan message {
	return u.internal.receiver
}

func (u *User) Rooms() (rooms []Room) {
	for _, room := range u.internal.rooms {
		rooms = append(rooms, Room{
			Id:   room.id,
			Name: room.name,
		})
	}
	return
}

func (u *User) JoinRoom(roomId RoomId) error {
	room, err := u.internal.service.getRoom(roomId)
	if err != nil {
		return err
	}

	room.joiner <- u.internal
	return nil
}

func (u *User) LeaveRoom(roomId RoomId) {
	room, ok := u.internal.rooms[roomId]
	if ok {
		room.leaver <- u.internal
	}
}

func (u *User) WriteMessage(roomId RoomId, msg string) error {
	room, ok := u.internal.service.rooms[roomId]
	if !ok {
		return fmt.Errorf("Room not found, roomId -> %s", roomId)
	}

	message := message{
		UserId:    u.Id,
		UserName:  u.Name,
		Timestamp: time.Now(),
		Payload:   msg,
	}

	room.forwarder <- message
	return nil
}
