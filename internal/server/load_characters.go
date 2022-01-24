package server

import (
	"net/http"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
)

func GetLoadGameCharactersHandler(s accounts.Storage) gin.HandlerFunc {
	type ChatacterPayload struct {
		Chars []*Chatacter `json:"chars"`
	}

	return func(c *gin.Context) {
		rawClaim, _ := c.Get("claim")
		claims := rawClaim.(*CustomClaims)

		if !s.Exists(claims.Subject) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &ErrorRespose{
				Message: "invalid Authorization schema",
			})
			return
		}

		acc := s.Get(claims.Subject)
		resp := &ChatacterPayload{}

		for _, c := range acc.Characters {
			resp.Chars = append(resp.Chars, &Chatacter{
				Name:        c.Name,
				CurrentZone: c.StartingZone,
				BaseColor:   c.BaseColor,
			})
		}

		c.JSON(http.StatusOK, resp)
	}
}
