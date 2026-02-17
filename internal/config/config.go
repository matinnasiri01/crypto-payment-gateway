package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var doOnce sync.Once
var singleton *Config


type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
}



func Load() *Config {
	doOnce.Do(func() {
		var err error = godotenv.Load()
		
		if err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
		
		singleton =  &Config{
			DatabaseURL: os.Getenv("DATABASE_URL"),
			Port:        os.Getenv("PORT"),
			JWTSecret:   os.Getenv("JWT_SECRET"),
		}
		
	})
	return singleton
}