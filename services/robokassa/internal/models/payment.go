package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusCreated PaymentStatus = "created"
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	UserID    uuid.UUID     `db:"user_id" json:"user_id"`
	Amount    float64       `db:"amount" json:"amount"`
	Currency  string        `db:"currency" json:"currency"`
	Status    PaymentStatus `db:"status" json:"status"`
	CreatedAt time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" json:"updated_at"`
}
