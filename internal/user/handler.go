package user

import (
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

}
func (h *Handler) Login(c *gin.Context) {

}
func (h *Handler) GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, MeResponse{})
}

func (h *Handler) UpdateMe(c *gin.Context) {
}
