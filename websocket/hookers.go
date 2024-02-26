package websocket

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nichtsam/go-chat/view/components"
	"github.com/nichtsam/go-chat/websocket/connection"
)

func (h *websocketHandler) hookChat(c *client, wg *sync.WaitGroup, done <-chan struct{}) {
	go func() {
		defer func() {
			log.Println("Info: Stopped chat hook.")
			wg.Done()
		}()

		wg.Add(1)
		for {
			select {
			case msg := <-c.user.Receiver():
				msgString := fmt.Sprintf(
					"%v %s: %s",
					msg.Timestamp.Format(time.DateTime),
					msg.UserName,
					msg.Payload,
				)

				buf := new(bytes.Buffer)
				if err := components.ChatMessage(msgString).Render(context.Background(), buf); err != nil {
					log.Println("Error: Failed to render chat message htmx -> ", err)
					continue
				}

				htmx := buf.String()

				event := connection.OutgoingEvent{
					EventType: connection.Event.WRITE_MESSAGE,
					Payload:   htmx,
				}

				c.egress <- event

			case <-done:
				return
			}
		}
	}()
}
