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

	authRoutes := rg.Group("/auth", auth.GuestOnly())
	{
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

// @Summary Signup User
// @Description Signup New User
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SignupRequest true "Signup Request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/signup [post]
func (h *Handler) Signup(c *gin.Context) {

	var sr SignupRequest
	if err := c.ShouldBindJSON(&sr); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Fail(err.Error()))
		return
	}

	e := h.userService.Signup(c.Request.Context(), &sr)
	if e != nil {
		c.JSON(http.StatusConflict, response.Error(e.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("User created, use /login"))
}

// @Summary Login User
// @Description Login User and Generate Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {

	var lr LoginRequest
	if err := c.ShouldBindJSON(&lr); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Fail(err.Error()))
		return
	}

	user, serr := h.userService.Login(c.Request.Context(), &lr)
	if serr != nil {
		c.JSON(http.StatusBadRequest, response.Fail(serr.Error()))
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

// @Summary Logout User
// @Description Logout Current User
// @Tags Auth
// @Produce json
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/logout [post]
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

// @Summary Get User Data
// @Description get Current User data
// @Tags User
// @Produce json
// @Success 200 {object} invoice.Response
// @Failure 404 {object} response.Response
// @Router /me [get]
func (h *Handler) GetMe(c *gin.Context) {

	me, sErr := h.userService.GetByID(c.Request.Context(), middleware.UserID(c))
	if sErr != nil {
		c.JSON(http.StatusNotFound, response.Fail(sErr.Error()))
		return
	}

	c.JSON(http.StatusOK, me)
}

// @Summary Update user
// @Description Update current user
// @Tags User
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param request body UpdateRequest true "Update Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /me [patch]
func (h *Handler) UpdateMe(c *gin.Context) {
	var ur UpdateRequest
	if err := c.ShouldBindJSON(&ur); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Fail(err.Error()))
		return
	}

	serr := h.userService.Update(c.Request.Context(), middleware.UserID(c), &ur)
	if serr != nil {
		c.JSON(http.StatusBadRequest, response.Fail(serr.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Done!"))
}
