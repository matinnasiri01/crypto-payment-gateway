package http

import (
	"crypto-payment-gateway/internal/user/model"
	"crypto-payment-gateway/internal/user/repository"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Wallet   string `json:"wallet" binding:"required"`
}

func RegisterHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		var registerRequest RegisterRequest
		if err := c.BindJSON(&registerRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "jsonReq" + err.Error()})
			return
		}

		if !IsValidEthereumAddress(registerRequest.Wallet) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Ethereum wallet address"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password " + err.Error()})
			return
		}

		user := &models.User{
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
			Wallet:   registerRequest.Wallet,
		}

		err = repository.CreateUser(pool, user)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"Status": "Created!", "message": "Login to get token!"})

	}
}
