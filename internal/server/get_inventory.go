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
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no account provided"})
			return
		}
		if c.Query("char") == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no char provided"})
			return
		}

		account, err := s.GetByUsername(c.Query("account"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error loading account",
			})
			return
		}

		char, err := s.GetCharacterByName(account.ID, c.Query("char"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrorRespose{
				Message: "error loading character",
			})
			return
		}

		items, err := s.GetItems(char.ID)
		if err != nil {
			logrus.WithError(err).Error("unable to add item to char inventory")
			c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{
				Message: "Something went wrong",
			})
			return
		}

		type (
			inventorySlotPayload struct {
				ItemID   int64 `json:"item_id"`
				Quantity int64 `json:"quantity"`
				SlotID   int64 `json:"slot_id"`
			}
			inventoryPayload struct {
				AccountName   string                 `json:"account"`
				CharacterName string                 `json:"character"`
				ItemIDs       []inventorySlotPayload `json:"items"`
			}
		)

		slotReturnData := make([]inventorySlotPayload, len(items))
		for i, slotData := range items {
			slotReturnData[i] = inventorySlotPayload{
				ItemID:   slotData.ItemID,
				Quantity: slotData.Quantity,
				SlotID:   slotData.SlotID,
			}
		}

		c.JSON(http.StatusOK, &inventoryPayload{
			AccountName:   c.Query("account"),
			CharacterName: c.Query("char"),
			ItemIDs:       slotReturnData,
		})
	}
}
