package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(p *pgxpool.Pool) *userPostgresRepo {
	return &userPostgresRepo{pool: p}
}

func (r *userPostgresRepo) Create(ctx context.Context, user *User) error {

	query := `
	INSERT INTO users (
		id,
		email,
		password_hash,
		withdraw_address,
		balance,
		created_at,
		updated_at
	)
	VALUES (
		@id,
		@email,
		@password_hash,
		@withdraw_address,
		@balance,
		@created_at,
		@updated_at
	)
	`

	args := pgx.NamedArgs{
		"id":               user.ID,
		"email":            user.Email,
		"password_hash":    user.PasswordHash,
		"withdraw_address": user.WithdrawAddress,
		"balance":          user.Balance,
		"created_at":       user.CreatedAt,
		"updated_at":       user.UpdatedAt,
	}

	_, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *userPostgresRepo) Update(cnx context.Context, user *User) error {
	return nil
}
func (r *userPostgresRepo) GetByID(cnx context.Context, id string) (*User, error) {
	var user User

	err := r.pool.QueryRow(cnx, `
        SELECT
            id,
            email,
            password_hash,
            withdraw_address,
            balance,
            created_at,
            updated_at
        FROM users
        WHERE id = $1
    `, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.WithdrawAddress,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userPostgresRepo) GetByEmail(cnx context.Context, email string) (*User, error) {
	var user User

	err := r.pool.QueryRow(cnx, `
        SELECT
            id,
            email,
            password_hash,
            withdraw_address,
            balance,
            created_at,
            updated_at
        FROM users
        WHERE email = $1
    `, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.WithdrawAddress,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
