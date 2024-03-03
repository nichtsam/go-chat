package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"github.com/nichtsam/go-chat/services/chat"
	. "github.com/nichtsam/go-chat/websocket/connection"
)

// Intent

type (
	intentHandler  func(*client, IncomingIntent) error
	intentHandlers map[IntentType]intentHandler
)

func handleSendMessage(c *client, i IncomingIntent) error {
	data := struct {
		RoomId  chat.RoomId `json:"room_id"`
		Message string      `json:"message"`
	}{}
	if err := json.Unmarshal(i.Raw, &data); err != nil {
		return err
	}

	if hasRoom := c.user.HasRoom(data.RoomId); !hasRoom {
		if err := c.user.JoinRoom(data.RoomId); err != nil {
			log.Println("Error: Room should exist -> ", err)
		}
	}

	if err := c.user.WriteMessage(data.RoomId, data.Message); err != nil {
		return err
	}

	return nil
}

func getIntentHandlers() intentHandlers {
	handlers := make(intentHandlers)
	handlers[Intent.SEND_MESSAGE] = handleSendMessage
	return handlers
}

var IntentHandlers = getIntentHandlers()

// Event

type (
	eventHandler  func(*client, OutgoingEvent) error
	eventHandlers map[EventType]eventHandler
)

func handleWriteMessage(c *client, e OutgoingEvent) error {
	msg, ok := e.Payload.(string)
	if !ok {
		return fmt.Errorf("Invalid payload, expected string, received %T", e.Payload)
	}

	return c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func getEventHandlers() eventHandlers {
	handlers := make(eventHandlers)
	handlers[Event.WRITE_MESSAGE] = handleWriteMessage
	return handlers
}

var EventHandlers = getEventHandlers()
