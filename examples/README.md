# Tensor SDK Examples

This directory contains example programs demonstrating how to use the Tensor Go SDK.

## Available Examples

### 1. Basic Usage Example (`basic_usage/`)

The `basic_usage` example demonstrates the core functionality of the SDK with real API calls:

### Features Demonstrated

1. **Basic Client Creation**: Creating a client with default configuration
2. **Simple Portfolio Request**: Making a basic request to get user portfolio
3. **Advanced Configuration**: Using custom client configuration with timeouts
4. **Advanced Request Options**: Using all available request parameters
5. **Error Handling**: Proper error handling for various scenarios
6. **Context Management**: Using context for request cancellation and timeouts

### Running the Example

```bash
# From the project root directory
go run examples/basic_usage/main.go
```

### Example Output

The program will demonstrate:

- Fetching a portfolio for a sample wallet address
- Displaying collection information including floor prices and volumes
- Advanced usage with custom configuration
- Error handling for invalid inputs
- Context cancellation scenarios

### Configuration Options

The example shows how to configure the client:

```go
config := &client.Config{
    BaseURL: "https://api.mainnet.tensordev.io", // API base URL
    Timeout: 45 * time.Second,                   // Request timeout
    APIKey:  "your-api-key-here",               // Optional API key
}
```

### Request Parameters

The example demonstrates all available portfolio request parameters:

```go
request := &user.PortfolioRequest{
    Wallet:                "wallet-address-here",
    IncludeBidCount:       &includeBidCount,       // Include bid counts
    IncludeFavouriteCount: &includeFavCount,       // Include favorite counts
    IncludeUnverified:     &includeUnverified,     // Include unverified collections
    IncludeCompressed:     &includeCompressed,     // Include compressed NFTs
    Currencies:            []string{"SOL"},        // Filter by currencies
}
```

### Error Handling

The example shows proper error handling patterns:

- Validation errors for invalid wallet addresses
- Network errors and timeouts
- API errors from the server
- Context cancellation

### Best Practices Demonstrated

1. **Resource Management**: Proper use of `defer client.Close()`
2. **Context Usage**: Using context with timeouts for all requests
3. **Error Handling**: Comprehensive error handling with user-friendly messages
4. **Configuration**: Showing both default and custom configuration options
5. **Validation**: Client-side validation before making API requests

### 2. Offline Demo Example (`offline_demo/`)

The `offline_demo` example demonstrates SDK usage patterns without making real API calls, perfect for learning and testing:

#### Features Demonstrated

1. **Client Configuration**: Shows how to create and configure the SDK client
2. **Request Building**: Demonstrates building both basic and advanced requests
3. **Validation**: Shows client-side validation of request parameters
4. **Context Patterns**: Demonstrates proper context usage patterns
5. **Response Structure**: Shows the structure of API responses with sample data
6. **Error Handling Patterns**: Demonstrates comprehensive error handling approaches

#### Running the Offline Demo

```bash
# From the project root directory
go run examples/offline_demo/main.go
```

#### What You'll Learn

- How to structure requests with all available parameters
- Client-side validation patterns
- Response data structure and handling
- Context management for timeouts and cancellation
- Error handling strategies for different scenarios

This example is ideal for:
- Learning the SDK API without network dependencies
- Understanding data structures and validation
- Testing integration patterns in your own code
- Offline development and testing

## Environment Variables

You can set the following environment variables to modify the example behavior:

- `TENSOR_SDK_DEBUG=true`: Enable debug mode for additional logging

## Notes

- The example uses a well-known Solana wallet address for demonstration
- All requests include proper timeout handling
- The example is designed to work without requiring an API key
- Error scenarios are demonstrated to show proper error handling patterns