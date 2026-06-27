package user

import (
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

	e := h.userService.Signup(&sr)
	if e != nil {
		c.JSON(http.StatusCreated, response.Error(e.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("User created, use /login"))
}

func (h *Handler) Login(c *gin.Context) {

	var lr LoginRequest
	if err := c.ShouldBindJSON(&lr); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Error(err.Error()))
		return
	}

	user, serr := h.userService.Login(&lr)
	if serr != nil {
		c.JSON(http.StatusBadRequest, response.Error(serr.Error()))
		return
	}

	token, tErr := h.jwt.Generate(user.ID)
	if tErr != nil {
		c.JSON(http.StatusInternalServerError, response.Error(tErr.Error()))
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"Authorization",
		token,
		86400,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, response.Success("login success"))
}

func (h *Handler) GetMe(c *gin.Context) {

	token, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("first /login"))
		return
	}

	id, jErr := h.jwt.Parse(token)
	if jErr != nil {
		c.JSON(http.StatusInternalServerError, response.Error("internal error"))
		return
	}

	me, sErr := h.userService.GetByID(id.UserID)
	if sErr != nil {
		c.JSON(http.StatusNotFound, response.Error(sErr.Error()))
		return
	}

	c.JSON(http.StatusOK, me)
}

func (h *Handler) UpdateMe(c *gin.Context) {
}
