// client/transport.go
package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/srpvpn/tensor-go-sdk/internal/errors"
	"github.com/srpvpn/tensor-go-sdk/internal/transport"
)

// HTTPTransport implements the transport.Transport interface using HTTP.
type HTTPTransport struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// NewTransport creates a new HTTPTransport with the given configuration.
func NewTransport(cfg Config) transport.Transport {
	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	return &HTTPTransport{
		client:  client,
		baseURL: strings.TrimSuffix(cfg.BaseURL, "/"),
		apiKey:  cfg.APIKey,
	}
}

// Get performs a GET request with context support and query parameters.
func (t *HTTPTransport) Get(ctx context.Context, path string, params url.Values) (*http.Response, error) {
	// Build the full URL
	fullURL := t.baseURL + path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, &errors.NetworkError{
			Op:  "create_request",
			Err: fmt.Errorf("failed to create HTTP request: %w", err),
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "tensor-go-sdk/1.0.0")

	if t.apiKey != "" {
		req.Header.Set("x-tensor-api-key", t.apiKey)
	}

	// Perform the request
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, &errors.NetworkError{
			Op:  "http_request",
			Err: fmt.Errorf("HTTP request failed: %w", err),
		}
	}

	// Check for HTTP errors and parse API errors
	if resp.StatusCode >= 400 {
		apiErr := errors.ParseAPIError(resp)
		resp.Body.Close() // Close the body since we're returning an error
		return nil, apiErr
	}

	return resp, nil
}
