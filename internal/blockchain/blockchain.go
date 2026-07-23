package blockchain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Blockchain interface {
	IsValidateAddress(address string) bool

	Balance(ctx context.Context, address string) (decimal.Decimal, error)

	Transfer(ctx context.Context, from, to, pk string, amount decimal.Decimal) (string, error)

	Transactions(ctx context.Context, address string, from time.Time) ([]Transaction, error)
}
