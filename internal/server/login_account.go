package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	*jwt.StandardClaims
}

func createToken(user string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims = &CustomClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			Subject:   user,
		},
	}

	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GetLoginHandler(accountStorage accounts.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		type LoginPayload struct {
			Username string `json:"Username"`
			Password string `json:"Password"`
		}

		payload := &LoginPayload{}
		if err := c.BindJSON(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorRespose{Message: "Invalid login data"})
			return
		}

		if !accountStorage.Exists(payload.Username) {
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorRespose{Message: "Username is not registered"})
			return
		}

		jwtToken, err := createToken(payload.Username)
		if err != nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRespose{Message: "Username is not registered"})
			return
		}
		c.String(http.StatusOK, jwtToken)
	}
}
