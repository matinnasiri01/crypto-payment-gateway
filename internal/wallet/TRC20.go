package wallet

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/ryanbekhen/tronwallet"
)

type TRC20Wallet struct {
	wallet *tronwallet.TronWallet
}

func New(mnemonic string) *TRC20Wallet {
	w, er := tronwallet.RestoreWallet(mnemonic)
	if er != nil {
		logger.Fatal(er.Error())
	}

	return &TRC20Wallet{wallet: w}
}

func (tw *TRC20Wallet) Address(index uint32) (string, error) {

	derive, err := tw.wallet.Derive(index)
	if err != nil {
		return "", err
	}

	return tronwallet.TronAddressFromPrivate(derive), nil
}

func (tw *TRC20Wallet) PrivateKey(index uint32) (string, error) {
	derive, err := tw.wallet.Derive(index)
	if err != nil {
		return "", err
	}

	return tronwallet.PrivateKeyToHex(derive), nil
}
