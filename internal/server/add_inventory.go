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
			ItemID      int64  `json:"item_id,string"`
			Quantity    int64  `json:"quantity,string"`
		}
		payload := &addItemPayload{}
		if err := c.BindJSON(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Invalid add item payload",
			})
			return
		}

		account, err := s.GetByUsername(payload.AccountName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error loading account",
			})
			return
		}

		char, err := s.GetCharacterByName(account.ID, payload.Character)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error loading character",
			})
			return
		}

		if err := s.AddItem(&accounts.InventorySlot{
			CharacterID: char.ID,
			ItemID:      payload.ItemID,
			Quantity:    payload.Quantity,
		}); err != nil {
			logrus.WithError(err).Error("unable to add item to char inventory")
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Something went wrong",
			})
			return
		}

		c.Status(http.StatusOK)
	}
}
