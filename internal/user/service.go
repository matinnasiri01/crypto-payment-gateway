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
	ErrPasswordHashing    = fmt.Errorf("problem with password encryption")
)

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (ser Service) SignUp(req *SignupRequest) error {
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

func (ser Service) LogIn(req *LoginRequest) error {
	return nil
}

func (ser Service) Update(req *LoginRequest) error {
	return nil
}

func (ser Service) GetUserByID(ID string) (*MeResponse, error) {
	return nil, nil
}
