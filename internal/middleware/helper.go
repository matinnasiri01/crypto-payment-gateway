package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserID(c *gin.Context) uuid.UUID {
	return c.MustGet("user_id").(uuid.UUID)
}

// todo add redis and make rate limits for ips
