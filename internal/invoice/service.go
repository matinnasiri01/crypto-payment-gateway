package invoice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

const defaultLifetime = 1800

var (
	ErrLifetime = fmt.Errorf("1800 < lifetime < 86400")
)

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(ctx context.Context, ID uuid.UUID, req *CreateInvoiceRequest) error {

	lifetime := defaultLifetime
	if req.Lifetime != 0 {
		if req.Lifetime < 1800 || req.Lifetime > 86400 {
			return ErrLifetime
		}

		lifetime = int(req.Lifetime)
	}

	invo := &Invoice{
		UserID:      ID,
		Amount:      req.Amount,
		Description: req.Description,
		CallbackURL: req.CallbackURL,
		ExpiredAt:   time.Now().Add(time.Duration(lifetime) * time.Second),
	}

	invo.BeforeCreate()
	err := s.repo.Create(ctx, invo)
	if err != nil {
		return err
	}

	return nil
}
