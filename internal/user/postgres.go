package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(p *pgxpool.Pool) *userPostgresRepo {
	return &userPostgresRepo{pool: p}
}

func (r *userPostgresRepo) Create(cnx context.Context, user *User) error {
	return nil
}
func (r *userPostgresRepo) Update(cnx context.Context, user *User) error {
	return nil
}
func (r *userPostgresRepo) GetUserByID(cnx context.Context, id string) (*User, error) {
	return nil, nil
}
