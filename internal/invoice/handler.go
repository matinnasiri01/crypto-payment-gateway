package invoice

import (
	"crypto-payment-gateway/internal/middleware"
	"crypto-payment-gateway/pkg/jwt"
	"crypto-payment-gateway/pkg/response"
	"net/http"

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

	rg.POST("/invoice", h.CreateInvoice)

}

func (h *Handler) CreateInvoice(c *gin.Context) {

	var ci CreateInvoiceRequest
	if err := c.ShouldBindJSON(&ci); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Error(err.Error()))
		return
	}

	sErr := h.userService.Create(c.Request.Context(), middleware.UserID(c), &ci)
	if sErr != nil {
		c.JSON(http.StatusBadRequest, response.Error(sErr.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("invoice created."))
}
