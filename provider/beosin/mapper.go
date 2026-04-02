package beosin

import (
	"time"

	"github.com/ABT-Tech-Limited/beosin-go"
	"github.com/ABT-Tech-Limited/kytunified/kyt"
)

// Mapper handles conversion between Beosin V2/V3 responses and unified KYT responses.
type Mapper struct{}

// NewMapper creates a new Mapper.
func NewMapper() *Mapper {
	return &Mapper{}
}

// MapAddressRisk converts Beosin AddressRiskResponse (V3) to unified RiskResult.
func (m *Mapper) MapAddressRisk(resp *beosin.AddressRiskResponse) *kyt.RiskResult {
	if resp == nil || resp.Data == nil {
		return &kyt.RiskResult{
			Level:    kyt.RiskLevelUnknown,
			Metadata: m.buildMetadata("v3"),
		}
	}

	return &kyt.RiskResult{
		Level:    mapBeosinRiskLevel(resp.Data.RiskLevel),
		Score:    resp.Data.Score,
		Metadata: m.buildMetadata("v3"),
		Detail:   resp.Data,
	}
}

// MapTransactionRisk converts Beosin TransactionRiskResponse (V2) to unified RiskResult.
func (m *Mapper) MapTransactionRisk(resp *beosin.TransactionRiskResponse) *kyt.RiskResult {
	if resp == nil || resp.Data == nil {
		return &kyt.RiskResult{
			Level:    kyt.RiskLevelUnknown,
			Metadata: m.buildMetadata("v2"),
		}
	}

	return &kyt.RiskResult{
		Level:    mapBeosinRiskLevel(resp.Data.RiskLevel),
		Score:    resp.Data.Score,
		Metadata: m.buildMetadata("v2"),
		Detail:   resp.Data,
	}
}

// mapBeosinRiskLevel maps Beosin risk level string to unified RiskLevel.
func mapBeosinRiskLevel(level string) kyt.RiskLevel {
	switch level {
	case beosin.RiskLevelSevere:
		return kyt.RiskLevelCritical
	case beosin.RiskLevelHigh:
		return kyt.RiskLevelHigh
	case beosin.RiskLevelMedium:
		return kyt.RiskLevelMedium
	case beosin.RiskLevelLow:
		return kyt.RiskLevelLow
	default:
		return kyt.RiskLevelUnknown
	}
}

// buildMetadata creates response metadata.
func (m *Mapper) buildMetadata(apiVersion string) kyt.Metadata {
	return kyt.Metadata{
		Provider:    ProviderName,
		ProcessedAt: time.Now().UTC(),
		APIVersion:  apiVersion,
	}
}
