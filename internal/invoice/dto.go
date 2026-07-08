package invoice

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateRequest struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Description string          `json:"description"`
	CallbackURL string          `json:"callback_url"`
	Lifetime    int64           `json:"lifetime"`
}

type UpdateRequest struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Description string          `json:"description"`
}

type ListResponse struct {
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Count int        `json:"count"`
	List  *[]Invoice `json:"list"`
}

type Response struct {
	ID            uuid.UUID       `json:"id"`
	Status        Status          `json:"status"`
	Amount        decimal.Decimal `json:"amount"`
	Description   string          `json:"description"`
	PayToAddress  string          `json:"pay_to_address"`
	PaidByAddress string          `json:"paid_by_address,omitempty"`
	Overpayment   decimal.Decimal `json:"overpayment,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	ExpiredAt     time.Time       `json:"expired_at"`
}
