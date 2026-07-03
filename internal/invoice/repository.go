package invoice

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, invoice *Invoice) error
	Update(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, invoiceID uuid.UUID) (*Invoice, error)
	ListByUser(ctx context.Context, userID uuid.UUID, p Pagination) (*[]Invoice, error)
	Delete(ctx context.Context, invoiceID, userID uuid.UUID) error
}

type Pagination struct {
	Page  int
	Limit int
}
