package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/devnull-twitch/game-api/internal/middleware"
	"github.com/devnull-twitch/game-api/internal/server"
	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/devnull-twitch/game-api/pkg/k8s"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	godotenv.Load(".env.yaml")

	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to database: %w", err))
	}

	s := accounts.NewStorage(conn)

	var portFinder func(string) int
	if os.Getenv("USE_K8S") != "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		portFinder = k8s.GetPortFinder(clientset)
	} else {
		portFinder = func(zone string) int {
			switch zone {
			case "starting_zone":
				return 50125
			case "world_1":
				return 50126
			}

			logrus.WithError(fmt.Errorf("unknown zone %s", zone)).Warn("missing zone port")
			return 0
		}
	}

	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://devnullga.me"
	}
	client, err := helix.NewClient(&helix.Options{
		ClientID:       os.Getenv("TW_CLIENTID"),
		AppAccessToken: os.Getenv("TW_APP_ACCESS"),
		RedirectURI:    fmt.Sprintf("%s/rpg/twitch/confirm", apiBaseURL),
	})
	if err != nil {
		log.Fatal("unable to create twitch api client")
	}

	go server.TokenProcessor()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.LoadHTMLGlob("templates/*")

	group := r.Group("/rpg")
	{
		group.POST("/twitch/start", server.GetSetupNewGameToken(client))
		group.GET("/twitch/confirm", server.GetConfirmGameToken(client, s))
		group.GET("/twitch/check", server.GetCheckGameToken())
		group.GET("/game/characters", middleware.TokenMW, server.GetLoadGameCharactersHandler(s))
		group.POST("/game/play", middleware.TokenMW, server.GetLoadGameserverHandler(s, portFinder))
		group.POST("/character/inventory", middleware.SrverAuthMW, server.GetAddInventory(s))
		group.GET("/character/inventory", middleware.SrverAuthMW, server.GetGetInventory(s))
		group.POST("/character/inventory/slot_change", middleware.SrverAuthMW, server.GetSlotChangeInventoryHandler(s))
		group.DELETE("/character/inventory", middleware.SrverAuthMW, server.GetDeleteInventoryHandler(s))
		group.PUT("/character/inventory", middleware.SrverAuthMW, server.GetUpdateInventoryHandler(s))
	}

	r.Run(os.Getenv("WEBSERVER_BIND"))
}
