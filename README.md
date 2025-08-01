# Tensor Go SDK (MVP)

A minimal Go SDK for interacting with the [Tensor API](https://docs.tensor.trade/), built with clean architecture and designed for future growth.

> ⚠️ This is an MVP version. Currently only the `/user/portfolio` endpoint is supported.

---

## 📦 Installation

```bash
go get github.com/srpvpn/tensor-go-sdk
```
🚀 Quick Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/srpvpn/tensor-go-sdk/api/user"
    "github.com/srpvpn/tensor-go-sdk/client"
)

func main() {
    cfg := &client.Config{
        APIKey:  "your-api-key",
        Timeout: 30 * time.Second,
    }

    tensorClient := client.New(cfg)

    req := &user.PortfolioRequest{
        Wallet: "DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy",
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    resp, err := tensorClient.User.GetPortfolio(ctx, req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d collections\n", len(resp.Collections))
}
```

✅ Features

Clean architecture with modular packages

Support for context.Context and timeouts

Typed models for request and response

Custom error types: APIError, NetworkError, ValidationError

📚 Supported Endpoints

    GET /user/portfolio

🧪 Running Tests
```go
go test ./...
```
🔜 Roadmap

TSwap pool support

Retry logic with backoff

WebSocket data stream

Caching layer

More API endpoints

📝 License

This project is licensed under the GNU General Public License v3.0 – see LICENSE for details.

Built with Go and ❤️ for the Solana NFT ecosystem