// client/transport_test.go
package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	apierrors "github.com/srpvpn/tensor-go-sdk/internal/errors"
)

func TestNewTransport(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   *HTTPTransport
	}{
		{
			name: "basic config",
			config: Config{
				BaseURL: "https://api.example.com",
				APIKey:  "test-key",
				Timeout: 30 * time.Second,
			},
			want: &HTTPTransport{
				baseURL: "https://api.example.com",
				apiKey:  "test-key",
			},
		},
		{
			name: "config with trailing slash",
			config: Config{
				BaseURL: "https://api.example.com/",
				APIKey:  "test-key",
				Timeout: 30 * time.Second,
			},
			want: &HTTPTransport{
				baseURL: "https://api.example.com",
				apiKey:  "test-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTransport(tt.config)

			// Cast to concrete type to access fields for testing
			httpTransport, ok := got.(*HTTPTransport)
			if !ok {
				t.Fatal("Expected *HTTPTransport")
			}

			if httpTransport.baseURL != tt.want.baseURL {
				t.Errorf("NewTransport().baseURL = %v, want %v", httpTransport.baseURL, tt.want.baseURL)
			}
			if httpTransport.apiKey != tt.want.apiKey {
				t.Errorf("NewTransport().apiKey = %v, want %v", httpTransport.apiKey, tt.want.apiKey)
			}
			if httpTransport.client == nil {
				t.Error("NewTransport().client is nil")
			}
			if httpTransport.client.Timeout != tt.config.Timeout {
				t.Errorf("NewTransport().client.Timeout = %v, want %v", httpTransport.client.Timeout, tt.config.Timeout)
			}
		})
	}
}

func TestHTTPTransport_Get_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("User-Agent") != "tensor-go-sdk/1.0.0" {
			t.Errorf("Expected User-Agent: tensor-go-sdk/1.0.0, got %s", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("x-tensor-api-key") != "test-api-key" {
			t.Errorf("Expected x-tensor-api-key: test-api-key, got %s", r.Header.Get("x-tensor-api-key"))
		}

		// Verify query parameters
		if r.URL.Query().Get("wallet") != "test-wallet" {
			t.Errorf("Expected wallet=test-wallet, got %s", r.URL.Query().Get("wallet"))
		}
		if r.URL.Query().Get("includeBidCount") != "true" {
			t.Errorf("Expected includeBidCount=true, got %s", r.URL.Query().Get("includeBidCount"))
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message": "success", "collections": []}`)
	}))
	defer server.Close()

	// Create transport
	transport := NewTransport(Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
		Timeout: 30 * time.Second,
	})

	// Prepare query parameters
	params := url.Values{}
	params.Set("wallet", "test-wallet")
	params.Set("includeBidCount", "true")

	// Make request
	ctx := context.Background()
	resp, err := transport.Get(ctx, "/api/v1/user/portfolio", params)

	// Verify results
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	resp.Body.Close()
}

func TestHTTPTransport_Get_WithoutAPIKey(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that API key header is not set
		if r.Header.Get("x-tensor-api-key") != "" {
			t.Errorf("Expected no API key header, got %s", r.Header.Get("x-tensor-api-key"))
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message": "success"}`)
	}))
	defer server.Close()

	// Create transport without API key
	transport := NewTransport(Config{
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
	})

	// Make request
	ctx := context.Background()
	resp, err := transport.Get(ctx, "/test", nil)

	// Verify results
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	resp.Body.Close()
}

func TestHTTPTransport_Get_WithoutQueryParams(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify no query parameters
		if len(r.URL.Query()) != 0 {
			t.Errorf("Expected no query parameters, got %v", r.URL.Query())
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message": "success"}`)
	}))
	defer server.Close()

	// Create transport
	transport := NewTransport(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	})

	// Make request without query parameters
	ctx := context.Background()
	resp, err := transport.Get(ctx, "/test", nil)

	// Verify results
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	resp.Body.Close()
}

func TestHTTPTransport_Get_APIErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedErrMsg string
	}{
		{
			name:           "400 Bad Request with JSON error",
			statusCode:     400,
			responseBody:   `{"message": "Invalid wallet address", "details": "Wallet format is incorrect"}`,
			expectedErrMsg: "API error 400: Invalid wallet address (Wallet format is incorrect)",
		},
		{
			name:           "401 Unauthorized",
			statusCode:     401,
			responseBody:   `{"message": "Invalid API key"}`,
			expectedErrMsg: "API error 401: Invalid API key",
		},
		{
			name:           "422 Validation Error",
			statusCode:     422,
			responseBody:   `{"error": "Validation failed", "details": "Required field missing"}`,
			expectedErrMsg: "API error 422: Validation failed (Required field missing)",
		},
		{
			name:           "429 Rate Limit",
			statusCode:     429,
			responseBody:   `{"message": "Too many requests"}`,
			expectedErrMsg: "API error 429: Too many requests",
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     500,
			responseBody:   `{"message": "Internal server error"}`,
			expectedErrMsg: "API error 500: Internal server error",
		},
		{
			name:           "404 with default message",
			statusCode:     404,
			responseBody:   `{}`,
			expectedErrMsg: "API error 404: Not Found",
		},
		{
			name:           "503 with invalid JSON",
			statusCode:     503,
			responseBody:   `invalid json`,
			expectedErrMsg: "API error 503: Service Unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that returns the error
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer server.Close()

			// Create transport
			transport := NewTransport(Config{
				BaseURL: server.URL,
				APIKey:  "test-key",
				Timeout: 30 * time.Second,
			})

			// Make request
			ctx := context.Background()
			resp, err := transport.Get(ctx, "/test", nil)

			// Verify error
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if resp != nil {
				t.Error("Expected nil response on error")
			}

			// Check if it's an API error
			var apiErr *apierrors.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("Expected APIError, got %T: %v", err, err)
			}

			if apiErr.Error() != tt.expectedErrMsg {
				t.Errorf("Expected error message %q, got %q", tt.expectedErrMsg, apiErr.Error())
			}
		})
	}
}

func TestHTTPTransport_Get_NetworkErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupServer func() *httptest.Server
		expectedErr string
	}{
		{
			name: "connection refused",
			setupServer: func() *httptest.Server {
				// Create server but close it immediately to simulate connection refused
				server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
				server.Close()
				return server
			},
			expectedErr: "network error during http_request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()

			// Create transport
			transport := NewTransport(Config{
				BaseURL: server.URL,
				APIKey:  "test-key",
				Timeout: 1 * time.Second, // Short timeout for faster tests
			})

			// Make request
			ctx := context.Background()
			resp, err := transport.Get(ctx, "/test", nil)

			// Verify error
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if resp != nil {
				t.Error("Expected nil response on error")
			}

			// Check if it's a network error
			var netErr *apierrors.NetworkError
			if !errors.As(err, &netErr) {
				t.Fatalf("Expected NetworkError, got %T: %v", err, err)
			}

			if !strings.Contains(netErr.Error(), tt.expectedErr) {
				t.Errorf("Expected error to contain %q, got %q", tt.expectedErr, netErr.Error())
			}
		})
	}
}

func TestHTTPTransport_Get_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message": "success"}`)
	}))
	defer server.Close()

	// Create transport
	transport := NewTransport(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	})

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Make request
	resp, err := transport.Get(ctx, "/test", nil)

	// Verify error
	if err == nil {
		t.Fatal("Expected error due to context cancellation, got nil")
	}
	if resp != nil {
		t.Error("Expected nil response on error")
	}

	// Check if it's a network error with context cancellation
	var netErr *apierrors.NetworkError
	if !errors.As(err, &netErr) {
		t.Fatalf("Expected NetworkError, got %T: %v", err, err)
	}

	if !strings.Contains(netErr.Error(), "context deadline exceeded") {
		t.Errorf("Expected context deadline exceeded error, got %q", netErr.Error())
	}
}

func TestHTTPTransport_Get_InvalidURL(t *testing.T) {
	// Create transport with invalid characters that will cause URL parsing to fail
	transport := NewTransport(Config{
		BaseURL: "http://example.com",
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	})

	// Make request with invalid path that contains characters that break URL parsing
	ctx := context.Background()
	resp, err := transport.Get(ctx, "/test\x00invalid", nil)

	// Verify error
	if err == nil {
		t.Fatal("Expected error due to invalid URL, got nil")
	}
	if resp != nil {
		t.Error("Expected nil response on error")
	}

	// Check if it's a network error
	var netErr *apierrors.NetworkError
	if !errors.As(err, &netErr) {
		t.Fatalf("Expected NetworkError, got %T: %v", err, err)
	}

	if netErr.Op != "create_request" {
		t.Errorf("Expected operation 'create_request', got %q", netErr.Op)
	}
}
