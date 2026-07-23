package blockchain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Address struct {
	Base58 string
	index  uint32
}

type Transaction struct {
	Hash      string
	Sender    string
	Receiver  string
	Amount    decimal.Decimal
	Timestamp time.Time
}
