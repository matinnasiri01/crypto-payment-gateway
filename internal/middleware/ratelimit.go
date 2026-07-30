package middleware

import (
	"github.com/matinnasiri01/gcpg/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"

	"golang.org/x/time/rate"
)

func RateLimiter(burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(1, burst)
	return func(c *gin.Context) {

		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, response.Error("Limit exceeded"))

		}
	}
}
