package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/devnull-twitch/game-api/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TokenMW(c *gin.Context) {
	authStr := c.Request.Header.Get("Authorization")
	authParts := strings.Split(authStr, " ")

	if len(authParts) != 2 {
		c.AbortWithStatusJSON(http.StatusForbidden, &server.ErrorRespose{
			Message: "invalid Authorization header val",
		})
		return
	}

	if strings.ToLower(strings.Trim(authParts[0], " ")) != "bearer" {
		c.AbortWithStatusJSON(http.StatusForbidden, &server.ErrorRespose{
			Message: "invalid Authorization schema",
		})
		return
	}

	claims := &server.CustomClaims{}
	token, err := jwt.ParseWithClaims(strings.Trim(authParts[1], " "), claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, &server.ErrorRespose{
			Message: "invalid token",
		})
		return
	}
	if !token.Valid && claims.Valid() == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, &server.ErrorRespose{
			Message: "invalid token",
		})
		return
	}

	c.Set("claim", claims)
}
