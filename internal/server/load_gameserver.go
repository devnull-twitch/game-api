package server

import (
	"net/http"
	"os"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
)

func GetLoadGameserverHandler(s accounts.Storage, portFindFn func(string) int) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawClaim, _ := c.Get("claim")
		claims := rawClaim.(*CustomClaims)

		if !s.Exists(claims.Subject) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &ErrorRespose{
				Message: "invalid Authorization schema",
			})
			return
		}

		if c.Query("selected_char") == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no char selected"})
			return
		}

		accountData := s.Get(claims.Subject)
		requiredZone := ""
		charName := ""
		for _, char := range accountData.Characters {
			if char.Name == c.Query("selected_char") {
				if c.Query("target_zone") != "" {
					char.StartingZone = c.Query("target_zone")
				}
				requiredZone = char.StartingZone
				charName = char.Name
			}
		}
		if requiredZone == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no char zone"})
			return
		}

		port := portFindFn(requiredZone)
		if port == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{Message: "missing zone gameserver"})
			return
		}

		type GSResponse struct {
			Scene         string `json:"scene"`
			CharacterName string `json:"character_name"`
			IP            string `json:"ip"`
			Port          int    `json:"port"`
		}
		c.JSON(http.StatusOK, &GSResponse{
			Scene:         requiredZone,
			CharacterName: charName,
			IP:            os.Getenv("EXTERNAL_IP"),
			Port:          port,
		})
	}
}
