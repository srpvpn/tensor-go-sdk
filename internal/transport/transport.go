package transport

import (
	"context"
	"net/http"
	"net/url"
)

// Transport defines the interface for making HTTP requests.
type Transport interface {
	Get(ctx context.Context, path string, params url.Values) (*http.Response, error)
}
