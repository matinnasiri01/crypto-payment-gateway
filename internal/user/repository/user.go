package repository

import (
	"context"
	"crypto-payment-gateway/internal/user/model"
	"time"

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

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT *
		FROM users
		WHERE email = $1
	`
	var user models.User

	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Wallet,
		&user.Balance,
		&user.Currency,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
