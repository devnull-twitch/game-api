package main

import (
	"fmt"
	"os"

	"github.com/devnull-twitch/game-api/internal/middleware"
	"github.com/devnull-twitch/game-api/internal/server"
	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/devnull-twitch/game-api/pkg/k8s"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	godotenv.Load(".env.yaml")
	s := accounts.NewStorage()

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

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.POST("/account", server.GetCreateAccountHandler(s))
	r.POST("/account/login", server.GetLoginHandler(s))
	r.GET("/game/characters", middleware.TokenMW, server.GetLoadGameCharactersHandler(s))
	r.POST("/game/characters", middleware.TokenMW, server.GetCreateGameCharactersHandler(s))
	r.POST("/game/play", middleware.TokenMW, server.GetLoadGameserverHandler(s, portFinder))
	r.POST("/character/inventory", middleware.SrverAuthMW, server.GetAddInventory(s))
	r.GET("/character/inventory", middleware.SrverAuthMW, server.GetGetInventory(s))

	r.Run(os.Getenv("WEBSERVER_BIND"))
}
