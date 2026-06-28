package invoice

import (
	"crypto-payment-gateway/internal/middleware"
	"crypto-payment-gateway/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService *Service
	jwt         *jwt.Manager
}

func NewHandler(us *Service, j *jwt.Manager) *Handler {
	return &Handler{
		userService: us,
		jwt:         j,
	}
}

func (h *Handler) Register(rg *gin.RouterGroup, auth *middleware.Auth) {
	rg.Use(auth.Handler())
}
