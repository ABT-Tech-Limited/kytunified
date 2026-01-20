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
	Token string `json:"token,omitempty"`
}

// TransactionRiskRequest represents a request to assess the risk of a transaction.
type TransactionRiskRequest struct {
	// ChainID is the blockchain identifier.
	ChainID string `json:"chain_id"`

	// TxHash is the transaction hash to assess.
	TxHash string `json:"tx_hash"`

	// Token is the optional token contract address.
	Token string `json:"token,omitempty"`
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
	Metadata *Metadata `json:"metadata"`

	// Detail contains detailed risk information from the provider.
	Detail *Detail `json:"detail,omitempty"`
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

// Detail contains detailed risk information from the provider.
type Detail struct {
	// Factors contains the detected risk factors.
	Factors []RiskFactor `json:"factors,omitempty"`

	// Tags contains risk tags associated with this address/transaction.
	Tags []string `json:"tags,omitempty"`

	// IncomingRisk represents risk from incoming transactions (address only).
	IncomingRisk *DirectionalRisk `json:"incoming_risk,omitempty"`

	// OutgoingRisk represents risk from outgoing transactions (address only).
	OutgoingRisk *DirectionalRisk `json:"outgoing_risk,omitempty"`
}

// RiskFactor represents a single risk factor detected during assessment.
type RiskFactor struct {
	// Category is the risk category (e.g., "Sanction", "Mixer", "Hacker").
	Category string `json:"category"`

	// Severity indicates the severity of this specific factor.
	Severity RiskLevel `json:"severity"`

	// Description provides additional context about the risk.
	Description string `json:"description,omitempty"`

	// Rate is the proportion of funds associated with this risk (0.0 - 1.0).
	Rate float64 `json:"rate,omitempty"`

	// Amount is the absolute amount associated with this risk.
	Amount float64 `json:"amount,omitempty"`

	// Hops indicates how many transaction hops away the risk source is.
	Hops int `json:"hops,omitempty"`

	// Exposure indicates whether this is direct or indirect exposure.
	Exposure string `json:"exposure,omitempty"`
}

// DirectionalRisk represents risk assessment for a specific direction.
type DirectionalRisk struct {
	// Level is the risk level for this direction.
	Level RiskLevel `json:"level"`

	// Score is the normalized risk score (0-100).
	Score float64 `json:"score"`

	// Factors contains the detected risk factors for this direction.
	Factors []RiskFactor `json:"factors,omitempty"`
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

// HasFactor returns true if the result contains a factor with the given category.
func (r *RiskResult) HasFactor(category string) bool {
	if r.Detail == nil {
		return false
	}
	for _, f := range r.Detail.Factors {
		if f.Category == category {
			return true
		}
	}
	return false
}

// HasTag returns true if the result contains the given tag.
func (r *RiskResult) HasTag(tag string) bool {
	if r.Detail == nil {
		return false
	}
	for _, t := range r.Detail.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
