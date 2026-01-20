package kyt

import "context"

// Provider defines the unified interface for KYT service providers.
// Each provider implementation (Beosin, Chainalysis, etc.) must implement this interface.
type Provider interface {
	// Name returns the provider name (e.g., "beosin", "chainalysis").
	Name() string

	// AddressRisk performs risk assessment on a blockchain address.
	AddressRisk(ctx context.Context, req *AddressRiskRequest) (*RiskResult, error)

	// DepositRisk performs risk assessment on an incoming/deposit transaction.
	DepositRisk(ctx context.Context, req *TransactionRiskRequest) (*RiskResult, error)

	// WithdrawRisk performs risk assessment on an outgoing/withdrawal transaction.
	WithdrawRisk(ctx context.Context, req *TransactionRiskRequest) (*RiskResult, error)

	// Close releases any resources held by the provider.
	Close() error
}

// ProviderInfo contains information about a provider.
type ProviderInfo struct {
	// Name is the provider identifier.
	Name string

	// DisplayName is the human-readable name.
	DisplayName string

	// Description provides details about the provider.
	Description string

	// SupportedChains lists the blockchain IDs this provider supports.
	SupportedChains []string
}
