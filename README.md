# 🚀 Tensor Go SDK

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue?style=for-the-badge)](LICENSE)


**A powerful, type-safe Go SDK for the [Tensor Protocol](https://docs.tensor.trade/) - The leading NFT marketplace on Solana**

[Installation](#-installation) •
[Quick Start](#-quick-start) •
[Examples](#-examples) •
[Contributing](#-contributing)

</div>

---

## ✨ Features
- 🔥 **Complete API Coverage** - Support for all major Tensor endpoints
- 🛡️ **Type Safety** - Fully typed requests and responses
- ⚡ **High Performance** - Optimized for speed and efficiency
- 🔄 **Context Support** - Built-in timeout and cancellation support
- 📝 **Rich Error Handling** - Detailed error types and messages
- 🧪 **Well Tested** - Comprehensive test coverage
- 📚 **Great Documentation** - Clear examples and API docs
- 🔧 **Easy Integration** - Simple, intuitive API design

## 📦 Installation

```bash
go get github.com/srpvpn/tensor-go-sdk
```

## 🚀 Quick Start

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
    // Initialize client
    cfg := &client.Config{
        APIKey:  "your-api-key",
        Timeout: 30 * time.Second,
    }
    
    tensorClient := client.New(cfg)
    ctx := context.Background()

    // Get user portfolio
    portfolio, _, err := tensorClient.User.GetPortfolio(ctx, &user.PortfolioRequest{
        Wallet: "DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d collections in portfolio\n", len(portfolio.Collections))

    // Buy an NFT
    buyTx, _, err := tensorClient.Marketplace.BuyNFT(ctx, &marketplace.BuyNFTRequest{
        Buyer:     "buyer-wallet-address",
        Mint:      "nft-mint-address",
        Owner:     "current-owner-address",
        MaxPrice:  1.5,
        Blockhash: "recent-blockhash",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Generated %d transactions for NFT purchase\n", len(buyTx.Txs))
}
```

## 📚 API Reference

### 👤 User API

<details>
<summary><b>Portfolio Management</b></summary>

```go
// Get user's NFT portfolio
portfolio, _, err := client.User.GetPortfolio(ctx, &user.PortfolioRequest{
    Wallet:                "wallet-address",
    IncludeBidCount:       &[]bool{true}[0],
    IncludeFavouriteCount: &[]bool{true}[0],
    IncludeUnverified:     &[]bool{false}[0],
    Currencies:            []string{"SOL", "USDC"},
})
```
</details>

<details>
<summary><b>Active Listings</b></summary>

```go
// Get user's active NFT listings
listings, _, err := client.User.GetListings(ctx, &user.ListingsRequest{
    Wallets:    []string{"wallet1", "wallet2"},
    SortBy:     "PriceDesc",
    Limit:      50,
    CollId:     &[]string{"collection-id"}[0],
    Currencies: []string{"SOL"},
})
```
</details>

<details>
<summary><b>Bid Management</b></summary>

```go
// Get NFT bids
nftBids, _, err := client.User.GetNFTBids(ctx, &user.NFTBidsRequest{
    Owner:  "wallet-address",
    Limit:  100,
    CollId: &[]string{"collection-id"}[0],
})

// Get collection bids
collBids, _, err := client.User.GetCollectionBids(ctx, &user.CollectionBidsRequest{
    Owner:  "wallet-address",
    Limit:  100,
    CollId: &[]string{"collection-id"}[0],
})

// Get trait bids
traitBids, _, err := client.User.GetTraitBids(ctx, &user.TraitBidsRequest{
    Owner:  "wallet-address",
    Limit:  100,
    CollId: &[]string{"collection-id"}[0],
})
```
</details>

<details>
<summary><b>Transaction History</b></summary>

```go
// Get user's transaction history
transactions, _, err := client.User.GetTransactions(ctx, &user.TransactionsRequest{
    Wallets: []string{"wallet-address"},
    Limit:   100,
    TxTypes: []string{"SALE_BUY_NOW", "SALE_ACCEPT_BID", "LIST"},
    Collid:  "collection-id",
})
```
</details>

<details>
<summary><b>Pool Management</b></summary>

```go
// Get TSwap pools
tswapPools, _, err := client.User.GetTSwapPools(ctx, &user.TSwapsPoolsRequest{
    Owner:         "wallet-address",
    PoolAddresses: []string{"pool1", "pool2"},
    Limit:         50,
})

// Get TAmm pools
tammPools, _, err := client.User.GetTAmmPools(ctx, &user.TAmmPoolsRequest{
    Owner:         "wallet-address",
    PoolAddresses: []string{"pool1", "pool2"},
    Limit:         50,
})
```
</details>

<details>
<summary><b>Escrow & Inventory</b></summary>

```go
// Get escrow accounts
escrow, _, err := client.User.GetEscrowAccounts(ctx, &user.EscrowAccountsRequest{
    Owner: "wallet-address",
})

// Get inventory for collection
inventory, _, err := client.User.GetInventoryForCollection(ctx, &user.InventoryForCollectionRequest{
    Wallets: []string{"wallet-address"},
    CollId:  &[]string{"collection-id"}[0],
    Limit:   &[]int32{100}[0],
})
```
</details>

### 🛒 Marketplace API

<details>
<summary><b>NFT Trading</b></summary>

```go
// Buy NFT
buyTx, _, err := client.Marketplace.BuyNFT(ctx, &marketplace.BuyNFTRequest{
    Buyer:              "buyer-wallet",
    Mint:               "nft-mint",
    Owner:              "current-owner",
    MaxPrice:           1.5,
    Blockhash:          "recent-blockhash",
    OptionalRoyaltyPct: &[]int32{5}[0],
})

// Sell NFT (accept bid)
sellTx, _, err := client.Marketplace.SellNFT(ctx, &marketplace.SellNFTRequest{
    Seller:     "seller-wallet",
    Mint:       "nft-mint",
    BidAddress: "bid-address",
    MinPrice:   1.0,
    Blockhash:  "recent-blockhash",
})
```
</details>

<details>
<summary><b>Listing Management</b></summary>

```go
// List NFT
listTx, _, err := client.Marketplace.ListNFT(ctx, &marketplace.ListNFTRequest{
    Mint:      "nft-mint",
    Owner:     "owner-wallet",
    Price:     2.5,
    Blockhash: "recent-blockhash",
    ExpireIn:  &[]int32{3600}[0], // 1 hour
})

// Edit listing
editTx, _, err := client.Marketplace.EditListing(ctx, &marketplace.EditListingRequest{
    Mint:      "nft-mint",
    Owner:     "owner-wallet",
    Price:     3.0, // New price
    Blockhash: "recent-blockhash",
})

// Delist NFT
delistTx, _, err := client.Marketplace.DelistNFT(ctx, &marketplace.DelistNFTRequest{
    Mint:      "nft-mint",
    Owner:     "owner-wallet",
    Blockhash: "recent-blockhash",
})
```
</details>

<details>
<summary><b>Bidding</b></summary>

```go
// Place NFT bid
nftBidTx, _, err := client.Marketplace.PlaceNFTBid(ctx, &marketplace.PlaceNFTBidRequest{
    Owner:           "bidder-wallet",
    Price:           1.5,
    Mint:            "nft-mint",
    Blockhash:       "recent-blockhash",
    UseSharedEscrow: &[]bool{true}[0],
})

// Place collection bid
collBidTx, _, err := client.Marketplace.PlaceCollectionBid(ctx, &marketplace.PlaceCollectionBidRequest{
    Owner:     "bidder-wallet",
    Price:     1.0,
    Quantity:  5,
    CollId:    "collection-id",
    Blockhash: "recent-blockhash",
})

// Place trait bid
traitBidTx, _, err := client.Marketplace.PlaceTraitBid(ctx, &marketplace.PlaceTraitBidRequest{
    Owner:     "bidder-wallet",
    Price:     0.8,
    Quantity:  3,
    CollId:    "collection-id",
    Traits:    []string{"trait1", "trait2"},
    Blockhash: "recent-blockhash",
})
```
</details>

<details>
<summary><b>Bid Management</b></summary>

```go
// Edit bid
editBidTx, _, err := client.Marketplace.EditBid(ctx, &marketplace.EditBidRequest{
    BidStateAddress: "bid-state-address",
    Blockhash:       "recent-blockhash",
    Price:           &[]float64{2.0}[0], // New price
    Quantity:        &[]int32{10}[0],    // New quantity
})

// Cancel bid
cancelTx, _, err := client.Marketplace.CancelBid(ctx, &marketplace.CancelBidRequest{
    BidStateAddress: "bid-state-address",
    Blockhash:       "recent-blockhash",
})
```
</details>

## 🎯 Implementation Status

### ✅ Implemented APIs

| API Category | Status | Endpoints |
|-------------|--------|-----------|
| **User API** | ✅ Complete | Portfolio, Listings, Bids, Transactions, Pools, Escrow, Inventory |
| **Marketplace API** | ✅ Complete | Buy, Sell, List, Delist, Edit, Bid, Cancel |

### 🚧 Roadmap

| API Category | Status | Priority | Description |
|-------------|--------|----------|-------------|
| **TSwap API** | 🔄 In Progress | High | AMM pool operations and swaps |
| **Shared Escrow API** | 📋 Planned | High | Escrow account management |
| **TAmm API** | 📋 Planned | Medium | Advanced AMM features |
| **Data API - NFTs** | 📋 Planned | Medium | NFT metadata and analytics |
| **Data API - Orders** | 📋 Planned | Medium | Order book and market data |
| **Data API - Collections** | 📋 Planned | Medium | Collection statistics |
| **Data API - RPC** | 📋 Planned | Low | Direct RPC calls |
| **Refresh API** | 📋 Planned | Low | Data refresh endpoints |
| **SDK API - Mint Proof** | 📋 Planned | Medium | Mint proof generation |
| **SDK API - Trait Bids** | 📋 Planned | Medium | Advanced trait bidding |
| **SDK API - Whitelist** | 📋 Planned | Low | Whitelist management |

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./api/user
go test ./api/marketplace
```

## 📖 Examples

Check out the [examples](./examples) directory for complete working examples:

- [Basic Usage](./examples/basic_usage/main.go) - Simple portfolio and trading operations
- [Advanced Trading](./examples/offline_demo/main.go) - Complex trading scenarios

## 🤝 Contributing

We welcome contributions! This project is growing fast and we'd love your help to make it even better.

### Ways to Contribute

- 🐛 **Report Bugs** - Found an issue? Let us know!
- 💡 **Feature Requests** - Have an idea? We'd love to hear it!
- 📝 **Documentation** - Help improve our docs
- 🔧 **Code Contributions** - Submit PRs for new features or fixes
- 🧪 **Testing** - Help us improve test coverage
- 🌟 **Spread the Word** - Star the repo and tell others!

### Getting Started

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for your changes
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Setup

```bash
# Clone the repo
git clone https://github.com/srpvpn/tensor-go-sdk.git
cd tensor-go-sdk

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter (if you have golangci-lint installed)
golangci-lint run
```

## 📄 License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.


<div align="center">

**Built with ❤️ for the Solana NFT ecosystem**

[⭐ Star us on GitHub](https://github.com/srpvpn/tensor-go-sdk) 

</div>