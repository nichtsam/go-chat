package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"github.com/nichtsam/go-chat/view"
)

type config struct {
	Port int `env:"PORT" envDefault:"8080"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	server := gin.Default()
	server.Static("/public/", "./public/")
	server.GET("/", func(ctx *gin.Context) {
		if err := view.Home().Render(ctx, ctx.Writer); err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
	})

	log.Fatal(server.Run(addr))
}
