package blockchain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Hash      string
	Sender    string
	Receiver  string
	Amount    decimal.Decimal
	Timestamp time.Time
}
