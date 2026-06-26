package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID              uuid.UUID
	Email           string
	PasswordHash    string
	WithdrawAddress string
	Balance         decimal.Decimal
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
