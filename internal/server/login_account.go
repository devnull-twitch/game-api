package server

import (
	"net/http"
	"os"
	"time"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type CustomClaims struct {
	AccountID int64 `json:"account_id,omitempty"`
	*jwt.StandardClaims
}

func createToken(userID int64, accountName string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims = &CustomClaims{
		AccountID: userID,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			Subject:   accountName,
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

		acc, err := accountStorage.GetByUsername(payload.Username)
		if err != nil {
			logrus.WithError(err).Error("unmable to load account data by username")
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRespose{Message: "database error"})
			return
		}

		// TODO add password hash check ?!

		jwtToken, err := createToken(acc.ID, acc.Username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRespose{Message: "Username is not registered"})
			return
		}
		c.String(http.StatusOK, jwtToken)
	}
}
