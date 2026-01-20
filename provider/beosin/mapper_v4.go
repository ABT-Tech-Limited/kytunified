package beosin

import (
	"time"

	"github.com/ABT-Tech-Limited/beosin-go"
	"github.com/ABT-Tech-Limited/kytunified/kyt"
)

// MapperV4 handles conversion between Beosin V4 responses and unified KYT responses.
type MapperV4 struct {
	calculator *RiskCalculator
}

// NewMapperV4 creates a new V4 Mapper.
func NewMapperV4() *MapperV4 {
	return &MapperV4{
		calculator: NewRiskCalculator(),
	}
}

// MapAddressRisk converts Beosin V4AddressRiskResponse to unified RiskResult.
func (m *MapperV4) MapAddressRisk(resp *beosin.V4AddressRiskResponse) *kyt.RiskResult {
	if resp == nil || resp.Data == nil {
		return &kyt.RiskResult{
			Level:    kyt.RiskLevelUnknown,
			Metadata: m.buildMetadata(),
		}
	}

	data := resp.Data

	// Extract tags from RiskTagDetails
	tags := data.RiskTagDetails

	// Extract factors from incoming and outgoing details
	var allFactors []kyt.RiskFactor
	incomingFactors := m.extractV4StrategyFactors(data.IncomingDetail)
	outgoingFactors := m.extractV4StrategyFactors(data.OutgoingDetail)
	allFactors = append(allFactors, incomingFactors...)
	allFactors = append(allFactors, outgoingFactors...)

	// Add risk tags as factors
	for _, tag := range tags {
		allFactors = append(allFactors, kyt.RiskFactor{
			Category: tag,
			Severity: m.calculator.GetCategorySeverity(tag),
		})
	}

	// Calculate unified risk level
	riskLevel := m.calculator.CalculateRiskLevel(data.Score, allFactors, tags)

	return &kyt.RiskResult{
		Level:    riskLevel,
		Score:    data.Score,
		Metadata: m.buildMetadata(),
		Detail: &kyt.Detail{
			Factors: allFactors,
			Tags:    tags,
			IncomingRisk: &kyt.DirectionalRisk{
				Level:   m.mapBeosinRiskLevel(data.IncomingLevel),
				Score:   data.IncomingScore,
				Factors: incomingFactors,
			},
			OutgoingRisk: &kyt.DirectionalRisk{
				Level:   m.mapBeosinRiskLevel(data.OutgoingLevel),
				Score:   data.OutgoingScore,
				Factors: outgoingFactors,
			},
		},
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

	data := resp.Data
	factors := m.extractV4RiskFactors(data.Risks)
	tags := m.extractV4Tags(data.Risks)

	// Calculate unified risk level
	riskLevel := m.calculator.CalculateTransactionRiskLevel(data.Score, factors)

	return &kyt.RiskResult{
		Level:    riskLevel,
		Score:    data.Score,
		Metadata: m.buildMetadata(),
		Detail: &kyt.Detail{
			Factors: factors,
			Tags:    tags,
		},
	}
}

// extractV4StrategyFactors extracts risk factors from V4 strategy details.
func (m *MapperV4) extractV4StrategyFactors(details []beosin.V4StrategyDetail) []kyt.RiskFactor {
	var factors []kyt.RiskFactor

	for _, detail := range details {
		// Create a factor for each entity detail
		for _, entity := range detail.EntityDetails {
			factor := kyt.RiskFactor{
				Category:    detail.StrategyName,
				Severity:    m.mapBeosinRiskLevel(detail.RiskLevel),
				Description: entity.EntityName,
				Rate:        detail.Rate,
				Amount:      detail.Amount,
				Hops:        detail.Hops,
				Exposure:    detail.Exposure,
			}

			if entity.Hops > 0 {
				factor.Hops = entity.Hops
			}

			factors = append(factors, factor)
		}

		// If no entity details, create a single factor for the strategy
		if len(detail.EntityDetails) == 0 {
			factors = append(factors, kyt.RiskFactor{
				Category: detail.StrategyName,
				Severity: m.mapBeosinRiskLevel(detail.RiskLevel),
				Rate:     detail.Rate,
				Amount:   detail.Amount,
				Hops:     detail.Hops,
				Exposure: detail.Exposure,
			})
		}
	}

	return factors
}

// extractV4RiskFactors extracts risk factors from V4 risks (for transactions).
func (m *MapperV4) extractV4RiskFactors(risks []beosin.V4Risk) []kyt.RiskFactor {
	var factors []kyt.RiskFactor

	for _, risk := range risks {
		// Create a factor for each entity detail
		for _, entity := range risk.EntityDetails {
			factor := kyt.RiskFactor{
				Category:    risk.RiskStrategy,
				Severity:    m.mapBeosinRiskLevel(risk.RiskLevel),
				Description: entity.EntityName,
				Rate:        risk.Rate,
				Amount:      risk.Amount,
				Hops:        risk.Hops,
				Exposure:    risk.Exposure,
			}

			if entity.Hops > 0 {
				factor.Hops = entity.Hops
			}

			factors = append(factors, factor)
		}

		// If no entity details, create a single factor for the risk
		if len(risk.EntityDetails) == 0 {
			factors = append(factors, kyt.RiskFactor{
				Category: risk.RiskStrategy,
				Severity: m.mapBeosinRiskLevel(risk.RiskLevel),
				Rate:     risk.Rate,
				Amount:   risk.Amount,
				Hops:     risk.Hops,
				Exposure: risk.Exposure,
			})
		}
	}

	return factors
}

// extractV4Tags extracts unique tags from V4 risks.
func (m *MapperV4) extractV4Tags(risks []beosin.V4Risk) []string {
	tagSet := make(map[string]struct{})

	for _, risk := range risks {
		if risk.RiskStrategy != "" {
			tagSet[risk.RiskStrategy] = struct{}{}
		}

		for _, entity := range risk.EntityDetails {
			if entity.EntityName != "" {
				tagSet[entity.EntityName] = struct{}{}
			}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags
}

// mapBeosinRiskLevel maps Beosin risk level string to unified RiskLevel.
func (m *MapperV4) mapBeosinRiskLevel(level string) kyt.RiskLevel {
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

// buildMetadata creates response metadata for V4 API.
func (m *MapperV4) buildMetadata() *kyt.Metadata {
	return &kyt.Metadata{
		Provider:    ProviderName,
		ProcessedAt: time.Now().UTC(),
		APIVersion:  "v4",
	}
}
