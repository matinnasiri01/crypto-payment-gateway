package invoice

import (
	"context"
	"crypto-payment-gateway/internal/blockchain"
	"crypto-payment-gateway/internal/wallet"
	"log"

	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service struct {
	repo   Repository
	chain  blockchain.Blockchain
	wallet wallet.Wallet
}

const defaultLifetime = 1800

var (
	ErrLifetime      = fmt.Errorf("1800 < lifetime < 86400")
	ErrNotFind       = fmt.Errorf("can`t find this invoice")
	ErrInvalidAmount = fmt.Errorf("amount error")
)

func NewService(r Repository, c blockchain.Blockchain, w wallet.Wallet) *Service {
	return &Service{
		repo:   r,
		chain:  c,
		wallet: w,
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

	// todo mutex!
	index, ier := s.repo.GetLastIndex(ctx)
	if ier != nil {
		return ier
	}
	index += 1

	add, der := s.wallet.Address(index)
	if der != nil {
		return der
	}

	invo := &Invoice{
		UserID:       ID,
		HDIndex:      index,
		PayToAddress: add,
		Amount:       req.Amount,
		Description:  req.Description,
		CallbackURL:  req.CallbackURL,
		ExpiredAt:    time.Now().Add(time.Duration(lifetime) * time.Second),
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
	log.Println("Invoice Watcher is running!")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		list, err := s.repo.GetPending(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, item := range *list {
			s.scanTransactions(
				ctx,
				&item,
				func(t *blockchain.Transaction) {

					item.Status = StatusPaid
					item.PaidByAddress = t.Hash

					err := s.repo.Update(ctx, &item)
					if err != nil {
						return
					}
				},
				func(t *[]blockchain.Transaction, over decimal.Decimal) {
					item.Status = StatusPaid
					item.Overpayment = over

					err := s.repo.Update(ctx, &item)
					if err != nil {
						return
					}
				},
				func(t *[]blockchain.Transaction, under decimal.Decimal) {
					// expire check system
				})
		}
	}
}

func (s *Service) scanTransactions(
	ctx context.Context,
	invoice *Invoice,
	isEqual func(t *blockchain.Transaction),
	isOver func(t *[]blockchain.Transaction, over decimal.Decimal),
	isUnder func(t *[]blockchain.Transaction, under decimal.Decimal),
) {

	transactions, err := s.chain.Transactions(ctx, invoice.PayToAddress, invoice.CreatedAt)
	if err != nil {
		log.Println(err)
		return
	}

	if len(transactions) == 0 {
		return

	}
	target := decimal.Zero
	for _, trans := range transactions {

		if trans.Amount.Equal(invoice.Amount) {
			isEqual(&trans)
			break
		}
		target.Add(trans.Amount)
	}

	if target.GreaterThanOrEqual(invoice.Amount) {
		isOver(&transactions, target)
		return
	}
	if target.LessThanOrEqual(invoice.Amount) {
		isUnder(&transactions, target)
		return
	}

}

func (s *Service) StartWorker(ctx context.Context) {
	log.Println("Invoice Worker is running!")

	for {

		err := s.expireChecker(ctx)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(10 * time.Second)
	}
}

func (s *Service) overpaymentChecker(ctx context.Context) error {
	return nil

}

func (s *Service) expireChecker(ctx context.Context) error {

	inv, er := s.repo.GetPending(ctx)
	if er != nil {
		return er
	}

	for _, item := range *inv {

		if item.IsExpired() {
			item.Status = StatusExpired

			if e := s.repo.UpdateStatus(ctx, &item); e != nil {
				return e
			}
		}
	}

	return nil
}
