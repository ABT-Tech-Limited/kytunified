# KYT Unified

A unified Know Your Transaction (KYT) framework for Go that abstracts multiple blockchain compliance providers behind a single interface.

## Features

- **Unified Provider Interface** — Single `kyt.Provider` interface for all KYT providers
- **Multi-API Version Support** — Beosin V2/V3/V4 APIs supported out of the box
- **Unified Error System** — Typed errors (Validation, Retryable, RateLimit, Provider) with `errors.Is`/`errors.As` support
- **Provider Registry** — Thread-safe registry for managing provider instances
- **18 Blockchains Supported** — Ethereum, Bitcoin, BSC, Tron, Solana, TON, Polygon, Arbitrum, Optimism, Avalanche, and more

## Installation

```bash
go get github.com/ABT-Tech-Limited/kytunified
```

## Quick Start

### Direct Provider Creation

```go
package main

import (
    "context"
    "fmt"

    "github.com/ABT-Tech-Limited/beosin-go"
    "github.com/ABT-Tech-Limited/kytunified/kyt"
    beosinprovider "github.com/ABT-Tech-Limited/kytunified/provider/beosin"
)

func main() {
    client := beosin.NewClient(appID, appSecret)
    provider := beosinprovider.New(client, beosinprovider.WithV4())
    defer provider.Close()

    result, err := provider.AddressRisk(context.Background(), &kyt.AddressRiskRequest{
        ChainID: kyt.ChainIDETH,
        Address: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
    })
    if err != nil {
        // Handle error
    }

    fmt.Printf("Risk Level: %s, Score: %.2f\n", result.Level, result.Score)
}
```

### Using the Registry

```go
import (
    "github.com/ABT-Tech-Limited/beosin-go"
    beosinprovider "github.com/ABT-Tech-Limited/kytunified/provider/beosin"
    "github.com/ABT-Tech-Limited/kytunified/registry"
)

// Register
client := beosin.NewClient(appID, appSecret)
registry.MustRegisterBeosin(client, beosinprovider.WithV4())

// Retrieve and use
provider, err := registry.GetBeosin()
if err != nil {
    // Handle error
}
defer provider.Close()
```

## Provider Interface

```go
type Provider interface {
    Name() string
    AddressRisk(ctx context.Context, req *AddressRiskRequest) (*RiskResult, error)
    DepositRisk(ctx context.Context, req *TransactionRiskRequest) (*RiskResult, error)
    WithdrawRisk(ctx context.Context, req *TransactionRiskRequest) (*RiskResult, error)
    Close() error
}
```

## Risk Levels

| Level    | Description                                          |
|----------|------------------------------------------------------|
| Low      | Minimal risk, no significant exposure                |
| Medium   | Moderate risk, some exposure to risky activities     |
| High     | Significant risk, substantial exposure detected      |
| Critical | Severe risk, direct exposure to sanctioned entities  |

## Supported Chains

**Full Support:** Bitcoin, Ethereum, BNB Smart Chain, Tron, Polygon, Solana, TON, Arbitrum One, Optimism, Avalanche C-Chain, Fantom Opera, Cronos, zkSync Era

**Basic Support:** Base, Linea, Scroll, Sui, Sonic

## Error Handling

Errors are classified into typed categories for easy handling:

```go
result, err := provider.AddressRisk(ctx, req)
if err != nil {
    switch {
    case kyt.IsRetryable(err):
        // Assessment in progress, retry later
    case kyt.IsValidation(err):
        // Bad input (invalid address, unsupported chain, etc.)
    case kyt.IsRateLimit(err):
        // Rate limited, back off and retry
    default:
        // Provider error
    }
}
```

## License

[MIT](LICENSE)
