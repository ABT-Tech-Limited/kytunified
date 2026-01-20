package beosin

import (
	"time"

	"github.com/ABT-Tech-Limited/beosin-go"
	"github.com/ABT-Tech-Limited/kytunified/kyt"
)

// Mapper handles conversion between Beosin V2/V3 responses and unified KYT responses.
type Mapper struct {
	calculator *RiskCalculator
}

// NewMapper creates a new Mapper.
func NewMapper() *Mapper {
	return &Mapper{
		calculator: NewRiskCalculator(),
	}
}

// MapAddressRisk converts Beosin AddressRiskResponse (V3) to unified RiskResult.
func (m *Mapper) MapAddressRisk(resp *beosin.AddressRiskResponse) *kyt.RiskResult {
	if resp == nil || resp.Data == nil {
		return &kyt.RiskResult{
			Level:    kyt.RiskLevelUnknown,
			Metadata: m.buildMetadata("v3"),
		}
	}

	data := resp.Data

	// Extract tags and factors
	tags := m.extractAddressTags(data)
	factors := m.extractAddressRiskFactors(data)

	// Calculate unified risk level
	riskLevel := m.calculator.CalculateRiskLevel(data.Score, factors, tags)

	return &kyt.RiskResult{
		Level:    riskLevel,
		Score:    data.Score,
		Metadata: m.buildMetadata("v3"),
		Detail: &kyt.Detail{
			Factors: factors,
			Tags:    tags,
			IncomingRisk: &kyt.DirectionalRisk{
				Level:   m.mapBeosinRiskLevel(data.IncomingLevel),
				Score:   data.IncomingScore,
				Factors: m.extractStrategyRiskFactors(data.IncomingDetail),
			},
			OutgoingRisk: &kyt.DirectionalRisk{
				Level:   m.mapBeosinRiskLevel(data.OutgoingLevel),
				Score:   data.OutgoingScore,
				Factors: m.extractStrategyRiskFactors(data.OutgoingDetail),
			},
		},
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

	data := resp.Data
	factors := m.extractTransactionRiskFactors(data.Risks)
	tags := m.extractTransactionTags(data.Risks)

	// Calculate unified risk level
	riskLevel := m.calculator.CalculateTransactionRiskLevel(data.Score, factors)

	return &kyt.RiskResult{
		Level:    riskLevel,
		Score:    data.Score,
		Metadata: m.buildMetadata("v2"),
		Detail: &kyt.Detail{
			Factors: factors,
			Tags:    tags,
		},
	}
}

// mapBeosinRiskLevel maps Beosin risk level string to unified RiskLevel.
func (m *Mapper) mapBeosinRiskLevel(level string) kyt.RiskLevel {
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

// extractAddressTags extracts risk tags from address data.
func (m *Mapper) extractAddressTags(data *beosin.AddressRiskData) []string {
	if data == nil {
		return nil
	}
	return data.RiskTagDetails
}

// extractAddressRiskFactors extracts all risk factors from address data.
func (m *Mapper) extractAddressRiskFactors(data *beosin.AddressRiskData) []kyt.RiskFactor {
	if data == nil {
		return nil
	}

	var factors []kyt.RiskFactor

	// Add risk tag details as factors
	for _, tag := range data.RiskTagDetails {
		factors = append(factors, kyt.RiskFactor{
			Category: tag,
			Severity: m.calculator.GetCategorySeverity(tag),
		})
	}

	// Add incoming risk factors
	for _, detail := range data.IncomingDetail {
		factors = append(factors, m.strategyDetailToFactors(detail)...)
	}

	// Add outgoing risk factors
	for _, detail := range data.OutgoingDetail {
		factors = append(factors, m.strategyDetailToFactors(detail)...)
	}

	return factors
}

// strategyDetailToFactors converts a strategy detail to risk factors.
func (m *Mapper) strategyDetailToFactors(detail beosin.StrategyRiskDetail) []kyt.RiskFactor {
	var factors []kyt.RiskFactor
	for _, rd := range detail.RiskDetails {
		factors = append(factors, kyt.RiskFactor{
			Category:    rd.RiskName,
			Severity:    m.calculator.GetCategorySeverity(rd.RiskName),
			Rate:        rd.Rate,
			Amount:      rd.Amount,
			Description: detail.StrategyName,
		})
	}
	return factors
}

// extractStrategyRiskFactors extracts factors from strategy details.
func (m *Mapper) extractStrategyRiskFactors(details []beosin.StrategyRiskDetail) []kyt.RiskFactor {
	var factors []kyt.RiskFactor
	for _, detail := range details {
		factors = append(factors, m.strategyDetailToFactors(detail)...)
	}
	return factors
}

// extractTransactionRiskFactors extracts risk factors from Beosin risks.
func (m *Mapper) extractTransactionRiskFactors(risks []beosin.Risk) []kyt.RiskFactor {
	var factors []kyt.RiskFactor
	for _, risk := range risks {
		for _, detail := range risk.RiskDetails {
			factors = append(factors, kyt.RiskFactor{
				Category:    detail.RiskName,
				Severity:    m.calculator.GetCategorySeverity(detail.RiskName),
				Rate:        detail.Rate,
				Amount:      detail.Amount,
				Description: risk.RiskStrategy,
			})
		}
	}
	return factors
}

// extractTransactionTags extracts unique tags from transaction risks.
func (m *Mapper) extractTransactionTags(risks []beosin.Risk) []string {
	tagSet := make(map[string]struct{})
	for _, risk := range risks {
		for _, detail := range risk.RiskDetails {
			tagSet[detail.RiskName] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags
}

// buildMetadata creates response metadata.
func (m *Mapper) buildMetadata(apiVersion string) *kyt.Metadata {
	return &kyt.Metadata{
		Provider:    ProviderName,
		ProcessedAt: time.Now().UTC(),
		APIVersion:  apiVersion,
	}
}
