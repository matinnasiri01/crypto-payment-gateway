package main

// @title Crypto Payment Gateway API
// @version 1.0
// @description Crypto Payment Gateway for accepting USDT payments.

// @contact.name Matin Nasiri
// @contact.email mnasirii829@gmail.com

// @license.name MIT

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token

import (
	"context"
	"crypto-payment-gateway/internal/blockchain"
	"crypto-payment-gateway/internal/invoice"
	"crypto-payment-gateway/internal/middleware"
	"crypto-payment-gateway/internal/user"
	"crypto-payment-gateway/pkg/database"
	"crypto-payment-gateway/pkg/jwt"
	"crypto-payment-gateway/pkg/tron"
	"log"
	"net/http"
	"os"

	"crypto-payment-gateway/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("can`t find .env")
	}
}

func main() {

	pool, err := database.NewPostgresDB(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can`t connect to db")
	}

	jwtManager := jwt.New(os.Getenv("JWT_SECRET"))
	auth := middleware.NewAuth(jwtManager)

	tn := tron.NewClient(tron.Config{
		Network: tron.Nila,
		APIKey:  os.Getenv("TRON_GRID_APIKEY"),
	})

	// blockchain
	chain := blockchain.NewTRC20(tn)

	// user
	ur := user.NewPostgresRepo(pool.Pool)
	us := user.NewService(ur, chain)
	uh := user.NewHandler(us, jwtManager)

	// invoice
	ir := invoice.NewPostgresRepo(pool.Pool)
	is := invoice.NewService(ir)
	ih := invoice.NewHandler(is)

	is.StartWatcher(context.Background())
	go is.StartWorker(context.Background())

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", HealthCheck)

	api := r.Group("/api/v1")
	docs.SwaggerInfo.BasePath = "/api/v1"

	uh.Register(api, auth)
	ih.Register(api, auth)

	_ = r.Run(":" + os.Getenv("PORT"))

}

// @Summary Health Check
// @Description Check server status
// @Tags Health
// @Produce plain
// @Success 200 {string} string "OK!"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK!")
}
