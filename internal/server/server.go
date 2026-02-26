package server

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"crypto-payment-gateway/internal/config"


)
type Server struct {
	eng  *gin.Engine
	cfg  *config.Config
	db   *pgxpool.Pool
}

func NewServer (db   *pgxpool.Pool,cfg  *config.Config) *Server{
	return &Server{ 
		eng: gin.Default(),
		cfg: cfg,
		db: db,
	}
}

func (s Server) Run() error {
	_ = s.eng.SetTrustedProxies(nil)

	s.eng.GET("/helth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "CPG API is running well!",
			"status":   "success",
			"database": "connected",
		})
	})
	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}
	s.eng.Run(":" + s.cfg.Port)

	return nil
} 


func (s Server) MapRoutes() error {
	_ = s.eng.Group("/api/v1")
	
	return nil
}