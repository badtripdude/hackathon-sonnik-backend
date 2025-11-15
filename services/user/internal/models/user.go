package models

import "time"

type User struct {
	ID                 string     `db:"id" json:"id" example:"37d86688-021e-41ce-932a-841a047f0454"`
	Email              string     `db:"email" json:"email" example:"test@example.com"`
	PasswordHash       string     `db:"password_hash" json:"password_hash" example:"$2a$12$1eBqZjVvAY.gmfcRUlR24uR8sSZHcF/DIeJu/eM7z0USl6OVCFVFy"`
	Username           *string    `db:"username" json:"username" example:"JohnDoe"`
	BirthDate          *time.Time `db:"birth_date" json:"birth_date" example:"1990-01-01"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at" example:"2025-11-15T09:48:37Z"`

	// for API
	SubscriptionStatus *string `json:"subscription_status,omitempty"`
}
