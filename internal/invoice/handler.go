package invoice

import (
	"crypto-payment-gateway/internal/middleware"
	"crypto-payment-gateway/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	invoiceService *Service
}

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

func NewHandler(us *Service) *Handler {
	return &Handler{
		invoiceService: us,
	}
}

func (h *Handler) Register(rg *gin.RouterGroup, auth *middleware.Auth) {

	invoices := rg.Group("/invoices", auth.Handler())
	{
		invoices.POST("", h.Create)
		invoices.GET("", h.List)
		invoices.GET("/:id", h.GetByID)
		invoices.PATCH("/:id", h.Update)
		invoices.DELETE("/:id", h.Delete)
	}

	rg.GET("/pay/:id", h.Pay)

}

func (h *Handler) Create(c *gin.Context) {

	var ci CreateRequest
	if err := c.ShouldBindJSON(&ci); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Error(err.Error()))
		return
	}

	sErr := h.invoiceService.Create(c.Request.Context(), middleware.UserID(c), &ci)
	if sErr != nil {
		c.JSON(http.StatusBadRequest, response.Error(sErr.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("invoice created."))
}

func (h *Handler) List(c *gin.Context) {

	page := 1
	limit := DefaultLimit

	if p := c.Query("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}

	if l := c.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	list, err := h.invoiceService.List(c.Request.Context(), middleware.UserID(c), page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	// todo Adding page limit and total record count values to the response
	c.JSON(http.StatusOK, list)

}

func (h *Handler) GetByID(c *gin.Context) {

	id, uErr := uuid.Parse(c.Param("id"))
	if uErr != nil {
		c.JSON(http.StatusBadRequest, response.Error(uErr.Error()))
		return
	}

	res, err := h.invoiceService.GetByID(c.Request.Context(), id, middleware.UserID(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusFound, res)

}

func (h *Handler) Update(c *gin.Context) {

	var ci UpdateRequest
	if err := c.ShouldBindJSON(&ci); err != nil {
		c.JSON(http.StatusBadRequest,
			response.Error(err.Error()))
		return
	}

	sErr := h.invoiceService.Update(c.Request.Context(), middleware.UserID(c), &ci)
	if sErr != nil {
		c.JSON(http.StatusBadRequest, response.Error(sErr.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("updated."))
}
func (h *Handler) Delete(c *gin.Context) {

	id, uErr := uuid.Parse(c.Param("id"))
	if uErr != nil {
		c.JSON(http.StatusBadRequest, response.Error(uErr.Error()))
		return
	}

	if err := h.invoiceService.Delete(c.Request.Context(), id, middleware.UserID(c)); err != nil {
		c.JSON(http.StatusNotFound, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("deleted"))
}

func (h *Handler) Pay(c *gin.Context) {

	id, uErr := uuid.Parse(c.Param("id"))
	if uErr != nil {
		c.JSON(http.StatusBadRequest, response.Error(uErr.Error()))
		return
	}

	res, err := h.invoiceService.GetForPay(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusFound, res)
}
