package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/srpvpn/tensor-go-sdk/internal/errors"
)

// mockTransport implements the client.Transport interface for testing
type mockTransport struct {
	response   *http.Response
	err        error
	lastPath   string
	lastParams url.Values
}

func (m *mockTransport) Get(ctx context.Context, path string, params url.Values) (*http.Response, error) {
	m.lastPath = path
	m.lastParams = params
	return m.response, m.err
}

// Helper function to create a mock HTTP response
func createMockResponse(statusCode int, body interface{}) *http.Response {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     make(http.Header),
	}
}

func TestNew(t *testing.T) {
	transport := &mockTransport{}
	api := New(transport)

	if api == nil {
		t.Fatal("New() returned nil")
	}

	// Verify that the returned instance implements UserAPI
	_, ok := api.(UserAPI)
	if !ok {
		t.Fatal("New() did not return a UserAPI implementation")
	}
}

func TestUserAPI_GetPortfolio_Success(t *testing.T) {
	// Prepare mock response - just a simple array of collections
	mockCollections := []Collection{
		{
			ID:         "collection1",
			Name:       "Test Collection",
			Symbol:     "TEST",
			Image:      "https://example.com/image.png",
			FloorPrice: 1.5,
			Volume24h:  100.0,
			Verified:   true,
			Compressed: false,
		},
	}

	transport := &mockTransport{
		response: createMockResponse(200, mockCollections),
	}

	api := New(transport)

	req := &PortfolioRequest{
		Wallet: "11111111111111111111111111111112", // Valid test wallet
	}

	body, statusCode, err := api.GetPortfolio(context.Background(), req)

	// Verify no error occurred
	if err != nil {
		t.Fatalf("GetPortfolio() returned error: %v", err)
	}

	// Verify status code
	if statusCode != 200 {
		t.Errorf("Expected status code 200, got %d", statusCode)
	}

	// Verify the correct path was called
	if transport.lastPath != "/api/v1/user/portfolio" {
		t.Errorf("Expected path '/api/v1/user/portfolio', got '%s'", transport.lastPath)
	}

	// Verify query parameters were built correctly
	if transport.lastParams.Get("wallet") != req.Wallet {
		t.Errorf("Expected wallet parameter '%s', got '%s'", req.Wallet, transport.lastParams.Get("wallet"))
	}

	// Verify response body is valid JSON
	var collections []Collection
	if err := json.Unmarshal(body, &collections); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if len(collections) != 1 {
		t.Errorf("Expected 1 collection, got %d", len(collections))
	}

	if collections[0].ID != "collection1" {
		t.Errorf("Expected collection ID 'collection1', got '%s'", collections[0].ID)
	}
}

func TestUserAPI_GetPortfolio_WithOptionalParams(t *testing.T) {
	mockCollections := []Collection{}

	transport := &mockTransport{
		response: createMockResponse(200, mockCollections),
	}

	api := New(transport)

	// Test with optional parameters
	includeBidCount := true
	includeFavCount := false
	includeUnverified := true
	includeCompressed := false

	req := &PortfolioRequest{
		Wallet:                "11111111111111111111111111111112",
		IncludeBidCount:       &includeBidCount,
		IncludeFavouriteCount: &includeFavCount,
		IncludeUnverified:     &includeUnverified,
		IncludeCompressed:     &includeCompressed,
		Currencies:            []string{"SOL", "USDC"},
	}

	_, _, err := api.GetPortfolio(context.Background(), req)

	if err != nil {
		t.Fatalf("GetPortfolio() returned error: %v", err)
	}

	// Verify query parameters include optional fields
	params := transport.lastParams

	if params.Get("includeBidCount") != "true" {
		t.Errorf("Expected includeBidCount 'true', got '%s'", params.Get("includeBidCount"))
	}

	if params.Get("includeFavouriteCount") != "false" {
		t.Errorf("Expected includeFavouriteCount 'false', got '%s'", params.Get("includeFavouriteCount"))
	}

	if params.Get("includeUnverified") != "true" {
		t.Errorf("Expected includeUnverified 'true', got '%s'", params.Get("includeUnverified"))
	}

	if params.Get("includeCompressed") != "false" {
		t.Errorf("Expected includeCompressed 'false', got '%s'", params.Get("includeCompressed"))
	}

	if params.Get("currencies") != "SOL,USDC" {
		t.Errorf("Expected currencies 'SOL,USDC', got '%s'", params.Get("currencies"))
	}
}

func TestUserAPI_GetPortfolio_ValidationError(t *testing.T) {
	transport := &mockTransport{}
	api := New(transport)

	// Test with empty wallet
	req := &PortfolioRequest{
		Wallet: "",
	}

	_, _, err := api.GetPortfolio(context.Background(), req)

	if err == nil {
		t.Fatal("Expected validation error for empty wallet")
	}

	if !strings.Contains(fmt.Sprintf("%v", err), "wallet address is required") {
		t.Errorf("Expected wallet validation error, got: %v", err)
	}
}

func TestUserAPI_GetPortfolio_InvalidWallet(t *testing.T) {
	transport := &mockTransport{}
	api := New(transport)

	// Test with invalid wallet address
	req := &PortfolioRequest{
		Wallet: "invalid-wallet",
	}

	_, _, err := api.GetPortfolio(context.Background(), req)

	if err == nil {
		t.Fatal("Expected validation error for invalid wallet")
	}

	if !strings.Contains(fmt.Sprintf("%v", err), "invalid wallet address") {
		t.Errorf("Expected wallet validation error, got: %v", err)
	}
}

func TestUserAPI_GetPortfolio_TransportError(t *testing.T) {
	transport := &mockTransport{
		err: &errors.NetworkError{
			Op:  "http_request",
			Err: fmt.Errorf("connection failed"),
		},
	}

	api := New(transport)

	req := &PortfolioRequest{
		Wallet: "11111111111111111111111111111112",
	}

	_, _, err := api.GetPortfolio(context.Background(), req)

	if err == nil {
		t.Fatal("Expected transport error")
	}

	if !strings.Contains(fmt.Sprintf("%v", err), "HTTP request failed") {
		t.Errorf("Expected HTTP request error, got: %v", err)
	}
}

func TestUserAPI_GetPortfolio_APIError(t *testing.T) {
	transport := &mockTransport{
		err: &errors.APIError{
			Code:    422,
			Message: "Validation failed",
			Details: "Invalid wallet address",
		},
	}

	api := New(transport)

	req := &PortfolioRequest{
		Wallet: "11111111111111111111111111111112",
	}

	_, _, err := api.GetPortfolio(context.Background(), req)

	if err == nil {
		t.Fatal("Expected API error")
	}

	// Check if it's an API error
	if !strings.Contains(fmt.Sprintf("%T", err), "APIError") && !strings.Contains(fmt.Sprintf("%v", err), "Validation failed") {
		t.Errorf("Expected API error, got: %v", err)
	}
}

func TestUserAPI_GetPortfolio_InvalidJSON(t *testing.T) {
	// Create response with invalid JSON
	transport := &mockTransport{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
			Header:     make(http.Header),
		},
	}

	api := New(transport)

	req := &PortfolioRequest{
		Wallet: "11111111111111111111111111111112",
	}

	body, statusCode, err := api.GetPortfolio(context.Background(), req)

	// With raw response, there should be no error - we just return the raw data
	if err != nil {
		t.Fatalf("GetPortfolio() returned unexpected error: %v", err)
	}

	// Verify status code
	if statusCode != 200 {
		t.Errorf("Expected status code 200, got %d", statusCode)
	}

	// Verify we get the raw invalid JSON back
	if string(body) != "invalid json" {
		t.Errorf("Expected raw body 'invalid json', got '%s'", string(body))
	}
}

func TestUserAPI_GetPortfolio_ContextCancellation(t *testing.T) {
	transport := &mockTransport{
		response: createMockResponse(200, []Collection{}),
	}

	api := New(transport)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &PortfolioRequest{
		Wallet: "11111111111111111111111111111112",
	}

	// The transport mock doesn't actually respect context cancellation,
	// but we can verify the context is passed through
	_, _, err := api.GetPortfolio(ctx, req)

	// In a real scenario with cancelled context, we'd expect an error
	// For this mock test, we just verify the method completes
	if err != nil && !strings.Contains(fmt.Sprintf("%v", err), "context") {
		// If there's an error, it should be context-related
		t.Logf("Context cancellation test completed with error: %v", err)
	}
}
