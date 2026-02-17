package main

import (
	"log"



	"crypto-payment-gateway/internal/config"
	"crypto-payment-gateway/internal/database"
	
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)


func main() {
	var cfg *config.Config 
	var err error


	cfg = config.Load()

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "CPG API is running well!",
			"status":   "success",
			"database": "connected",
		})
	})

	router.Run(":" + cfg.Port)
	
}
