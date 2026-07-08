package user

import (
	"context"
	"crypto-payment-gateway/pkg/password"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

var (
	ErrEmailAlreadyExists = fmt.Errorf("email exists")
	ErrUserNotExists      = fmt.Errorf("user not exists")

	ErrWrongPassword   = fmt.Errorf("wrong password")
	ErrPasswordHashing = fmt.Errorf("problem with password encryption")
)

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (ser Service) Signup(ctx context.Context, req *SignupRequest) error {

	ur, err := ser.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if ur != nil {
		return ErrEmailAlreadyExists
	}

	hash, er := password.Hash(req.Password)
	if er != nil {
		return ErrPasswordHashing
	}

	user := &User{
		Email:           req.Email,
		PasswordHash:    hash,
		WithdrawAddress: req.Wallet,
	}

	user.BeforeCreate()

	return ser.repo.Create(ctx, user)
}

func (ser Service) Login(ctx context.Context, req *LoginRequest) (*User, error) {

	ur, err := ser.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if ur == nil {
		return nil, ErrUserNotExists
	}

	if verify := password.Verify(ur.PasswordHash, req.Password); !verify {
		return nil, ErrWrongPassword
	}

	return ur, nil
}

func (ser Service) Update(ctx context.Context, ID uuid.UUID, req *UpdateRequest) error {

	return ser.repo.Update(ctx, &User{
		ID:              ID,
		WithdrawAddress: req.Wallet,
	})

}

func (ser Service) GetByID(ctx context.Context, ID uuid.UUID) (*Response, error) {

	ur, err := ser.repo.GetByID(ctx, ID.String())
	if err != nil {
		return nil, err
	}
	if ur == nil {
		return nil, ErrUserNotExists
	}

	return &Response{
		ID:      ur.ID.String(),
		Email:   ur.Email,
		Wallet:  ur.WithdrawAddress,
		Balance: ur.Balance.String(),
	}, nil
}
