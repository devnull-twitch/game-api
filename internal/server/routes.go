package server

import (
	"context"
	"log"
	"net/http"

	"github.com/devnull-twitch/gameserver-manager/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func GetGameserverHandler(c *gin.Context) {
	playerChars := []*Chatacter{
		{Name: "TestCharacter", CurrentZone: "overworld"},
	}

	if c.Query("selected_char") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no char selected"})
	}

	var targetChar *Chatacter
	for _, playerChar := range playerChars {
		if playerChar.Name == c.Query("selected_char") {
			targetChar = playerChar
			break
		}
	}

	if targetChar == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no valid char name given"})
	}

	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewGameserverManagerClient(conn)

	response, err := client.GetGameserver(context.Background(), &proto.GetRequest{
		Zone: targetChar.CurrentZone,
	})
	if err != nil {
		log.Fatalf("rpc error: %v", err)
	}

	type GSResponse struct {
		Scene string `json:"scene"`
		IP    string `json:"ip"`
		Port  int    `json:"port"`
	}
	c.JSON(200, &GSResponse{IP: "127.0.0.1", Port: int(response.GetGsPort()), Scene: "overworld"})
}

func GetGameserverChangeHandler(c *gin.Context) {
	if c.Query("target_scene") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorRespose{Message: "no char selected"})
	}

	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewGameserverManagerClient(conn)

	response, err := client.GetGameserver(context.Background(), &proto.GetRequest{
		Zone: c.Query("target_scene"),
	})
	if err != nil {
		log.Fatalf("rpc error: %v", err)
	}

	type GSResponse struct {
		Scene string `json:"scene"`
		IP    string `json:"ip"`
		Port  int    `json:"port"`
	}
	c.JSON(200, &GSResponse{IP: "127.0.0.1", Port: int(response.GetGsPort()), Scene: c.Query("target_scene")})
}
