package main

import (
	"os"

	"github.com/devnull-twitch/game-api/internal/middleware"
	"github.com/devnull-twitch/game-api/internal/server"
	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.yaml")
	s := accounts.NewStorage()

	r := gin.Default()
	r.POST("/account", server.GetCreateAccountHandler(s))
	r.POST("/account/login", server.GetLoginHandler(s))
	r.GET("/game/characters", middleware.TokenMW, server.GetLoadGameCharactersHandler(s))
	r.POST("/game/characters", middleware.TokenMW, server.GetCreateGameCharactersHandler(s))
	r.POST("/game/play", middleware.TokenMW, server.GetGameserverHandler)
	r.POST("/game/change_scene", middleware.TokenMW, server.GetGameserverChangeHandler)

	r.Run(os.Getenv("WEBSERVER_BIND"))
}
