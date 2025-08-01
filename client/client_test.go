package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/srpvpn/tensor-go-sdk/api/user"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		config          *Config
		expectedBaseURL string
		expectedTimeout time.Duration
	}{
		{
			name:            "nil config uses defaults",
			config:          nil,
			expectedBaseURL: "https://api.mainnet.tensordev.io",
			expectedTimeout: 30 * time.Second,
		},
		{
			name:            "empty config uses defaults",
			config:          &Config{},
			expectedBaseURL: "https://api.mainnet.tensordev.io",
			expectedTimeout: 30 * time.Second,
		},
		{
			name: "custom config is used",
			config: &Config{
				BaseURL: "https://custom.api.com",
				APIKey:  "test-key",
				Timeout: 10 * time.Second,
			},
			expectedBaseURL: "https://custom.api.com",
			expectedTimeout: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.config)

			// Verify client is not nil
			if client == nil {
				t.Fatal("expected client to be non-nil")
			}

			// Verify User API is initialized
			if client.User == nil {
				t.Fatal("expected User API to be initialized")
			}

			// Verify transport is initialized
			if client.transport == nil {
				t.Fatal("expected transport to be initialized")
			}

			// Verify transport configuration (we need to cast to concrete type)
			httpTransport, ok := client.transport.(*HTTPTransport)
			if !ok {
				t.Fatal("expected HTTPTransport")
			}

			if httpTransport.baseURL != tt.expectedBaseURL {
				t.Errorf("expected baseURL %s, got %s", tt.expectedBaseURL, httpTransport.baseURL)
			}

			if httpTransport.client.Timeout != tt.expectedTimeout {
				t.Errorf("expected timeout %v, got %v", tt.expectedTimeout, httpTransport.client.Timeout)
			}
		})
	}
}
func TestClient_Close(t *testing.T) {
	client := New(nil)
	err := client.Close()
	if err != nil {
		t.Errorf("expected Close() to return nil, got %v", err)
	}
}

func TestClient_IntegrationFlow(t *testing.T) {
	// Create a test server that mimics the Tensor API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/user/portfolio" {
			t.Errorf("expected path /api/v1/user/portfolio, got %s", r.URL.Path)
		}

		// Check query parameters
		wallet := r.URL.Query().Get("wallet")
		if wallet != "11111111111111111111111111111111" {
			t.Errorf("expected wallet parameter '11111111111111111111111111111111', got '%s'", wallet)
		}

		// Check headers
		userAgent := r.Header.Get("User-Agent")
		if userAgent != "tensor-go-sdk/1.0.0" {
			t.Errorf("expected User-Agent 'tensor-go-sdk/1.0.0', got '%s'", userAgent)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "success",
			"collections": [
				{
					"id": "test-collection-1",
					"name": "Test Collection",
					"symbol": "TEST",
					"image": "https://example.com/image.png",
					"floorPrice": 1.5,
					"volume24h": 100.0,
					"verified": true,
					"compressed": false
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := New(config)

	// Test the full flow
	ctx := context.Background()
	req := &user.PortfolioRequest{
		Wallet: "11111111111111111111111111111111", // Valid 32-character wallet address
	}

	body, statusCode, err := client.User.GetPortfolio(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify status code
	if statusCode != 200 {
		t.Errorf("expected status code 200, got %d", statusCode)
	}

	// Verify response body is not empty
	if len(body) == 0 {
		t.Fatal("expected response body to be non-empty")
	}

	// For simplicity, just check that the response contains key elements
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "test-collection-1") {
		t.Errorf("expected response to contain 'test-collection-1', got: %s", bodyStr)
	}
}
func TestClient_IntegrationFlow_WithAPIKey(t *testing.T) {
	// Create a test server that checks for API key
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check API key header
		apiKey := r.Header.Get("x-tensor-api-key")
		if apiKey != "test-api-key" {
			t.Errorf("expected API key 'test-api-key', got '%s'", apiKey)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "success",
			"collections": []
		}`))
	}))
	defer server.Close()

	// Create client with API key
	config := &Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
		Timeout: 5 * time.Second,
	}
	client := New(config)

	// Test request with API key
	ctx := context.Background()
	req := &user.PortfolioRequest{
		Wallet: "11111111111111111111111111111111", // Valid 32-character wallet address
	}

	_, _, err := client.User.GetPortfolio(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestClient_IntegrationFlow_ErrorHandling(t *testing.T) {
	// Create a test server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{
			"code": 422,
			"message": "Validation error",
			"details": "Invalid wallet address"
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := New(config)

	// Test error handling
	ctx := context.Background()
	req := &user.PortfolioRequest{
		Wallet: "invalid-wallet",
	}

	_, _, err := client.User.GetPortfolio(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// The error should be wrapped and contain information about the API error
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestClient_IntegrationFlow_ContextCancellation(t *testing.T) {
	// Create a test server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success", "collections": []}`))
	}))
	defer server.Close()

	// Create client with test server URL
	config := &Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := New(config)

	// Test context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req := &user.PortfolioRequest{
		Wallet: "11111111111111111111111111111111", // Valid 32-character wallet address
	}

	_, _, err := client.User.GetPortfolio(ctx, req)
	if err == nil {
		t.Fatal("expected context timeout error, got nil")
	}

	// Should be a context deadline exceeded error
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected context deadline exceeded, got %v", ctx.Err())
	}
}
