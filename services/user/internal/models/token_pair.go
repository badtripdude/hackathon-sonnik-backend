package models

import "time"

type TokenPair struct {
	AccessToken  string    `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refresh_token" example:"auGdrFA8HKbLrEY-2X9KTAtvzh9Qh-obHElRMtEg..."`
	ExpiresAt    time.Time `json:"expires_at" example:"2025-11-15T10:18:37Z"`
	UserID       string    `json:"user_id" example:"37d86688-021e-41ce-932a-841a047f0454"`
}
