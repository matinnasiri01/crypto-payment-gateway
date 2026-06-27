package user

import (
	"context"
	"crypto-payment-gateway/pkg/password"
	"fmt"
)

type Service struct {
	repo Repository
}

var (
	ErrEmailAlreadyExists = fmt.Errorf("email exists")
	ErrEmailNotExists     = fmt.Errorf("email not exists")

	ErrWrongPassword   = fmt.Errorf("wrong password")
	ErrPasswordHashing = fmt.Errorf("problem with password encryption")
)

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (ser Service) Signup(req *SignupRequest) error {
	ctx := context.Background()

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

func (ser Service) Login(req *LoginRequest) (*User, error) {

	ctx := context.Background()

	ur, err := ser.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if ur == nil {
		return nil, ErrEmailNotExists
	}

	if verify := password.Verify(ur.PasswordHash, req.Password); !verify {
		return nil, ErrWrongPassword
	}

	return ur, nil
}

func (ser Service) Update(req *LoginRequest) error {
	return nil
}

func (ser Service) GetByID(ID string) (*MeResponse, error) {
	return nil, nil
}
