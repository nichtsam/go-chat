package chat

import "time"

type message struct {
	UserId    UserId
	UserName  string
	Timestamp time.Time
	Payload   string
}
