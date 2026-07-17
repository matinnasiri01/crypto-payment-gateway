package invoice

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPaid      Status = "paid"
	StatusCancelled Status = "cancelled"
	StatusExpired   Status = "expired"
	StatusDeleted   Status = "deleted"
)

type Invoice struct {
	ID     uuid.UUID
	UserID uuid.UUID

	HDIndex     uint32
	Status      Status
	Amount      decimal.Decimal
	Overpayment decimal.Decimal

	Description string
	CallbackURL string

	PayToAddress  string
	PaidByAddress string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
}

func (i *Invoice) BeforeCreate() {
	now := time.Now().UTC()

	i.ID = uuid.New()
	i.Status = StatusPending
	i.CreatedAt = now
	i.UpdatedAt = now
}

func (i *Invoice) IsExpired() bool {
	return time.Now().UTC().After(i.ExpiredAt)
}
