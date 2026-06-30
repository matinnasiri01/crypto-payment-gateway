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
	ErrNotFind  = fmt.Errorf("can`t find this invoice")
	ErrNoRecord = fmt.Errorf("there is no record")
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

func (s *Service) List(ctx context.Context, ID uuid.UUID, page, limit int) (*[]Invoice, error) {

	list, err := s.repo.ListByUser(ctx, ID, Pagination{Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Service) GetByID(ctx context.Context, invoiceID, userID uuid.UUID) (*Invoice, error) {

	res, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if res.UserID != userID {
		return nil, ErrNotFind
	}
	return res, nil
}

func (s *Service) GetForPay(ctx context.Context, invoiceID uuid.UUID) (*InvoiceResponse, error) {

	res, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	ri := InvoiceResponse{
		ID:             res.ID,
		Status:         res.Status,
		Amount:         res.Amount,
		Description:    res.Description,
		PaymentAddress: res.PayToAddress,
		ExpiredAt:      res.ExpiredAt,
	}

	return &ri, nil
}

func (s *Service) Delete(ctx context.Context, invoiceID, userID uuid.UUID) error {
	return s.repo.Delete(ctx, invoiceID, userID)
}

func (s *Service) Update(ctx context.Context, ID uuid.UUID, req *UpdateInvoiceRequest) error {
	return nil
}
