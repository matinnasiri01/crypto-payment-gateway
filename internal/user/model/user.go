package models

import (
	"time"

	"github.com/google/uuid"
	_"github.com/jackc/pgx/v5/pgxpool"
)

type Cu string

const (
	ETH Cu = "eth"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Wallet    string    `json:"wallet"`
	Balance   float32   `json:"balance"`
	Currency  Cu        `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func (user *User) BeforeCreate() error {
	user.ID = uuid.New().String()
	if user.Currency == "" {
		user.Currency = ETH 
	}
	return nil
}