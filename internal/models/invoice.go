package models

import (
	"time"

	_ "github.com/google/uuid"
)

type Status string

const (
	Completed Status = "Completed"
	Pending   Status = "Pending"
	Failed    Status = "Failed"
)

type Invoice struct {
	ID        string    `json:"id"`
	CreatorID string    `json:"creatorid"`
	Address   string    `json:"address"`
	Amount    float32   `json:"amount"`
	//Currency  Cu        `json:"currency"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
