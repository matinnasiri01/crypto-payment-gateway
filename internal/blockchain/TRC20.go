package blockchain

import (
	"context"
	"os"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	address2 "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/contract"
	"github.com/fbsobreira/gotron-sdk/pkg/keys"
	"github.com/fbsobreira/gotron-sdk/pkg/signer"
	"github.com/fbsobreira/gotron-sdk/pkg/standards/trc20"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TRC20 struct {
	network *Network
	client  *client.GrpcClient
	token   *trc20.Token
}

var Nile = Network{
	Endpoint: "grpc.nile.trongrid.io:50051",
	Contract: "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf",
	FeeLimit: 100_000_000,
}

func NewTRC20(network *Network) *TRC20 {

	conn := client.NewGrpcClient(network.Endpoint)
	if err := conn.Start(grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		logger.Fatal(os.Stderr, "connect error: %v\n", err)
	}

	token := trc20.New(conn, network.Contract)

	return &TRC20{network: network, client: conn, token: token}
}

func (t *TRC20) IsValidateAddress(address string) bool {

	toAddress, err := address2.Base58ToAddress(address)
	if err != nil {
		return false
	}

	return toAddress.IsValid()
}

func (t *TRC20) Balance(ctx context.Context, address string) (decimal.Decimal, error) {

	bl, err := t.token.BalanceOf(ctx, address)
	if err != nil {
		return decimal.Zero, err
	}

	dec, err := decimal.NewFromString(bl.Display)
	if err != nil {
		return decimal.Zero, err
	}

	return dec, nil
}

func (t *TRC20) Transfer(ctx context.Context, from, to, pk string, amount decimal.Decimal) (string, error) {

	trc20Tx := t.token.Transfer(from, to, amount.BigInt(),
		contract.WithFeeLimit(t.network.FeeLimit))

	s, err := t.newSigner(pk)
	if err != nil {
		return "", err
	}

	res, err := trc20Tx.Send(ctx, s)
	if err != nil {
		return "", err
	}

	return res.TxID, nil
}

func (t *TRC20) newSigner(privateKey string) (signer.Signer, error) {
	key, err := keys.GetPrivateKeyFromHex(privateKey)
	if err != nil {
		return nil, err
	}

	s, err := signer.NewPrivateKeySignerFromBTCEC(key)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (t *TRC20) Transactions(ctx context.Context, address string, after time.Time) ([]Transaction, error) {

	return nil, nil
}
