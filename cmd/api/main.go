package main

import (
	"github.com/devnull-twitch/game-api/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/game/login", server.LoginHandler)
	r.GET("/game/characters", server.CharacterLoaderHandler)
	r.POST("/game/play", server.GetGameserverHandler)
	r.POST("/game/change_scene", server.GetGameserverChangeHandler)

	r.Run(":8082")
}
