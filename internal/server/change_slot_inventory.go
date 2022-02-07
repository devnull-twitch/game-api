package server

import (
	"net/http"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetSlotChangeInventoryHandler(s accounts.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		type changeItemPayload struct {
			AccountName   string `json:"account"`
			Character     string `json:"character"`
			SlotID        int64  `json:"slot_id,string"`
			ItemIDToSlot  int64  `json:"set_in_slot,string"`
			ItemIDToUnlot int64  `json:"remove_from_slot,string"`
		}
		payload := &changeItemPayload{}
		if err := c.BindJSON(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Invalid change item payload",
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

		if payload.ItemIDToUnlot > 0 {
			if err := s.UnslotEquipment(char.ID, payload.SlotID, payload.ItemIDToUnlot); err != nil {
				logrus.WithError(err).Error("unable to unequip item")
				c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
					Message: "Something went wrong",
				})
				return
			}
		}
		if payload.ItemIDToSlot > 0 {
			if err := s.SlotEquipment(char.ID, payload.SlotID, payload.ItemIDToSlot); err != nil {
				logrus.WithError(err).Error("unable to slot item")
				c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
					Message: "Something went wrong",
				})
				return
			}
		}

		c.Status(http.StatusOK)
	}
}
