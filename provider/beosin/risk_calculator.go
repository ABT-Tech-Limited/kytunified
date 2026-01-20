package beosin

import (
	"github.com/ABT-Tech-Limited/kytunified/kyt"
)

// RiskCalculator implements the Beosin-specific risk calculation rules.
// It converts Beosin's raw scores and tags into unified risk levels.
type RiskCalculator struct{}

// NewRiskCalculator creates a new RiskCalculator.
func NewRiskCalculator() *RiskCalculator {
	return &RiskCalculator{}
}

// CriticalTags are tags that immediately elevate risk to Critical.
// These represent severe threats like sanctions or direct involvement in attacks.
var CriticalTags = map[string]bool{
	"Sanction":          true,
	"Sanctioned":        true,
	"OFAC":              true,
	"Hacker":            true,
	"Terrorist":         true,
	"Ransomware":        true,
	"ChildAbuseMaterial": true,
}

// HighTags are tags that indicate high-risk activities.
var HighTags = map[string]bool{
	"Darknet":         true,
	"Theft":           true,
	"Scam":            true,
	"FraudShop":       true,
	"Drug":            true,
	"UndergroundBank": true,
	"MoneyMule":       true,
	"Trojan":          true,
}

// MediumTags are tags considered "medium severity".
var MediumTags = map[string]bool{
	"Mixing":            true,
	"Mixer":             true,
	"Gambling":          true,
	"HighRiskExchange":  true,
	"HighRiskJurisdictionFATF": true,
	"GreyListFATF":      true,
	"Piracy":            true,
	"ProtocolPiracy":    true,
}

// CalculateRiskLevel calculates unified risk level based on Beosin mapping rules.
//
// Mapping Rules (from requirements):
//   - Score 0-30 + no tags          -> Low
//   - Score 31-70 OR 1-2 medium tags -> Medium
//   - Score 71-99 OR >=3 tags        -> High
//   - Score 100 OR sanctioned/hacker -> Critical
//
// Parameters:
//   - score: The raw risk score from Beosin (0-100)
//   - factors: Extracted risk factors
//   - tags: Risk tags from the response
//
// Returns the unified RiskLevel.
func (c *RiskCalculator) CalculateRiskLevel(score float64, factors []kyt.RiskFactor, tags []string) kyt.RiskLevel {
	// Rule 1: Check for critical tags first (highest priority)
	if c.hasCriticalTags(tags) || c.hasCriticalFactors(factors) {
		return kyt.RiskLevelCritical
	}

	// Rule 2: Score of 100 is always Critical
	if score >= 100 {
		return kyt.RiskLevelCritical
	}

	// Count total significant tags/factors
	tagCount := len(tags)
	mediumTagCount := c.countMediumTags(tags)
	highTagCount := c.countHighTags(tags)

	// Rule 3: Score 71-99 OR >=3 tags -> High
	if score >= 71 || tagCount >= 3 || highTagCount > 0 {
		return kyt.RiskLevelHigh
	}

	// Rule 4: Score 31-70 OR 1-2 medium tags -> Medium
	if score >= 31 || (mediumTagCount >= 1 && mediumTagCount <= 2) {
		return kyt.RiskLevelMedium
	}

	// Rule 5: Score 0-30 + no significant tags -> Low
	if score <= 30 && tagCount == 0 {
		return kyt.RiskLevelLow
	}

	// Default: if there are any tags but score is low, still consider Medium
	if tagCount > 0 {
		return kyt.RiskLevelMedium
	}

	return kyt.RiskLevelLow
}

// CalculateTransactionRiskLevel calculates risk level for transactions.
// This is similar to CalculateRiskLevel but operates on transaction-specific data.
func (c *RiskCalculator) CalculateTransactionRiskLevel(score float64, factors []kyt.RiskFactor) kyt.RiskLevel {
	// Check for critical factors
	for _, f := range factors {
		if f.Severity == kyt.RiskLevelCritical || CriticalTags[f.Category] {
			return kyt.RiskLevelCritical
		}
	}

	if score >= 100 {
		return kyt.RiskLevelCritical
	}

	factorCount := len(factors)
	highFactorCount := c.countHighFactors(factors)

	if score >= 71 || factorCount >= 3 || highFactorCount > 0 {
		return kyt.RiskLevelHigh
	}

	mediumFactorCount := c.countMediumFactors(factors)
	if score >= 31 || (mediumFactorCount >= 1 && mediumFactorCount <= 2) {
		return kyt.RiskLevelMedium
	}

	if score <= 30 && factorCount == 0 {
		return kyt.RiskLevelLow
	}

	if factorCount > 0 {
		return kyt.RiskLevelMedium
	}

	return kyt.RiskLevelLow
}

// hasCriticalTags checks if any critical tags are present.
func (c *RiskCalculator) hasCriticalTags(tags []string) bool {
	for _, tag := range tags {
		if CriticalTags[tag] {
			return true
		}
	}
	return false
}

// hasCriticalFactors checks if any factors have critical severity.
func (c *RiskCalculator) hasCriticalFactors(factors []kyt.RiskFactor) bool {
	for _, f := range factors {
		if f.Severity == kyt.RiskLevelCritical || CriticalTags[f.Category] {
			return true
		}
	}
	return false
}

// countMediumTags counts tags considered medium severity.
func (c *RiskCalculator) countMediumTags(tags []string) int {
	count := 0
	for _, tag := range tags {
		if MediumTags[tag] {
			count++
		}
	}
	return count
}

// countHighTags counts tags considered high severity.
func (c *RiskCalculator) countHighTags(tags []string) int {
	count := 0
	for _, tag := range tags {
		if HighTags[tag] {
			count++
		}
	}
	return count
}

// countMediumFactors counts factors with medium severity.
func (c *RiskCalculator) countMediumFactors(factors []kyt.RiskFactor) int {
	count := 0
	for _, f := range factors {
		if f.Severity == kyt.RiskLevelMedium || MediumTags[f.Category] {
			count++
		}
	}
	return count
}

// countHighFactors counts factors with high severity.
func (c *RiskCalculator) countHighFactors(factors []kyt.RiskFactor) int {
	count := 0
	for _, f := range factors {
		if f.Severity == kyt.RiskLevelHigh || HighTags[f.Category] {
			count++
		}
	}
	return count
}

// GetCategorySeverity determines severity based on a category name.
// This is used when the provider doesn't explicitly provide severity.
func (c *RiskCalculator) GetCategorySeverity(category string) kyt.RiskLevel {
	if CriticalTags[category] {
		return kyt.RiskLevelCritical
	}
	if HighTags[category] {
		return kyt.RiskLevelHigh
	}
	if MediumTags[category] {
		return kyt.RiskLevelMedium
	}
	return kyt.RiskLevelLow
}
