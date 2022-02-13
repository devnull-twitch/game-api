package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/devnull-twitch/game-api/pkg/accounts"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4"
	"github.com/nicklaw5/helix"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
)

type (
	TwitchTuple struct {
		GameJWT   string
		GameToken string
		AccountID int64
	}
	CheckRequest struct {
		GameToken string
		ReplyChan chan string
	}
	CustomClaims struct {
		AccountID int64 `json:"account_id,omitempty"`
		*jwt.RegisteredClaims
	}
)

var (
	CreateTokenChan chan string       = make(chan string)
	UnlockTokenChan chan TwitchTuple  = make(chan TwitchTuple)
	GetUnlockedChan chan CheckRequest = make(chan CheckRequest)
)

func TokenProcessor() {
	waitingToken := make(map[string]string)

	for {
		select {
		case newGameToken := <-CreateTokenChan:
			waitingToken[newGameToken] = ""
		case tuple := <-UnlockTokenChan:
			waitingToken[tuple.GameToken] = tuple.GameJWT
		case checkReq := <-GetUnlockedChan:
			accessToken := waitingToken[checkReq.GameToken]
			select {
			case checkReq.ReplyChan <- accessToken:
			case <-time.After(200 * time.Millisecond):
			}
		}
	}
}

func GetSetupNewGameToken(client *helix.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		gt := randstr.Base64(12)
		select {
		case CreateTokenChan <- gt:
		case <-time.After(time.Millisecond * 200):
		}

		authURL := client.GetAuthorizationURL(&helix.AuthorizationURLParams{
			State:        gt,
			ResponseType: "code",
			Scopes:       []string{},
		})
		c.JSON(http.StatusOK, struct {
			AuthURL string `json:"auth_url"`
			Token   string `json:"wait_token"`
		}{
			AuthURL: authURL,
			Token:   gt,
		})
	}
}

func GetConfirmGameToken(client *helix.Client, accountStorage accounts.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		twUserAccess, err := client.RequestUserAccessToken(c.Query("code"))
		if err != nil {
			logrus.WithError(err).Error("unable to get user access token")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ok, resp, err := client.ValidateToken(twUserAccess.Data.AccessToken)
		if err != nil {
			logrus.WithError(err).Error("token validation error")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !ok {
			logrus.Error("token invalid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		accObj, err := accountStorage.GetByUsername(resp.Data.Login)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				if err := setupNewUser(accountStorage, resp.Data.Login); err != nil {
					logrus.WithError(err).Error("unable to rsetup new account")
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				} else {
					accObj, err = accountStorage.GetByUsername(resp.Data.Login)
					if err != nil {
						logrus.WithError(err).Error("unable to load fresh user account")
						c.AbortWithStatus(http.StatusInternalServerError)
						return
					}
				}
			} else {
				logrus.WithError(err).Error("unable to load account data")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		jwtToken, err := createToken(accObj.ID, accObj.Username)
		if err != nil {
			logrus.WithError(err).Error("error creating jwt")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		select {
		case UnlockTokenChan <- TwitchTuple{
			GameToken: c.Query("state"),
			GameJWT:   jwtToken,
		}:
		case <-time.After(time.Millisecond * 200):
		}

		c.HTML(http.StatusOK, "ok.tmpl", nil)
	}
}

func GetCheckGameToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("gametoken")
		reply := make(chan string)
		select {
		case GetUnlockedChan <- CheckRequest{GameToken: token, ReplyChan: reply}:
		case <-time.After(time.Millisecond * 200):
			logrus.Error("timeout on check game token request")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		select {
		case jwtToken := <-reply:
			if len(jwtToken) > 0 {
				c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(jwtToken))
			} else {
				c.AbortWithStatus(http.StatusNoContent)
			}
		case <-time.After(time.Millisecond * 200):
			logrus.Error("timeout on check game token response")
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func setupNewUser(accountStorage accounts.Storage, username string) error {
	if err := accountStorage.Add(username); err != nil {
		return fmt.Errorf("error adding new user: %w", err)
	}

	accObj, err := accountStorage.GetByUsername(username)
	if err != nil {
		return fmt.Errorf("error loading fresh user account: %w", err)
	}

	if err := accountStorage.AddCharacter(accObj.ID, username, "#000000"); err != nil {
		return fmt.Errorf("error adding new chracter: %w", err)
	}

	return nil
}

func createToken(userID int64, accountName string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims = &CustomClaims{
		AccountID: userID,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Subject:   accountName,
		},
	}

	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
