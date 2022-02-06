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

		characters, err := s.GetCharacters(claims.AccountID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error loading characters",
			})
			return
		}

		resp := &ChatacterPayload{}
		for _, c := range characters {
			resp.Chars = append(resp.Chars, &Chatacter{
				Name:        c.Name,
				CurrentZone: c.CurrentZone,
				BaseColor:   c.BaseColor,
			})
		}

		c.JSON(http.StatusOK, resp)
	}
}
