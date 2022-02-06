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

		if c.Query("selected_char") == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no char selected"})
			return
		}

		char, err := s.GetCharacterByName(claims.AccountID, c.Query("selected_char"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error loading character",
			})
			return
		}

		requiredZone := char.CurrentZone
		if c.Query("target_scene") != "" {
			if err := s.ChangeCurrentZone(char.ID, c.Query("target_scene")); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
					Message: "error changing character zone",
				})
				return
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
			CharacterName: char.Name,
			IP:            os.Getenv("EXTERNAL_IP"),
			Port:          port,
		})
	}
}
