package models

import "github.com/google/uuid"

// CreatePaymentRequest godoc
// @Description Request body for creating a payment
type CreatePaymentRequest struct {
	UserID   string  `json:"user_id" example:"94a20574-f6f5-4659-a57b-1c8eb23591ac"`
	Amount   float64 `json:"amount" example:"100.50"`
	Currency string  `json:"currency" example:"USD"`
}

// PaymentResponse godoc
// @Description Payment object returned by create/get
type PaymentResponse struct {
	ID        uuid.UUID     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    uuid.UUID     `json:"user_id" example:"94a20574-f6f5-4659-a57b-1c8eb23591ac"`
	Amount    float64       `json:"amount" example:"100.50"`
	Currency  string        `json:"currency" example:"USD"`
	Status    PaymentStatus `json:"status" example:"created"`
	CreatedAt string        `json:"created_at" example:"2025-11-15T20:13:25Z"`
	UpdatedAt string        `json:"updated_at" example:"2025-11-15T20:13:25Z"`
}

// GetPaymentsByUserResponse godoc
// @Description Response for payments by user. Returns status of last payment
type GetPaymentsByUserResponse struct {
	Status string `json:"status" example:"success"`
}

// UpdatePaymentStatusRequest godoc
// @Description Request body for updating payment status
type UpdatePaymentStatusRequest struct {
	PaymentID string        `json:"payment_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status    PaymentStatus `json:"status" example:"success"`
}

// ErrorResponse godoc
// @Description Standard error response
type ErrorResponse struct {
	Message string `json:"message" example:"payment not found"`
}
