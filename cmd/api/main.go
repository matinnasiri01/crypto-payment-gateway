package main

import (
	"context"
	"crypto-payment-gateway/internal/user"
	"crypto-payment-gateway/pkg/database"
	"fmt"
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
		fmt.Print("can`t connect to db")
	}

	ur := user.NewPostgresRepo(pool.Pool)
	us := user.NewService(ur)
	uh := user.NewHandler(us)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Oky!")

	})

	api := r.Group("/api/v1")
	uh.Register(api)

	_ = r.Run(":" + os.Getenv("PORT"))

}
