package user

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

	authRoutes := rg.Group("/auth")
	{
		authRoutes.Use(auth.GuestOnly())
		authRoutes.POST("/signup", h.Signup)
		authRoutes.POST("/login", h.Login)
	}

	userRoutes := rg.Group("", auth.Handler())
	{
		userRoutes.POST("/auth/logout", h.Logout)
		userRoutes.GET("/me", h.GetMe)
		userRoutes.PATCH("/me", h.UpdateMe)
	}
}

func (h *Handler) Signup(c *gin.Context) {

	var sr SignupRequest
	if err := c.ShouldBindJSON(&sr); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Error(err.Error()))
		return
	}

	// todo: Check Wallet Address

	e := h.userService.Signup(c.Request.Context(), &sr)
	if e != nil {
		c.JSON(http.StatusConflict, response.Error(e.Error()))
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

	user, serr := h.userService.Login(c.Request.Context(), &lr)
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
		"access_token",
		token,
		86400,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, response.Success("login success"))
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
	c.JSON(http.StatusOK, response.Success("logout success"))
}

func (h *Handler) GetMe(c *gin.Context) {

	me, sErr := h.userService.GetByID(c.Request.Context(), middleware.UserID(c))
	if sErr != nil {
		c.JSON(http.StatusNotFound, response.Error(sErr.Error()))
		return
	}

	c.JSON(http.StatusOK, me)
}

func (h *Handler) UpdateMe(c *gin.Context) {
	var ur UpdateRequest
	if err := c.ShouldBindJSON(&ur); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Error(err.Error()))
		return
	}

	serr := h.userService.Update(c.Request.Context(), middleware.UserID(c), &ur)
	if serr != nil {
		c.JSON(http.StatusBadRequest, response.Error(serr.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Done!"))
}
