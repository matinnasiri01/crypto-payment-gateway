package invoice

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateInvoiceRequest struct {
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	Description string          `json:"description"`
	CallbackURL string          `json:"callback_url"`
	Lifetime    int64           `json:"lifetime"`
}

type UpdateInvoiceRequest struct {
	ID          uuid.UUID       `json:"id"`
	Status      Status          `json:"status"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}

type InvoiceResponse struct {
	ID             uuid.UUID       `json:"id"`
	Status         Status          `json:"status"`
	Amount         decimal.Decimal `json:"amount"`
	Description    string          `json:"description"`
	PaymentAddress string          `json:"pay_to_address"`
	ExpiredAt      time.Time       `json:"expired_at"`
}
