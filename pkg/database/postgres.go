package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	*pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, url string) (*PostgresDB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &PostgresDB{pool}, nil
}
