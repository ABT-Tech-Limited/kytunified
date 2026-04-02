package kyt

import "time"

// ============================================
// Request Types
// ============================================

// AddressRiskRequest represents a request to assess the risk of an address.
type AddressRiskRequest struct {
	// ChainID is the blockchain identifier (e.g., "1" for Ethereum, "56" for BSC).
	ChainID string `json:"chain_id"`

	// Address is the blockchain address to assess.
	Address string `json:"address"`

	// Token is the optional token contract address.
	Token *string `json:"token,omitempty"`
}

// TransactionRiskRequest represents a request to assess the risk of a transaction.
type TransactionRiskRequest struct {
	// ChainID is the blockchain identifier.
	ChainID string `json:"chain_id"`

	// TxHash is the transaction hash to assess.
	TxHash string `json:"tx_hash"`

	// Token is the optional token contract address.
	Token *string `json:"token,omitempty"`
}

// ============================================
// Response Types
// ============================================

// RiskResult is the unified response for all risk assessments.
type RiskResult struct {
	// Level is the unified risk level.
	Level RiskLevel `json:"level"`

	// Score is the normalized risk score (0-100).
	Score float64 `json:"score"`

	// Metadata contains provider-specific metadata.
	Metadata Metadata `json:"metadata"`

	// Detail contains the raw response data from the provider.
	Detail any `json:"detail,omitempty"`
}

// Metadata contains metadata about the response.
type Metadata struct {
	// Provider is the name of the KYT provider.
	Provider string `json:"provider"`

	// ProcessedAt is when the assessment was processed.
	ProcessedAt time.Time `json:"processed_at"`

	// RequestID is an optional request identifier for tracing.
	RequestID string `json:"request_id,omitempty"`

	// APIVersion indicates which API version was used.
	APIVersion string `json:"api_version,omitempty"`
}

// ============================================
// Helper Methods
// ============================================

// IsHighRisk returns true if the result indicates high or critical risk.
func (r *RiskResult) IsHighRisk() bool {
	return r.Level.IsHigherOrEqual(RiskLevelHigh)
}

// IsCritical returns true if the result indicates critical risk.
func (r *RiskResult) IsCritical() bool {
	return r.Level == RiskLevelCritical
}
