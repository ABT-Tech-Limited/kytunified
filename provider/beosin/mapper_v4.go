package beosin

import (
	"time"

	"github.com/ABT-Tech-Limited/beosin-go"
	"github.com/ABT-Tech-Limited/kytunified/kyt"
)

// MapperV4 handles conversion between Beosin V4 responses and unified KYT responses.
type MapperV4 struct{}

// NewMapperV4 creates a new V4 Mapper.
func NewMapperV4() *MapperV4 {
	return &MapperV4{}
}

// MapAddressRisk converts Beosin V4AddressRiskResponse to unified RiskResult.
func (m *MapperV4) MapAddressRisk(resp *beosin.V4AddressRiskResponse) *kyt.RiskResult {
	if resp == nil || resp.Data == nil {
		return &kyt.RiskResult{
			Level:    kyt.RiskLevelUnknown,
			Metadata: m.buildMetadata(),
		}
	}

	return &kyt.RiskResult{
		Level:    mapBeosinRiskLevel(resp.Data.RiskLevel),
		Score:    resp.Data.Score,
		Metadata: m.buildMetadata(),
		Detail:   resp.Data,
	}
}

// MapTransactionRisk converts Beosin V4TransactionRiskResponse to unified RiskResult.
func (m *MapperV4) MapTransactionRisk(resp *beosin.V4TransactionRiskResponse) *kyt.RiskResult {
	if resp == nil || resp.Data == nil {
		return &kyt.RiskResult{
			Level:    kyt.RiskLevelUnknown,
			Metadata: m.buildMetadata(),
		}
	}

	return &kyt.RiskResult{
		Level:    mapBeosinRiskLevel(resp.Data.RiskLevel),
		Score:    resp.Data.Score,
		Metadata: m.buildMetadata(),
		Detail:   resp.Data,
	}
}

// buildMetadata creates response metadata for V4 API.
func (m *MapperV4) buildMetadata() kyt.Metadata {
	return kyt.Metadata{
		Provider:    ProviderName,
		ProcessedAt: time.Now().UTC(),
		APIVersion:  "v4",
	}
}
