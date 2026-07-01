package client

import (
	"context"
	"net/http"
	"time"
)

// NewTestClient returns a Client configured for unit tests against a custom HTTP server.
func NewTestClient(httpClient *http.Client, ruleBasePath string) *Client {
	return &Client{
		userAgent:          "test",
		requestAttempts:    1,
		requestTimeout:     5 * time.Second,
		context:            context.Background(),
		rateLimiters:       *NewClientRateLimiters(),
		backstoryAPIClient: httpClient,
		RuleBasePath:       ruleBasePath,
	}
}
