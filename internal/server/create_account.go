package server

import (
	"net/http"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
)

func GetCreateAccountHandler(accountStorage accounts.Storage) gin.HandlerFunc {
	// called once on server startup

	return func(c *gin.Context) {
		// called on every request

		type registrationPayload struct {
			Username string `json:"username"`
		}
		payload := &registrationPayload{}
		if err := c.BindJSON(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Invalid registration payload",
			})
			return
		}

		if len(payload.Username) < 3 {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Missing username or username too short",
			})
			return
		}

		if err := accountStorage.Add(payload.Username); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Unable to store account",
			})
			return
		}

		c.Status(http.StatusCreated)
	}
}
