// client/config.go
package client

import "time"

type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}
