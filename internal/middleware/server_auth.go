package middleware

import (
	"net/http"
	"os"

	"github.com/devnull-twitch/game-api/internal/server"
	"github.com/gin-gonic/gin"
)

func SrverAuthMW(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, &server.ErrorRespose{
			Message: "invalid Authorization header val",
		})
		return
	}

	if username != "gameserver" {
		c.AbortWithStatusJSON(http.StatusForbidden, &server.ErrorRespose{
			Message: "invalid Authorization header val",
		})
		return
	}

	if password != os.Getenv("GS_AUTH_PASSWORD") {
		c.AbortWithStatusJSON(http.StatusForbidden, &server.ErrorRespose{
			Message: "invalid Authorization header val",
		})
		return
	}
}
