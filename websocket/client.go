package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/nichtsam/go-chat/services/chat"
	"github.com/nichtsam/go-chat/websocket/connection"
)

type client struct {
	conn    *websocket.Conn
	handler *websocketHandler
	user    *chat.User
	egress  chan connection.OutgoingEvent
}

func newClient(conn *websocket.Conn, handler *websocketHandler) *client {
	user := handler.chat.CreateUser("anonymous")

	return &client{
		conn:    conn,
		handler: handler,
		user:    user,
		egress:  make(chan connection.OutgoingEvent),
	}
}

func (client *client) run() {
	defer func() { client.tidy() }()

	servicesWg := &sync.WaitGroup{}
	servicesDone := make(chan struct{})

	go client.startWriting()
	client.hookServices(servicesWg, servicesDone)
	client.startReading()

	servicesDone <- struct{}{}
	servicesWg.Wait()
}

func (c *client) tidy() {
	close(c.egress)
	c.handler.chat.DestoryUser(c.user)
}

func (c *client) hookServices(wg *sync.WaitGroup, done <-chan struct{}) {
	for _, h := range c.handler.serviceHookers {
		h(c, wg, done)
	}
}

func (c *client) startReading() {
	defer func() {
		log.Println("Info: Stopped reading.")
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Error: Connection unexpectedly closed -> ", err)
			} else {
				log.Println("Info: Connection closed")
			}

			return
		}

		event := connection.IncomingIntent{Raw: msg}
		if err := json.Unmarshal(msg, &event); err != nil {
			log.Println("Error: Invalid event -> ", err)
		}

		if err := c.handler.routeIncomingIntent(c, event); err != nil {
			log.Println("Error: Failed to route incoming intent -> ", err)
		}
	}
}

func (c *client) startWriting() {
	defer func() { log.Println("Info: Stopped writing.") }()
	stopPingPong := c.pingPong()
	defer func() { stopPingPong <- struct{}{} }()

	for event := range c.egress {
		if err := c.handler.routeOutgoingEvent(c, event); err != nil {
			log.Println("Error: Failed to route outgoing event -> ", err)
		}
	}

	log.Println("Info: Egress channel closed")
}

var (
	PONG_WAIT     = time.Second * 10
	PING_INTERVAL = PONG_WAIT * 9 / 10
)

func (c *client) pingPong() chan<- struct{} {
	if err := c.conn.SetReadDeadline(time.Now().Add(PONG_WAIT)); err != nil {
		log.Println("Error: Failed to set read deadline -> ", err)
	}
	c.conn.SetPongHandler(c.pongHandler)
	stop := make(chan struct{})
	go c.startPinging(stop)
	return stop
}

func (c *client) pongHandler(_ string) error {
	log.Println("Info: Pong!")
	return c.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
}

func (c *client) startPinging(stop <-chan struct{}) {
	ticker := time.NewTicker(PING_INTERVAL)
	defer func() {
		ticker.Stop()
		log.Println("Info: Stopped pinging.")
	}()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			log.Println("pinging client...")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error: Failed to ping the client -> ", err)
			}
		}
	}
}
