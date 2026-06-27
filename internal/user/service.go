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
		return nil, ErrUserNotExists
	}

	if verify := password.Verify(ur.PasswordHash, req.Password); !verify {
		return nil, ErrWrongPassword
	}

	return ur, nil
}

func (ser Service) Update(req *LoginRequest) error {
	return nil
}

func (ser Service) GetByID(ID uuid.UUID) (*MeResponse, error) {

	ctx := context.Background()

	ur, err := ser.repo.GetByID(ctx, ID.String())
	if err != nil {
		return nil, err
	}
	if ur == nil {
		return nil, ErrUserNotExists
	}

	return ur.Convert(), nil
}
