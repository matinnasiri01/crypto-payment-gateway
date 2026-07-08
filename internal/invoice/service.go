package invoice

import (
	"context"
	"crypto-payment-gateway/internal/blockchain"

	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service struct {
	repo  Repository
	chain blockchain.Blockchain
}

const defaultLifetime = 1800

var (
	ErrLifetime      = fmt.Errorf("1800 < lifetime < 86400")
	ErrNotFind       = fmt.Errorf("can`t find this invoice")
	ErrInvalidAmount = fmt.Errorf("amount error")
)

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(ctx context.Context, ID uuid.UUID, req *CreateRequest) error {

	lifetime := defaultLifetime
	if req.Lifetime != 0 {
		if req.Lifetime < 1800 || req.Lifetime > 86400 {
			return ErrLifetime
		}

		lifetime = int(req.Lifetime)
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
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

func (s *Service) List(ctx context.Context, ID uuid.UUID, page, limit int) (*ListResponse, error) {

	list, err := s.repo.ListByUser(ctx, ID, Pagination{Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}

	return &ListResponse{page, limit, len(*list), list}, nil
}

func (s *Service) GetByID(ctx context.Context, invoiceID, userID uuid.UUID) (*Response, error) {

	res, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if res.UserID != userID {
		return nil, ErrNotFind
	}
	return &Response{
		ID:            res.ID,
		Status:        res.Status,
		Amount:        res.Amount,
		Description:   res.Description,
		PayToAddress:  res.PayToAddress,
		PaidByAddress: res.PaidByAddress,
		Overpayment:   res.Overpayment,
		CreatedAt:     res.CreatedAt,
		ExpiredAt:     res.ExpiredAt,
	}, nil
}

func (s *Service) GetForPay(ctx context.Context, invoiceID uuid.UUID) (*Response, error) {

	res, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	ri := Response{
		ID:          res.ID,
		Status:      res.Status,
		Amount:      res.Amount,
		Description: res.Description,
		ExpiredAt:   res.ExpiredAt,
	}

	return &ri, nil
}

func (s *Service) Delete(ctx context.Context, invoiceID, userID uuid.UUID) error {
	return s.repo.Delete(ctx, invoiceID, userID)
}

func (s *Service) Update(ctx context.Context, ID uuid.UUID, req *UpdateRequest) error {
	return s.repo.Update(ctx, &Invoice{
		UserID:      ID,
		Amount:      req.Amount,
		Description: req.Description,
	})
}

func (s *Service) StartWatcher(ctx context.Context) {

}

func (s *Service) StartWorker(ctx context.Context) {

}
