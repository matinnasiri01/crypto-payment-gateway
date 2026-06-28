package middleware

import (
	"crypto-payment-gateway/pkg/response"
	"net/http"

	"crypto-payment-gateway/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	jwt *jwt.Manager
}

func NewAuth(jwt *jwt.Manager) *Auth {
	return &Auth{
		jwt: jwt,
	}
}

func (a *Auth) Handler() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("login required"))
			c.Abort()
			return
		}

		claims, err := a.jwt.Parse(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("invalid token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func (a *Auth) GuestOnly() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := c.Cookie("access_token")
		if err == nil {

			_, err = a.jwt.Parse(token)
			if err == nil {

				c.JSON(
					http.StatusConflict,
					response.Error("already logged in"),
				)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
