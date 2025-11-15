package models

import "time"

type LoginResponse struct {
    UserID             string  `json:"user_id"`
    AccessToken        string  `json:"access_token"`
    RefreshToken       string  `json:"refresh_token"`
    ExpiresAt          time.Time `json:"expires_at"`
    SubscriptionStatus *string `json:"subscription_status,omitempty"`
}
