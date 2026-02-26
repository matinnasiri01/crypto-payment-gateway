package repository

import (
	"context"
	"time"	
	"crypto-payment-gateway/internal/user/model"


	"github.com/jackc/pgx/v5/pgxpool"
)
func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.BeforeCreate()

	query := `
		INSERT INTO users (id, email, password, wallet) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, wallet, currency, balance, created_at
	`

	err := pool.QueryRow(ctx, query, 
		user.ID,        
		user.Email, 
		user.Password, 
		user.Wallet,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Wallet,
		&user.Currency,
		&user.Balance,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}