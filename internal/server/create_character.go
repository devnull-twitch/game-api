package server

import (
	"net/http"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
)

func GetCreateGameCharactersHandler(s accounts.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawClaim, _ := c.Get("claim")
		claims := rawClaim.(*CustomClaims)

		payload := &Chatacter{}
		if err := c.BindJSON(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Invalid character creation payload",
			})
			return
		}

		err := s.AddCharacter(claims.AccountID, payload.Name, payload.BaseColor)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error adding character",
			})
			return
		}

		c.Status(http.StatusCreated)
	}
}
