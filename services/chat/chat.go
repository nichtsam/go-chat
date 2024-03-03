package chat

import (
	"fmt"

	"github.com/google/uuid"
)

type ChatService struct {
	rooms map[RoomId]*room
	users map[UserId]User
}

func NewChatService() *ChatService {
	return &ChatService{
		rooms: make(map[RoomId]*room),
	}
}

func (s *ChatService) CreateRoom(name string) {
	room := &room{
		id:        RoomId(uuid.NewString()),
		name:      name,
		users:     make(map[UserId]*user),
		joiner:    make(chan *user),
		leaver:    make(chan *user),
		forwarder: make(chan message),
	}
	go room.run()

	s.rooms[room.id] = room
}

func (s *ChatService) CreateUser(name string) *User {
	internal := &user{
		id:       UserId(uuid.NewString()),
		name:     name,
		receiver: make(chan message),
		service:  s,
	}

	return &User{
		Id:       internal.id,
		Name:     internal.name,
		internal: internal,
	}
}

func (s *ChatService) DestoryUser(u *User) {
	for _, room := range u.Rooms() {
		u.LeaveRoom(room.Id)
	}
	delete(s.users, u.Id)
}

func (s *ChatService) Rooms() (rooms []Room) {
	for _, room := range s.rooms {
		rooms = append(rooms, Room{
			Id:   room.id,
			Name: room.name,
		})
	}
	return
}

func (s *ChatService) Users() (users []User) {
	for _, user := range s.users {
		users = append(users, user)
	}
	return
}

func (s *ChatService) getRoom(roomId RoomId) (*room, error) {
	room, ok := s.rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("Room not found, roomId -> %s", roomId)
	}

	return room, nil
}

func (s *ChatService) GetRoom(roomId RoomId) (*Room, error) {
	room, ok := s.rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("Room not found, roomId -> %s", roomId)
	}

	return &Room{
		Id:   room.id,
		Name: room.name,
	}, nil
}
