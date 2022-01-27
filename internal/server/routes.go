package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Chatacter struct {
	Name        string `json:"name"`
	BaseColor   string `json:"base_color"`
	CurrentZone string `json:"current_zone"`
}

type ErrorRespose struct {
	Message string `json:"msg"`
}

func LoginHandler(c *gin.Context) {
	type LoginPayload struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}
	type LoginResponse struct {
		Token string `json:"Token"`
	}

	payload := &LoginPayload{}
	if err := c.BindJSON(payload); err != nil {
		c.AbortWithStatusJSON(400, ErrorRespose{Message: "Yeah .. like .. what?"})
		return
	}

	if payload.Username == "test" {
		c.JSON(http.StatusOK, &LoginResponse{
			Token: "THISISATOKEN",
		})
		return
	}

	c.Status(http.StatusForbidden)
}
