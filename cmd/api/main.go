package main

import (
	"context"
	"crypto-payment-gateway/internal/invoice"
	"crypto-payment-gateway/internal/middleware"
	"crypto-payment-gateway/internal/user"
	"crypto-payment-gateway/pkg/database"
	"crypto-payment-gateway/pkg/jwt"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("can`t find .env")
	}
}

func main() {

	pool, err := database.NewPostgresDB(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can`t connect to db")
	}

	jwtManager := jwt.New(os.Getenv("JWT_SECRET"))
	auth := middleware.NewAuth(jwtManager)

	// user
	ur := user.NewPostgresRepo(pool.Pool)
	us := user.NewService(ur)
	uh := user.NewHandler(us, jwtManager)

	// invoice
	ir := invoice.NewPostgresRepo(pool.Pool)
	is := invoice.NewService(ir)
	ih := invoice.NewHandler(is, jwtManager)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Oky!")

	})

	api := r.Group("/api/v1")

	uh.Register(api, auth)
	ih.Register(api, auth)

	_ = r.Run(":" + os.Getenv("PORT"))

}
