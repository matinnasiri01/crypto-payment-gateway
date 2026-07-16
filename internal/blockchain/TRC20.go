package blockchain

import (
	"context"
	"crypto-payment-gateway/pkg/tron"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/ryanbekhen/tronwallet"
	"github.com/shopspring/decimal"
)

type TRC20 struct {
	client *tron.Client
	wallet *tronwallet.TronWallet
}

func NewTRC20(cli *tron.Client, mnemonic string) *TRC20 {
	w, er := tronwallet.RestoreWallet(mnemonic)
	if er != nil {
		logger.Fatal(er.Error())
	}

	return &TRC20{client: cli, wallet: w}
}

func (t *TRC20) ValidateAddress(address string) error {

	add, err := tron.Base58ToAddress(address)
	if err != nil {
		return err
	}

	if add.IsValid() {
		return nil
	}
	return err

}

func (t *TRC20) GenerateDepositAddress(index uint32) (string, error) {

	derive, err := t.wallet.Derive(index)
	if err != nil {
		return "", err
	}

	return tronwallet.TronAddressFromPrivate(derive), nil
}

func (t *TRC20) Balance(ctx context.Context, address string) (decimal.Decimal, error) {
	return t.client.GetBalance(ctx, []byte(address))
}

func (t *TRC20) Transfer(ctx context.Context, to string, amount decimal.Decimal) (string, error) {
	return "", nil
}

func (t *TRC20) Transactions(ctx context.Context, address string, from time.Time) ([]Transaction, error) {
	return nil, nil
}

func (t *TRC20) TransactionByHash(ctx context.Context, hash string) (*Transaction, error) {
	return nil, nil
}
