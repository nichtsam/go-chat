package websocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/nichtsam/go-chat/services/chat"
	"github.com/nichtsam/go-chat/websocket/connection"
)

var websocketUpgrader = websocket.Upgrader{}

type (
	serviceHooker    func(*client, *sync.WaitGroup, <-chan struct{})
	websocketHandler struct {
		chat           *chat.ChatService
		serviceHookers []serviceHooker
	}
)

func NewWebsocketHandler(chat *chat.ChatService) *websocketHandler {
	ws := &websocketHandler{chat: chat}
	ws.Setup()

	return ws
}

func (h *websocketHandler) Setup() {
	h.serviceHookers = append(h.serviceHookers, h.hookChat)
}

func (h *websocketHandler) HandleRequest(ctx *gin.Context) {
	conn, err := upgradeWS(ctx)
	if err != nil {
		log.Println("Error: Failed to upgrade to websocket -> ", err)
		return
	} else {
		defer func() { conn.Close() }()
	}

	client := newClient(conn, h)
	client.run()
}

func (wsh *websocketHandler) routeIncomingIntent(
	c *client,
	intent connection.IncomingIntent,
) error {
	h, ok := IntentHandlers[intent.IntentType]
	if !ok {
		return fmt.Errorf("Invalid intent \"%s\"", intent.IntentType)
	}

	return h(c, intent)
}

func (wsh *websocketHandler) routeOutgoingEvent(c *client, event connection.OutgoingEvent) error {
	h, ok := EventHandlers[event.EventType]
	if !ok {
		return fmt.Errorf("Unsupported event \"%s\"", event.EventType)
	}

	return h(c, event)
}

func upgradeWS(ctx *gin.Context) (*websocket.Conn, error) {
	w, r := ctx.Writer, ctx.Request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
