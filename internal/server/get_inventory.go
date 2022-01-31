package server

import (
	"net/http"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetGetInventory(s accounts.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("account") == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no char selected"})
			return
		}
		if c.Query("char") == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no char selected"})
			return
		}

		items, err := s.GetItems(c.Query("account"), c.Query("char"))
		if err != nil {
			logrus.WithError(err).Error("unable to add item to char inventory")
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Something went wrong",
			})
			return
		}

		type inventoryPayload struct {
			AccountName   string `json:"account"`
			CharacterName string `json:"character"`
			ItemIDs       []int  `json:"item_ids"`
		}
		c.JSON(http.StatusOK, &inventoryPayload{
			AccountName:   c.Query("account"),
			CharacterName: c.Query("char"),
			ItemIDs:       items,
		})
	}
}
