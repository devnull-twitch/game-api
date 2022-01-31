package server

import (
	"net/http"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetAddInventory(s accounts.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		type addItemPayload struct {
			AccountName string `json:"account"`
			Character   string `json:"character"`
			ItemID      int    `json:"item_id,string"`
		}
		payload := &addItemPayload{}
		if err := c.BindJSON(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Invalid add item payload",
			})
			return
		}

		if err := s.AddItem(payload.AccountName, payload.Character, payload.ItemID); err != nil {
			logrus.WithError(err).Error("unable to add item to char inventory")
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Something went wrong",
			})
			return
		}

		c.Status(http.StatusOK)
	}
}
