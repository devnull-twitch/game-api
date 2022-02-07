package server

import (
	"net/http"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetDeleteInventoryHandler(s accounts.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		type removeItemPayload struct {
			AccountName string `json:"account"`
			Character   string `json:"character"`
			ItemID      int64  `json:"item_id,string"`
			Quantity    int64  `json:"quantity,string"`
		}
		payload := &removeItemPayload{}
		if err := c.BindJSON(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Invalid remove item payload",
			})
			return
		}

		account, err := s.GetByUsername(payload.AccountName)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":        err,
				"account_name": payload.AccountName,
			}).Error("error loading account")
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error loading account",
			})
			return
		}

		char, err := s.GetCharacterByName(account.ID, payload.Character)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":      err,
				"account_id": account.ID,
				"char_name":  payload.Character,
			}).Error("error loading character")
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error loading character",
			})
			return
		}

		if err := s.RemoveItem(char.ID, payload.ItemID, payload.Quantity); err != nil {
			logrus.WithError(err).Error("unable to remove item from char inventory")
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Something went wrong",
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"character_id": char.ID,
			"item_id":      payload.ItemID,
			"quantity":     payload.Quantity,
		}).Info("removed some items")

		c.Status(http.StatusOK)
	}
}
