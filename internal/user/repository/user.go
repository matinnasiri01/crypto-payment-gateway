package repository

import (
	"time"
	"context"
	"crypto-payment-gateway/internal/user/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.BeforeCreate()

	query := `
		INSERT INTO users (id, email, password, wallet) 
		VALUES ($1, $2, $3, $4)
	`

	_, err := pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.Wallet,
	)

	if err != nil {
		return err
	}

	return nil
}
