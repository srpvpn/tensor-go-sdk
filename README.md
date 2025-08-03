# Tensor Go SDK

A Go SDK for interacting with the [Tensor API](https://docs.tensor.trade/)
## 🎯 Supported APIs

- **User API**: Portfolio, listings, bids, transactions, pools, escrow accounts, and inventory
- **Marketplace API**: NFT buying transactions

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

    "github.com/srpvpn/tensor-go-sdk/api/marketplace"
    "github.com/srpvpn/tensor-go-sdk/api/user"
    "github.com/srpvpn/tensor-go-sdk/client"
)

func main() {
    cfg := &client.Config{
        APIKey:  "your-api-key",
        Timeout: 30 * time.Second,
    }

    tensorClient := client.New(cfg)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // User API Example
    portfolioReq := &user.PortfolioRequest{
        Wallet: "DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy",
    }

    portfolioResp, _, err := tensorClient.User.GetPortfolio(ctx, portfolioReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d collections\n", len(portfolioResp))

    // Marketplace API Example
    buyReq := &marketplace.BuyNFTRequest{
        Buyer:     "buyer-wallet-address",
        Mint:      "nft-mint-address", 
        Owner:     "current-owner-address",
        MaxPrice:  1.5,
        Blockhash: "recent-blockhash",
    }

    buyResp, _, err := tensorClient.Marketplace.BuyNFT(ctx, buyReq)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated %d transactions\n", len(buyResp.Txs))
}
```

✅ Features

Clean architecture with modular packages

Support for context.Context and timeouts

Typed models for request and response

Custom error types: APIError, NetworkError, ValidationError

📚 Supported Endpoints

### User API
- `GET /api/v1/user/portfolio` - Get user's NFT portfolio
- `GET /api/v1/user/active_listings` - Get user's active listings
- `GET /api/v1/user/nft_bids` - Get user's single NFT bids
- `GET /api/v1/user/coll_bids` - Get user's collection bids
- `GET /api/v1/user/trait_bids` - Get user's trait bids
- `GET /api/v1/user/transactions` - Get user's NFT transactions
- `GET /api/v1/user/amm_pools` - Get user's TSwap pools
- `GET /api/v1/user/tamm_pools` - Get user's TAmm pools
- `GET /api/v1/user/escrow_accounts` - Get user's escrow accounts
- `GET /api/v1/user/inventory` - Get user's inventory for collection

### Marketplace API
- `GET /api/v1/tx/buy` - Create NFT purchase transaction
- `GET /api/v1/tx/sell` - Create NFT sell transaction (accept bid)
- `GET /api/v1/tx/list` - Create NFT listing transaction
- `GET /api/v1/tx/delist` - Create NFT delisting transaction
- `GET /api/v1/tx/edit` - Create NFT listing edit transaction
- `GET /api/v1/tx/bid` - Create single NFT bid transaction
- `GET /api/v1/tx/trait_bid` - Create trait bid transaction
- `GET /api/v1/tx/collection_bid` - Create collection bid transaction
- `GET /api/v1/tx/edit_bid` - Create bid edit transaction
- `GET /api/v1/tx/cancel_bid` - Create bid cancellation transaction

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