package client

import (
	"time"

	"github.com/srpvpn/tensor-go-sdk/api/escrow"
	"github.com/srpvpn/tensor-go-sdk/api/marketplace"
	"github.com/srpvpn/tensor-go-sdk/api/rpc"
	"github.com/srpvpn/tensor-go-sdk/api/tswap"
	"github.com/srpvpn/tensor-go-sdk/api/user"
	"github.com/srpvpn/tensor-go-sdk/internal/transport"
)

// Client is the main SDK client that provides access to all API endpoints
type Client struct {
	transport   transport.Transport
	User        user.UserAPI
	Marketplace marketplace.MarketplaceAPI
	TSwap       tswap.TSwapAPI
	RPC         rpc.RPCAPI
	Escrow      escrow.EscrowAPI
}

const (
	// DefaultBaseURL is the default Tensor API base URL
	DefaultBaseURL = "https://api.mainnet.tensordev.io"
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
)

// New creates a new Client with the provided configuration.
// If config is nil, default configuration will be used.
func New(config *Config) *Client {
	// Use default configuration if none provided
	if config == nil {
		config = &Config{
			BaseURL: DefaultBaseURL,
			Timeout: DefaultTimeout,
		}
	} else {
		// Create a copy to avoid mutating the original
		configCopy := *config
		config = &configCopy
	}

	// Set default values for missing configuration
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	// Create transport layer
	transport := NewTransport(*config)
	// Create User API with transport
	userAPI := user.New(transport)
	// Create Marketplace API with transport
	marketplaceAPI := marketplace.New(transport)
	// Create TSwap API with transport
	tswapAPI := tswap.New(transport)
	// Create RPC API with transport
	rpcAPI := rpc.New(transport)
	// Create Escrow API with transport
	escrowAPI := escrow.New(transport)

	// Return initialized client
	return &Client{
		transport:   transport,
		User:        userAPI,
		Marketplace: marketplaceAPI,
		TSwap:       tswapAPI,
		RPC:         rpcAPI,
		Escrow:      escrowAPI,
	}
}

// Close closes the client and releases any resources.
// Currently this is a no-op but provides future extensibility.
func (c *Client) Close() error {
	// Future: close HTTP connections, cleanup resources
	return nil
}
