package blockchain

import (
	"context"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/ryanbekhen/tronwallet"
	"github.com/shopspring/decimal"
)

type TRC20 struct {
	wallet *tronwallet.TronWallet
}

func NewTRC20(mnemonic string) *TRC20 {
	w, er := tronwallet.RestoreWallet(mnemonic)
	if er != nil {
		logger.Fatal(er.Error())
	}

	return &TRC20{wallet: w}
}

func (t *TRC20) ValidateAddress(address string) error {

	return nil

}

func (t *TRC20) GenerateDepositAddress(index uint32) (string, error) {
	derive, err := t.wallet.Derive(index)
	if err != nil {
		return "", err
	}

	return tronwallet.TronAddressFromPrivate(derive), nil
}

func (t *TRC20) Balance(ctx context.Context, address string) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

func (t *TRC20) Transfer(ctx context.Context, to string, amount decimal.Decimal) (string, error) {
	return "", nil
}

func (t *TRC20) Transactions(ctx context.Context, address string, after time.Time) ([]Transaction, error) {
	return nil, nil
}

func (t *TRC20) TransactionByHash(ctx context.Context, hash string) (*Transaction, error) {
	return nil, nil
}
