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

		if !s.Exists(claims.Subject) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &ErrorRespose{
				Message: "invalid Authorization schema",
			})
			return
		}

		payload := &Chatacter{}
		if err := c.BindJSON(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Invalid character creation payload",
			})
			return
		}

		acc := s.Get(claims.Subject)
		acc.Characters = append(acc.Characters, &accounts.GameCharacter{
			StartingZone: "starting_zone",
			Name:         payload.Name,
			BaseColor:    payload.BaseColor,
		})

		c.Status(http.StatusCreated)
	}
}
