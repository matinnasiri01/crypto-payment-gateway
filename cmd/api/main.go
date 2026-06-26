package main

import (
	"context"
	"crypto-payment-gateway/pkg/database"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("can`t find .env")
	}
}

func main() {

	_, err := database.NewPostgresDB(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Print("can`t connect to db")
	}
}
