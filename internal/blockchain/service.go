package blockchain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Blockchain interface {
	ValidateAddress(address string) error

	GenerateDepositAddress(index uint32) (string, error)

	Balance(ctx context.Context, address string) (decimal.Decimal, error)

	Transfer(ctx context.Context, to string, amount decimal.Decimal) (string, error)

	Transactions(ctx context.Context, address string, from time.Time) ([]Transaction, error)

	TransactionByHash(ctx context.Context, hash string) (*Transaction, error)
}
