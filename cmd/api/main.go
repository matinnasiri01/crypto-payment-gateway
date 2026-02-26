package main

import (
	"log"

	"crypto-payment-gateway/internal/config"
	"crypto-payment-gateway/internal/database"
	"crypto-payment-gateway/internal/server"

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

	


	server := server.NewServer(pool,cfg)
	server.Run()
	
}
