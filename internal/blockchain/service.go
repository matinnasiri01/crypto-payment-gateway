package blockchain

import (
	"context"

	"github.com/shopspring/decimal"
)

type Blockchain interface {
	ValidateAddress(address string) error

	GenerateDepositAddress(
		ctx context.Context,
	) (string, error)

	Transfer(
		ctx context.Context,
		to string,
		amount decimal.Decimal,
	) (string, error)

	Transactions(
		ctx context.Context,
		address string,
	) ([]Transaction, error)
}
