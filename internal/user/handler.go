package user

import (
	"crypto-payment-gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService *Service
}

func NewHandler(us *Service) *Handler {
	return &Handler{
		userService: us,
	}
}

func (h *Handler) Register(rg *gin.RouterGroup) {

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/signup", h.Signup)
		authGroup.POST("/login", h.Login)
	}

	rg.GET("/me", h.GetMe)
	rg.PATCH("/me", h.UpdateMe)

}

func (h *Handler) Signup(c *gin.Context) {

	var sr SignupRequest
	if err := c.ShouldBindJSON(&sr); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Error(err.Error()))
		return
	}

	// todo: Check Wallet Address:

	e := h.userService.SignUp(&sr)
	if e != nil {
		c.JSON(http.StatusCreated, response.Error(e.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("User created, use /login"))
}

func (h *Handler) Login(c *gin.Context) {

}

func (h *Handler) GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, MeResponse{})
}

func (h *Handler) UpdateMe(c *gin.Context) {
}
