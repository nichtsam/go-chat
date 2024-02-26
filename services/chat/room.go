package chat

type (
	RoomId string
	room   struct {
		id        RoomId
		name      string
		users     map[UserId]*user
		joiner    chan *user
		leaver    chan *user
		forwarder chan message
	}
	Room struct {
		Id   RoomId
		Name string
	}
)

func (r *room) run() {
	for {
		select {
		case user := <-r.joiner:
			r.addUser(user)
		case user := <-r.leaver:
			r.removeUser(user)
		case message := <-r.forwarder:
			r.forward(message)
		}
	}
}

func (r *room) addUser(u *user) {
	r.users[u.id] = u
}

func (r *room) removeUser(u *user) {
	delete(r.users, u.id)
}

func (r *room) forward(msg message) {
	for _, u := range r.users {
		go func(u *user) { u.receiver <- msg }(u)
	}
}
