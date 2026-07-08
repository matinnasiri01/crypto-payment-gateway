package blockchain

import (
	"context"

	"github.com/shopspring/decimal"
)

type TRC20 struct{}

func NewTRC20() *TRC20 {
	return &TRC20{}
}

func (t *TRC20) ValidateAddress(address string) error {
	return nil
}
func (t *TRC20) GenerateDepositAddress(ctx context.Context) (string, error) {
	return "", nil
}
func (t *TRC20) Transfer(ctx context.Context, to string, amount decimal.Decimal) (string, error) {
	return "", nil
}
func (t *TRC20) Transactions(ctx context.Context, address string) ([]Transaction, error) {
	return nil, nil
}
