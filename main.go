package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"github.com/nichtsam/go-chat/services/chat"
	"github.com/nichtsam/go-chat/view/pages"
	"github.com/nichtsam/go-chat/websocket"
)

type config struct {
	Port int `env:"PORT" envDefault:"8080"`
}

func main() {
	// config
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	addr := fmt.Sprintf(":%d", cfg.Port)

	// prepare
	server := gin.Default()
	_ = server.SetTrustedProxies(nil)

	chatService := chat.NewChatService()
	chatService.CreateRoom("default")

	handler := &handler{chat: chatService}
	ws := websocket.NewWebsocketHandler(chatService)

	// routes
	server.Static("/public/", "./public/")
	server.GET("/", handler.handleHome)
	server.GET("/chat/:roomId", handler.handleChatRoom)
	server.GET("/ws", ws.HandleRequest) // this gets logged when connection ends

	// run
	log.Fatal(server.Run(addr))
}

type handler struct {
	chat *chat.ChatService
}

func (h *handler) handleHome(ctx *gin.Context) {
	if err := pages.Home(h.chat.Rooms()).Render(ctx, ctx.Writer); err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}

func (h *handler) handleChatRoom(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	room, err := h.chat.GetRoom(chat.RoomId(roomId))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	}
	if err := pages.Chat(*room).Render(ctx, ctx.Writer); err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}
