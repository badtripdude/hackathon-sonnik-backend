package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SubscriptionClient interface {
	GetStatus(ctx context.Context, userID string) (string, error)
}

type HttpSubscriptionClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHttpSubscriptionClient(baseURL string) *HttpSubscriptionClient {
	return &HttpSubscriptionClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Payment struct {
	Status string `json:"status"`
}

func (c *HttpSubscriptionClient) GetStatus(ctx context.Context, userID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/payments/user/%s", c.baseURL, userID), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call subscription service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var payment Payment
	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if payment.Status == "" {
		return "none", nil
	}

	return payment.Status, nil
}
